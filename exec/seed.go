package exec

import (
	"errors"
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/mattn/go-shellwords"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	DatabaseSeedsRelativeDir = "seeds/.*"
	SchemasRelativeDir     = "schemas"
	SchemaSeedsRelativeDir = "schemas/%s/seeds/.*"
	TableSeedsRelativeDir  = `schemas/%s/tables/%s/seeds/.*`
	CMD                    = "cmd"
	SEED_DATABASE_NAMED    = "database_named"
	SEED_DATABASE_PROD     = "database_prod"
	SEED_SCHEMAS_NAMED     = "schemas_named"
	SEED_SCHEMAS_PROD      = "schemas_prod"
	SEED_SCHEMA_NAMED      = "schema_named"
	SEED_SCHEMA_PROD       = "schema_prod"
	SEED_TABLES_NAMED      = "tables_named"
	SEED_TABLES_PROD       = "tables_prod"
	SEED_TABLE_NAMED       = "table_named"
	SEED_TABLE_PROD        = "table_prod"
)

var seedPatterns = map[string]string{
	SEED_DATABASE_NAMED: `seeds/%s.*`,
	SEED_DATABASE_PROD:  `seeds/database\.seed\..*`,
	SEED_SCHEMAS_NAMED:  `schemas/?/seeds/%s\..*`,
	SEED_SCHEMAS_PROD:   `schemas/?/seeds/schema\..*`,
	SEED_SCHEMA_NAMED:   `schemas/%s/seeds/%s\..*`,
	SEED_SCHEMA_PROD:    `schemas/%s/seeds/schema\..*`,
	SEED_TABLES_NAMED:   `schemas/%s/tables/?/%s/seeds/.*`,
	SEED_TABLES_PROD:    `schemas/%s/tables/?/seeds/table\..*`,
	SEED_TABLE_NAMED:    `schemas/%s/tables/%s/seeds/%s\..*`,
	SEED_TABLE_PROD:     `schemas/%s/tables/%s/seeds/table\..*`,
}

var shellParser *shellwords.Parser

func init() {
	shellParser = shellwords.NewParser()
	shellParser.ParseEnv = true
	shellParser.ParseBacktick = true
}

