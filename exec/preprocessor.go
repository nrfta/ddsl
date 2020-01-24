package exec

import (
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/nrfta/ddsl/drivers/source"
	"github.com/nrfta/ddsl/log"
	"github.com/nrfta/ddsl/parser"
	"path"
	"sort"
	"strings"
)

const (
	// DDSL keywords
	CREATE           string = "create"
	DROP             string = "drop"
	MIGRATE          string = "migrate"
	SEED             string = "seed"
	SQL              string = "sql"
	GRANT            string = "grant"
	REVOKE           string = "revoke"
	LIST             string = "list"
	DATABASE         string = "database"
	DATABASE_PRIVS   string = "database-privs"
	EXTENSIONS       string = "extensions"
	ROLES            string = "roles"
	SCHEMAS          string = "schemas"
	FOREIGN_KEYS     string = "foreign-keys"
	FOREIGN_KEYS_ON  string = "foreign-keys-on"
	SCHEMA           string = "schema"
	SCHEMA_ITEMS     string = "schema-items"
	SCHEMA_PRIVS     string = "schema-privs"
	TABLES           string = "tables"
	TABLES_PRIVS     string = "tables-privs"
	VIEWS            string = "views"
	VIEWS_PRIVS      string = "views-privs"
	FUNCTIONS        string = "functions"
	FUNCTIONS_PRIVS  string = "functions-privs"
	PROCEDURES       string = "procedures"
	PROCEDURES_PRIVS string = "procedures-privs"
	TABLE            string = "table"
	TABLE_PRIVS      string = "table-privs"
	VIEW             string = "view"
	VIEW_PRIVS       string = "view-privs"
	FUNCTION         string = "function"
	FUNCTION_PRIVS   string = "function-privs"
	PROCEDURE        string = "procedure"
	PROCEDURE_PRIVS  string = "procedure-privs"
	INDEXES          string = "indexes"
	CONSTRAINTS      string = "constraints"
	PRIVILEGES       string = "privileges"
	TRIGGERS         string = "triggers"
	TYPE             string = "type"
	TYPES            string = "types"
	CMD              string = "cmd"
	IN               string = "in"
	EXCEPT_IN        string = "except in"

	// param keys
	FILE_PATH    string = "file_path"
	COMMAND      string = "command"
	ARGS         string = "args"
	SCHEMA_NAME  string = "schema_name"
	SCHEMA_NAMES string = "schema_names"
	TABLE_NAME   string = "table_name"
	SEED_NAME    string = "seed_name"
	ITEM_TYPE    string = "item_type"
)

var pathPatterns = map[string]string{
	DATABASE:         `database\.%s\.sql`,
	DATABASE_PRIVS:   `privileges\.%s\.sql`,
	ROLES:            `roles\.%s\.sql`,
	SCHEMAS:          `scheams\.*`,
	FOREIGN_KEYS:     `schemas/%s/tables/?/foreign-keys\.%s\.sql`,
	FOREIGN_KEYS_ON:  `schemas/%s/?/%s/foreign-keys\.%s\.sql`,
	EXTENSIONS:       `extensions\.%s\.sql`,
	SCHEMA:           `schemas/%s/schema\.%s\.sql`,
	SCHEMA_PRIVS:     `schemas/%s/privileges\.%s\.sql`,
	TABLES:           `schemas/%s/tables/?/table\.%s\.sql`,
	TABLES_PRIVS:     `schemas/%s/tables/?/privileges\.%s\.sql`,
	VIEWS:            `schemas/%s/views/?/view\.%s\.sql`,
	VIEWS_PRIVS:      `schemas/%s/views/?/privileges\.%s\.sql`,
	FUNCTIONS:        `schemas/%s/functions/?/function\.%s\.sql`,
	FUNCTIONS_PRIVS:  `schemas/%s/functions/?/privileges\.%s\.sql`,
	PROCEDURES:       `schemas/%s/procedures/?/procedure\.%s\.sql`,
	PROCEDURES_PRIVS: `schemas/%s/procedures/?/privileges\.%s\.sql`,
	TYPES:            `schemas/%s/types/.*\.%s\.sql`,
	TABLE:            `schemas/%s/tables/%s/table\.%s\.sql`,
	TABLE_PRIVS:      `schemas/%s/tables/%s/privileges\.%s\.sql`,
	VIEW:             `schemas/%s/views/%s/view\.%s\.sql`,
	VIEW_PRIVS:       `schemas/%s/views/%s/privileges\.%s\.sql`,
	FUNCTION:         `schemas/%s/functions/%s/function\.%s\.sql`,
	FUNCTION_PRIVS:   `schemas/%s/functions/%s/privileges\.%s\.sql`,
	PROCEDURE:        `schemas/%s/procedures/%s/procedure\.%s\.sql`,
	PROCEDURE_PRIVS:  `schemas/%s/procedures/%s/privileges\.%s\.sql`,
	INDEXES:          `schemas/%s/?/%s/indexes\.%s\.sql`,
	CONSTRAINTS:      `schemas/%s/tables/%s/constraints\.%s\.sql`,
	PRIVILEGES:       `schemas/%s/?/%s/privileges\.%s\.sql`,
	TRIGGERS:         `schemas/%s/tables/%s/triggers\.%s\.sql`,
	TYPE:             `schemas/%s/types/%s\.%s\.sql`,
}

