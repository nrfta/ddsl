package exec

import (
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/neighborly/ddsl/drivers/source"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"path"
	"sort"
	"strings"
)

const (
	// DDSL keywords
	CREATE       string = "create"
	DROP         string = "drop"
	MIGRATE      string = "migrate"
	SEED         string = "seed"
	SQL          string = "sql"
	GRANT        string = "grant"
	REVOKE       string = "revoke"
	DATABASE     string = "database"
	EXTENSIONS   string = "extensions"
	ROLES        string = "roles"
	SCHEMAS      string = "schemas"
	KEYS         string = "keys"
	FOREIGN_KEYS string = "foreign_keys"
	SCHEMA       string = "schema"
	TABLES       string = "tables"
	VIEWS        string = "views"
	FUNCTIONS    string = "functions"
	PROCEDURES   string = "procedures"
	TABLE        string = "table"
	VIEW         string = "view"
	FUNCTION     string = "function"
	PROCEDURE    string = "procedure"
	INDEXES      string = "indexes"
	CONSTRAINTS  string = "constraints"
	PRIVILEGES   string = "privileges"
	TRIGGERS     string = "triggers"
	TYPE         string = "type"
	TYPES        string = "types"
	CMD          string = "cmd"

	// param keys
	FILE_PATH   string = "file_path"
	COMMAND     string = "command"
	ARGS        string = "args"
	SCHEMA_NAME string = "schema_name"
	TABLE_NAME  string = "table_name"
	SEED_NAME   string = "seed_name"
)

var pathPatterns = map[string]string{
	DATABASE:     `database\.%s\.sql`,
	ROLES:        `roles\.%s\.sql`,
	SCHEMAS:      `scheams\.*`,
	FOREIGN_KEYS: `foreign_keys\.%s\.sql`,
	SCHEMA:       `schemas/%s/schema\.%s\.sql`,
	EXTENSIONS:   `schemas/%s/extensions\.%s\.sql`,
	TABLES:       `schemas/%s/tables/?/table\.%s\.sql`,
	VIEWS:        `schemas/%s/views/?/view\.%s\.sql`,
	FUNCTIONS:    `schemas/%s/functions/?/function\.%s\.sql`,
	PROCEDURES:   `schemas/%s/procedures/?/procedure\.%s\.sql`,
	TYPES:        `schemas/%s/types/.*\.%s\.sql`,
	TABLE:        `schemas/%s/tables/%s/table\.%s\.sql`,
	VIEW:         `schemas/%s/views/%s/view\.%s\.sql`,
	INDEXES:      `schemas/%s/?/%s/indexes\.%s\.sql`,
	CONSTRAINTS:  `schemas/%s/tables/%s/constraints\.%s\.sql`,
	PRIVILEGES:   `schemas/%s/?/%s/privileges\.%s\.sql`,
	TRIGGERS:     `schemas/%s/tables/%s/triggers\.%s\.sql`,
	TYPE:         `schemas/%s/types/%s\.%s\.sql`,
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
	case REVOKE:
		p.grantOrRevoke = REVOKE
	case SEED:
		count, err = p.preprocessSeed()
	//case MIGRATE:
	//	count, err = p.preprocessMigrate()
	case SQL:
		count, err = p.preprocessSql()
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

		readers, err := p.sourceDriver.ReadFiles(d, filePattern)
		if err != nil {
			return 0, err
		}

		count += len(readers)

		for _, fr := range readers {
			log.Debug("preprocessing %s", fr.FilePath)
			p.ctx.addInstructionWithParams(INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: fr.FilePath})
		}
	}
	return count, nil
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