func (ex *executor) executeSeed() (int, error) {
	switch ex.command.CommandDef.Name {
	case DATABASE:
		return ex.executeSeedDatabase()
	case SCHEMAS:
		return ex.executeSeedSchemas()
	case SCHEMA:
		return ex.executeSeedSchema()
	case TABLES:
		return ex.executeSeedTables()
	case TABLE:
		return ex.executeSeedTable()
	case SQL:
		return ex.executeSql()
	case CMD:
		return ex.seedFromShellCommand()
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) executeSeedDatabase() (int, error) {
	count := 0
	var seedNames []string
	var err error
	switch ex.command.Clause {
	case "with":
		seedNames, err = ex.getSeedNames(DatabaseSeedsRelativeDir, ex.command.ExtArgs, nil)
	case "without":
		seedNames, err = ex.getSeedNames(DatabaseSeedsRelativeDir, nil, ex.command.ExtArgs)
	default:
		return ex.executeSeedKey(SEED_DATABASE_PROD, map[string]string{})
	}
	if err != nil {
		return count, err
	}

	for _, seedName := range seedNames {
		params := map[string]string{
			"seedName": seedName,
		}
		c, err := ex.executeSeedKey(SEED_DATABASE_NAMED, params, seedName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil

}

func (ex *executor) executeSeedSchemas() (int, error) {
	if len(ex.command.Clause) == 0 {
		return ex.executeSeedKey(SEED_SCHEMAS_PROD, map[string]string{})
	}

	schemaNames, err := ex.getSchemaNames(nil, nil)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		relativeDir := fmt.Sprintf(SchemaSeedsRelativeDir, schemaName)

		var seedNames []string
		var err error
		switch ex.command.Clause {
		case "with":
			seedNames, err = ex.getSeedNames(relativeDir, ex.command.ExtArgs, nil)
		case "without":
			seedNames, err = ex.getSeedNames(relativeDir, nil, ex.command.ExtArgs)
		default:
			return count, fmt.Errorf("syntax error at '%s'", ex.command.Clause)
		}

		if err != nil {
			return count, err
		}

		for _, seedName := range seedNames {
			params := map[string]string{
				"seedName": seedName,
			}
			c, err := ex.executeSeedKey(SEED_SCHEMAS_NAMED, params, seedName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (ex *executor) executeSeedSchema() (int, error) {
	if len(ex.command.Args) < 1 {
		return 0, fmt.Errorf("schema name(s) must be specified")
	}

	count := 0
	for _, schemaName := range ex.command.Args {

		if len(ex.command.Clause) == 0 {
			return ex.executeSeedKey(SEED_SCHEMA_PROD, map[string]string{"schemaName": schemaName}, schemaName)
		}

		relativeDir := fmt.Sprintf(SchemaSeedsRelativeDir, schemaName)

		var seedNames []string
		var err error
		switch ex.command.Clause {
		case "with":
			seedNames, err = ex.getSeedNames(relativeDir, ex.command.ExtArgs, nil)
		case "without":
			seedNames, err = ex.getSeedNames(relativeDir, nil, ex.command.ExtArgs)
		default:
			return count, fmt.Errorf("syntax error at '%s'", ex.command.Clause)
		}
		if err != nil {
			return count, err
		}

		for _, seedName := range seedNames {
			params := map[string]string{
				"schemaName": schemaName,
				"seedName":   seedName,
			}
			c, err := ex.executeSeedKey(SEED_SCHEMA_NAMED, params, schemaName, seedName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil

}

func (ex *executor) executeSeedTables() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	case "":
		schemaNames, err = ex.getSchemaNames(nil, nil)
	default:
		return 0, fmt.Errorf("syntax error at '%s'", ex.command.Clause)
	}

	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		params := map[string]string{
			"schemaName": schemaName,
		}
		c, err := ex.executeSeedKey(SEED_TABLES_PROD, params, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSeedTable() (int, error) {
	if len(ex.command.Args) < 1 {
		return 0, fmt.Errorf("table name(s) must be specified")
	}

	count := 0
	for _, schemaItem := range ex.command.Args {
		schemaName, tableName, err := parseSchemaItemName(schemaItem)
		if err != nil {
			return count, err
		}
		if len(ex.command.Clause) == 0 {
			return ex.executeSeedKey(SEED_TABLE_PROD, map[string]string{"schemaName": schemaName, "tableName": tableName}, schemaName, tableName)
		}

		relativeDir := fmt.Sprintf(TableSeedsRelativeDir, schemaName, tableName)

		var seedNames []string
		switch ex.command.Clause {
		case "with":
			seedNames, err = ex.getSeedNames(relativeDir, ex.command.ExtArgs, nil)
		case "without":
			seedNames, err = ex.getSeedNames(relativeDir, nil, ex.command.ExtArgs)
		default:
			return count, fmt.Errorf("syntax error at '%s'", ex.command.Clause)
		}
		if err != nil {
			return count, err
		}

		for _, seedName := range seedNames {
			params := map[string]string{
				"schemaName": schemaName,
				"tableName": tableName,
				"seedName":   seedName,
			}
			c, err := ex.executeSeedKey(SEED_TABLE_NAMED, params, schemaName, tableName, seedName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil

}


func (ex *executor) executeSeedKey(patternKey string, paramsMap map[string]string, params ...interface{}) (int, error) {
	pathPattern := fmt.Sprintf(seedPatterns[patternKey], params...)

	if err := ex.ensureSourceDriverOpen(); err != nil {
		return 0, err
	}

	relativePath, filePattern := path.Split(pathPattern)
	readers, err := ex.sourceDriver.ReadFiles(relativePath, filePattern)
	if err != nil {
		return 0, err
	}

	count := 0

	for _, fr := range readers {
		ext := path.Ext(fr.FilePath)
		action := ""
		switch ext {
		case ".csv":
			if !strings.Contains(patternKey, "table") {
				return count, fmt.Errorf("only tables can be seeded with CSV: %s", fr.FilePath)
			}
			action = "seeding with CSV"
		case ".ddsl":
			action = "seeding with DDSL"
		case ".sql": // TODO ".sh", ".ddsl":
			action = "seeding with SQL"
		default:
			return count, fmt.Errorf("unsupported file %s", fr.FilePath)
		}

		log.Log(levelOrDryRun(ex.ctx, log.LEVEL_INFO), "%s %s", action, fr.FilePath)
		if ex.ctx.DryRun {
			continue
		}

		switch ext {
		case ".sql":
			count++
			if err = ex.ctx.dbDriver.Exec(fr.Reader); err != nil {
				return count, err
			}
		case ".csv":
			count++
			schemaName := paramsMap["schemaName"]
			tableName := paramsMap["tableName"]

			if err = ex.ctx.dbDriver.ImportCSV(fr.FilePath, schemaName, tableName, ",", true); err != nil {
				return count, err
			}
		case ".ddsl":
			commandBytes, err := ioutil.ReadFile(fr.FilePath)
			if err != nil {
				return count, err
			}

			command := string(commandBytes)
			cmds, _, _, err := parser.Parse(command)
			if err != nil {
				return count, err
			}
			ex.ctx.pushNesting()
			c, err := executeBatch(ex.ctx, cmds)
			ex.ctx.popNesting()
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (ex *executor) seedFromShellCommand() (int, error) {
	if len(ex.command.ExtArgs) != 1 {
		return 0, fmt.Errorf("the sql command requires one argument")
	}

	tokens, err := shellParser.Parse(ex.command.ExtArgs[0])
	if err != nil {
		return 0, err
	}

	log.Log(levelOrDryRun(ex.ctx, log.LEVEL_INFO), "seeding with shell command: %s", ex.command.ExtArgs[0])
	if ex.ctx.DryRun {
		return 1, nil
	}

	command := tokens[0]
	args := []string{}
	if len(tokens) > 1 {
		args = tokens[1:]
	}
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(),
		"DDSL_SOURCE="+ex.ctx.SourceRepo,
		"DDSL_DATABASE="+ex.ctx.DatbaseUrl,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, err
	}
	if err = cmd.Start(); err != nil {
		return 0, err
	}
	out, _ := ioutil.ReadAll(stdout)
	e, _ := ioutil.ReadAll(stderr)
	if err = cmd.Wait(); err != nil {
		s := getOutAndErr(out, e)
		if len(s) > 0 {
			return 0, fmt.Errorf(string(s))
		}
		return 0, err
	}

	s := getOutAndErr(out, e)
	if len(s) > 0 {
		fmt.Println(s)
	}
	return 1, nil

}

func getOutAndErr(out, e []uint8) string {
	s := ""
	if len(out) > 0 {
		s = string(out)
	}
	if len(e) > 0 {
		s += string(e)
	}
	return s
}

func (ex *executor) getSeedNames(relativeFilePathPattern string, with, without []string) ([]string, error) {
	if err := ex.ensureSourceDriverOpen(); err != nil {
		return nil, err
	}

	relativeDir, filePattern := path.Split(relativeFilePathPattern)
	dirs, err := ex.resolveDirectoryWildcards(relativeDir)
	if err != nil {
		return nil, err
	}

	seedNames := []string{}
	for _, d := range dirs {
		frs, err := ex.sourceDriver.ReadFiles(d, filePattern)
		if err != nil {
			return nil, err
		}

		if with == nil {
			with = []string{}
		}
		if without == nil {
			without = []string{}
		}

		for _, fr := range frs {
			i := strings.Index(path.Base(fr.FilePath), ".")
			seedName := path.Base(fr.FilePath)[:i]
			if sliceutil.Contains(seedNames, seedName) {
				continue
			}
			if (len(with) > 0 && sliceutil.Contains(with, seedName)) ||
				(len(without) > 0 && !sliceutil.Contains(without, seedName)) ||
				(len(with) == 0 && len(without) == 0) {
				seedNames = append(seedNames, seedName)
			}
		}
	}
	return seedNames, nil
}