type preprocessor struct {
	ctx           *Context
	sourceDriver  source.Driver
	command       *parser.Command
	createOrDrop  string
	grantOrRevoke string
	databaseName  string
}

type InstructionType int

const (
	INSTR_DDSL InstructionType = iota
	INSTR_SQL_FILE
	INSTR_CSV_FILE
	INSTR_DDSL_FILE
	INSTR_SH_FILE
	INSTR_DDSL_FILE_END
	INSTR_SH_SCRIPT
	INSTR_BEGIN
	INSTR_COMMIT
	INSTR_ROLLBACK
	INSTR_SQL_SCRIPT
	INSTR_LIST
)

type instruction struct {
	instrType InstructionType
	params    map[string]interface{}
}

func preprocessBatch(ctx *Context, cmds []*parser.Command) (int, error) {
	count := 0

	for _, cmd := range cmds {
		// blank lines and comments
		if cmd == nil {
			continue
		}

		ctx.clearPatterns()

		c, err := makeInstructions(ctx, cmd)
		if err != nil {
			return 0, err
		}

		s := "s"
		if c == 1 {
			s = ""
		}
		if c == 0 {
			return 0, fmt.Errorf("no matching files found for %s; patterns tried:\n%s", cmd.Text, ctx.getPatterns())
		}

		if c > 0 {
			count += c
			log.Debug("%d file%s processed", c, s)
			log.Debug("path patterns processed:\n%s", ctx.getPatterns())
		}
	}

	if count == 0 {
		return 0, fmt.Errorf("no files or commands executed")
	}

	return count, nil
}

func makeInstructions(ctx *Context, cmd *parser.Command) (int, error) {
	log.Debug("%sDDSL> %s", ctx.getNestingForLogging(), cmd.Text)
	ctx.addInstructionWithParams(INSTR_DDSL, map[string]interface{}{COMMAND: cmd.Text})

	cmdDef := cmd.CommandDef

	if cmdDef.Name == "begin" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot begin transaction in auto transaction context")
		}
		ctx.addInstruction(INSTR_BEGIN)
		return -1, nil
	}

	if cmdDef.Name == "commit" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot commit transaction in auto transaction context")
		}
		ctx.addInstruction(INSTR_COMMIT)
		return -1, nil
	}

	if cmdDef.Name == "rollback" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot rollback transaction in auto transaction context")
		}
		ctx.addInstruction(INSTR_ROLLBACK)
		return -1, nil
	}

	p := &preprocessor{
		ctx:     ctx,
		command: cmd,
	}
	defer func() {
		if p.sourceDriver != nil {
			p.sourceDriver.Close()
		}
	}()

	count, err := p.preprocess()
	if err != nil {
		return count, err
	}

	return count, nil

}

func (p *preprocessor) preprocess() (int, error) {
	topCmd := p.command.RootDef
	var err error
	var count int
	switch topCmd.Name {
	case CREATE:
		p.createOrDrop = CREATE
		p.grantOrRevoke = GRANT
		count, err = p.preprocessCreateOrDrop()
	case DROP:
		p.createOrDrop = DROP
		p.grantOrRevoke = REVOKE
		count, err = p.preprocessCreateOrDrop()
	case GRANT:
		p.grantOrRevoke = GRANT
		count, err = p.preprocessGrantOrRevoke()
	case REVOKE:
		p.grantOrRevoke = REVOKE
		count, err = p.preprocessGrantOrRevoke()
	case SEED:
		count, err = p.preprocessSeed()
	//case MIGRATE:
	//	count, err = p.preprocessMigrate()
	case SQL:
		count, err = p.preprocessSql()
	case LIST:
		count, err = p.preprocessList()
	default:
		return 0, fmt.Errorf("unknown command")
	}

	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *preprocessor) preprocessSql() (int, error) {
	if len(p.command.ExtArgs) != 1 {
		return 0, fmt.Errorf("the sql command requires one argument")
	}

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing SQL statement")
	if p.ctx.DryRun {
		return 1, nil
	}
	sql := p.command.ExtArgs[0]
	p.ctx.addInstructionWithParams(INSTR_SQL_SCRIPT, map[string]interface{}{SQL: sql})
	err := p.ctx.dbDriver.Exec(strings.NewReader(sql))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (p *preprocessor) makeFileInstructions(pathPattern string) (int, error) {
	if err := p.ensureSourceDriverOpen(); err != nil {
		return 0, err
	}

	relativePath, filePattern := path.Split(pathPattern)
	dirs, err := p.resolveDirectoryWildcards(relativePath)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, d := range dirs {
		p.ctx.addPattern(path.Join(d, filePattern))

		filePaths, err := p.sourceDriver.ReadFiles(d, filePattern)
		if err != nil {
			return 0, err
		}

		count += len(filePaths)

		for _, filePath := range filePaths {
			log.Debug("preprocessing %s", filePath)
			p.ctx.addInstructionWithParams(INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath})
		}
	}
	return count, nil
}

func (p *preprocessor) makeListInstruction(itemType string, params map[string]interface{}) {
	params[ITEM_TYPE] = itemType
	p.ctx.addInstructionWithParams(INSTR_LIST, params)
}

func (p *preprocessor) getSchemaNamesDB(in, except []string) ([]string, error) {
	if in == nil {
		in = []string{}
	}
	if except == nil {
		except = []string{}
	}

	schemaNames, err := p.ctx.dbDriver.Schemas()
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, schemaName := range schemaNames {
		if (len(in) > 0 && sliceutil.Contains(in, schemaName)) ||
			(len(except) > 0 && !sliceutil.Contains(except, schemaName)) ||
			(len(in) == 0 && len(except) == 0) {
			result = append(result, schemaName)
		}
	}

	sort.Strings(result)
	return result, nil
}

func (p *preprocessor) getSchemaNames(in, except []string) ([]string, error) {
	if err := p.ensureSourceDriverOpen(); err != nil {
		return nil, err
	}

	dirs, err := p.getSubdirectories(SCHEMAS_REL_DIR)
	if err != nil {
		return nil, err
	}

	if in == nil {
		in = []string{}
	}
	if except == nil {
		except = []string{}
	}

	schemaNames := []string{}
	for _, d := range dirs {
		schemaName := path.Base(d)
		if (len(in) > 0 && sliceutil.Contains(in, schemaName)) ||
			(len(except) > 0 && !sliceutil.Contains(except, schemaName)) ||
			(len(in) == 0 && len(except) == 0) {
			schemaNames = append(schemaNames, schemaName)
		}
	}

	sort.Strings(schemaNames)
	return schemaNames, nil
}

func (p *preprocessor) getSubdirectories(relativeDir string) ([]string, error) {
	dirs, err := p.resolveDirectoryWildcards(relativeDir)
	if err != nil {
		return nil, err
	}

	dirNames := []string{}
	for _, d := range dirs {

		if err := p.ensureSourceDriverOpen(); err != nil {
			return nil, err
		}

		dirReaders, err := p.sourceDriver.ReadDirectories(d, ".*")
		if err != nil {
			return nil, err
		}

		for _, dr := range dirReaders {
			dirName := path.Base(dr.DirectoryPath)
			dirNames = append(dirNames, dirName)
		}
	}
	return dirNames, nil
}

func (p *preprocessor) resolveDirectoryWildcards(relativeDir string) ([]string, error) {
	if !strings.Contains(relativeDir, "?") {
		return []string{relativeDir}, nil
	}

	i := strings.Index(relativeDir, "?")
	dirs, err := p.getSubdirectories(relativeDir[:i-1])
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, d := range dirs {
		base := path.Base(d)
		names, err := p.resolveDirectoryWildcards(strings.Replace(relativeDir, "?", base, 1))
		if err != nil {
			return nil, err
		}
		result = append(result, names...)
	}

	sort.Strings(result)
	return result, nil
}

func (p *preprocessor) ensureSourceDriverOpen() error {
	if p.sourceDriver != nil {
		return nil
	}

	url := strings.TrimRight(p.ctx.SourceRepo, "/")

	i := strings.LastIndex(url, "/")
	if i == -1 {
		return fmt.Errorf("database name must be last element of DDSL_SOURCE")
	}
	p.databaseName = url[i+1:]

	ref := p.command.Ref
	if ref != nil {
		url += "#" + *ref
	}
	sourceDriver, err := source.Open(url)
	if err != nil {
		return err
	}

	p.sourceDriver = sourceDriver

	return nil
}

func (p *preprocessor) preprocessKey(patternKey string, params ...interface{}) (int, error) {
	switch {
	case len(p.createOrDrop) > 0:
		return p.preprocessCreateOrDropKey(patternKey, params...)
	case len(p.grantOrRevoke) > 0:
		return p.preprocessGrantOrRevokeKey(patternKey, params...)
	default:
		panic("can only be used with create, drop, grant, and revoke")
	}
}

func (p *preprocessor) preprocessSchemas(itemType string) (int, error) {
	count := 0

	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "except":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return count, err
	}

	for _, schemaName := range schemaNames {
		c, err := p.preprocessKey(itemType, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSchema(itemType string) (int, error) {
	count := 0

	for _, schemaName := range p.command.ExtArgs {
		c, err := p.preprocessKey(itemType, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSchemaItem(itemType string) (int, error) {
	count := 0

	if len(p.command.ExtArgs) == 0 {
		return count, fmt.Errorf("comma-delimited list of %ss must be provided", itemType)
	}

	for _, n := range p.command.ExtArgs {
		schemaName, itemName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessKey(itemType, schemaName, itemName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 && p.createOrDrop == CREATE {
			// when creating a table, also create its constraints
			if itemType == TABLE {
				c, err := p.preprocessCreateOrDropKey(CONSTRAINTS, schemaName, itemName)
				count += c
				if err != nil {
					return count, err
				}
			}

			// when creating a table or view, also create its indexes
			if itemType == TABLE || itemType == VIEW {
				c, err = p.preprocessCreateOrDropKey(INDEXES, schemaName, itemName)
				count += c
				if err != nil {
					return count, err
				}
			}

			// when creating any schema item, also grant its privileges
			c, err = p.preprocessGrantOrRevokeKey(PRIVILEGES, schemaName, itemName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil
}

func (p *preprocessor) preprocessSchemaItems(itemType string) (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0

	// before dropping any tables, drop any foreign keys
	if itemType == TABLES && p.createOrDrop == DROP {
		c, err := p.preprocessForeignKeyAttachments(schemaNames)
		count += c
		if err != nil {
			return count, err
		}
	}

	for _, schemaName := range schemaNames {
		c, err := p.preprocessKey(itemType, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if c > 0 && p.createOrDrop == CREATE {
			c, err := p.preprocessSchemaItemAttachments(itemType, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	// after creating all tables, create any foreign keys
	if itemType == TABLES && p.createOrDrop == CREATE {
		c, err := p.preprocessForeignKeyAttachments(schemaNames)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessForeignKeyAttachments(schemaNames []string) (int, error) {
	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessKey(FOREIGN_KEYS, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func (p *preprocessor) preprocessListSchemaItems(itemType string) (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNamesDB(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNamesDB(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNamesDB(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	p.makeListInstruction(itemType, map[string]interface{}{SCHEMA_NAMES: schemaNames})

	return 1, nil
}

func parseSchemaItemName(item string) (schemaName string, tableOrViewName string, err error) {
	if len(item) == 0 {
		return "", "", fmt.Errorf("empty table or view name provided; check for trailing comma or space after comma in list arg")
	}
	nparts := strings.Split(item, ".")
	if len(nparts) != 2 {
		return "", "", fmt.Errorf("tables and views must be defined as <schema_name>.<table_or_view_name>")
	}

	return nparts[0], nparts[1], nil
}
