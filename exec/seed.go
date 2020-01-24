package exec

import (
	"errors"
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/mattn/go-shellwords"
	"github.com/nrfta/ddsl/log"
	"github.com/nrfta/ddsl/parser"
	"io/ioutil"
	"path"
	"sort"
	"strings"
)

const (
	DB_SEEDS_REL_DIR     = "seeds/.*"
	SCHEMAS_REL_DIR      = "schemas"
	SCHEMA_SEEDS_REL_DIR = "schemas/%s/seeds/.*"
	TABLE_SEEDS_REL_DIR  = `schemas/%s/tables/%s/seeds/.*`
	SEED_DATABASE_NAMED  = "database_named"
	SEED_DATABASE_PROD   = "database_prod"
	SEED_SCHEMAS_NAMED   = "schemas_named"
	SEED_SCHEMAS_PROD    = "schemas_prod"
	SEED_SCHEMA_NAMED    = "schema_named"
	SEED_SCHEMA_PROD     = "schema_prod"
	SEED_TABLES_NAMED    = "tables_named"
	SEED_TABLES_PROD     = "tables_prod"
	SEED_TABLE_NAMED     = "table_named"
	SEED_TABLE_PROD      = "table_prod"
)

var seedPatterns = map[string]string{
	SEED_DATABASE_NAMED: `seeds/%s.*`,
	SEED_DATABASE_PROD:  `seeds/database\..*`,
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

func (p *preprocessor) preprocessSeed() (int, error) {
	switch p.command.CommandDef.Name {
	case DATABASE:
		return p.preprocessSeedDatabase()
	case SCHEMAS:
		return p.preprocessSeedSchemas()
	case SCHEMA:
		return p.preprocessSeedSchema()
	case TABLES:
		return p.preprocessSeedTables()
	case TABLE:
		return p.preprocessSeedTable()
	case SQL:
		return p.preprocessSql()
	case CMD:
		return p.seedFromShellCommand()
	}

	return 0, errors.New("unknown command")
}

func (p *preprocessor) preprocessSeedDatabase() (int, error) {
	count := 0
	var seedNames []string
	var err error
	switch p.command.Clause {
	case "with":
		seedNames, err = p.getSeedNames(DB_SEEDS_REL_DIR, p.command.ExtArgs, nil)
	case "without":
		seedNames, err = p.getSeedNames(DB_SEEDS_REL_DIR, nil, p.command.ExtArgs)
	default:
		return p.preprocessSeedKey(SEED_DATABASE_PROD, map[string]interface{}{})
	}
	if err != nil {
		return count, err
	}

	for _, seedName := range seedNames {
		params := map[string]interface{}{
			SEED_NAME: seedName,
		}
		c, err := p.preprocessSeedKey(SEED_DATABASE_NAMED, params, seedName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil

}

func (p *preprocessor) preprocessSeedSchemas() (int, error) {
	if len(p.command.Clause) == 0 {
		return p.preprocessSeedKey(SEED_SCHEMAS_PROD, map[string]interface{}{})
	}

	if p.command.Clause != "except" {
		return 0, fmt.Errorf("syntax error at '%s'", p.command.Clause)
	}

	schemaNames, err := p.getSchemaNames(nil, p.command.ExtArgs)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		pathPattern := strings.Replace(seedPatterns[SEED_SCHEMAS_PROD], "?", schemaName, 1)

		c, err := p.preprocessSeedPath(pathPattern, map[string]interface{}{})
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSeedSchema() (int, error) {
	if len(p.command.Args) < 1 {
		return 0, fmt.Errorf("schema name(s) must be specified")
	}

	count := 0

	schemaNames := p.command.Args
	sort.Strings(schemaNames)
	for _, schemaName := range schemaNames {

		if len(p.command.Clause) == 0 {
			return p.preprocessSeedKey(SEED_SCHEMA_PROD, map[string]interface{}{SCHEMA_NAME: schemaName}, schemaName)
		}

		relativeDir := fmt.Sprintf(SCHEMA_SEEDS_REL_DIR, schemaName)

		var seedNames []string
		var err error
		switch p.command.Clause {
		case "with":
			seedNames, err = p.getSeedNames(relativeDir, p.command.ExtArgs, nil)
		case "without":
			seedNames, err = p.getSeedNames(relativeDir, nil, p.command.ExtArgs)
		default:
			return count, fmt.Errorf("syntax error at '%s'", p.command.Clause)
		}
		if err != nil {
			return count, err
		}

		for _, seedName := range seedNames {
			params := map[string]interface{}{
				SCHEMA_NAME: schemaName,
				SEED_NAME:   seedName,
			}
			c, err := p.preprocessSeedKey(SEED_SCHEMA_NAMED, params, schemaName, seedName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil

}

func (p *preprocessor) preprocessSeedTables() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	case "":
		schemaNames, err = p.getSchemaNames(nil, nil)
	default:
		return 0, fmt.Errorf("syntax error at '%s'", p.command.Clause)
	}

	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		pathPattern := fmt.Sprintf(seedPatterns[SEED_TABLES_PROD], schemaName)
		dirs, err := p.resolveDirectoryWildcards(pathPattern)
		if err != nil {
			return count, err
		}

		for _, dir := range dirs {
			tableName := strings.Split(dir, "/")[3]

			params := map[string]interface{}{
				SCHEMA_NAME: schemaName,
				TABLE_NAME:  tableName,
			}

			c, err := p.preprocessSeedPath(dir, params, schemaName, tableName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSeedTable() (int, error) {
	if len(p.command.Args) < 1 {
		return 0, fmt.Errorf("table name(s) must be specified")
	}

	count := 0

	schemaItems := p.command.Args
	sort.Strings(schemaItems)
	for _, schemaItem := range schemaItems {
		schemaName, tableName, err := parseSchemaItemName(schemaItem)
		if err != nil {
			return count, err
		}
		if len(p.command.Clause) == 0 {
			params := map[string]interface{}{
				SCHEMA_NAME: schemaName,
				TABLE_NAME:  tableName,
			}

			c, err := p.preprocessSeedKey(SEED_TABLE_PROD, params, schemaName, tableName)
			count += c
			if err != nil {
				return count, err
			}

			continue
		}

		relativeDir := fmt.Sprintf(TABLE_SEEDS_REL_DIR, schemaName, tableName)

		var seedNames []string
		switch p.command.Clause {
		case "with":
			seedNames, err = p.getSeedNames(relativeDir, p.command.ExtArgs, nil)
		case "without":
			seedNames, err = p.getSeedNames(relativeDir, nil, p.command.ExtArgs)
		default:
			return count, fmt.Errorf("syntax error at '%s'", p.command.Clause)
		}
		if err != nil {
			return count, err
		}

		for _, seedName := range seedNames {
			params := map[string]interface{}{
				SCHEMA_NAME: schemaName,
				TABLE_NAME:  tableName,
				SEED_NAME:   seedName,
			}
			c, err := p.preprocessSeedKey(SEED_TABLE_NAMED, params, schemaName, tableName, seedName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil

}

func (p *preprocessor) preprocessSeedKey(patternKey string, paramsMap map[string]interface{}, params ...interface{}) (int, error) {
	pathPattern := fmt.Sprintf(seedPatterns[patternKey], params...)
	return p.preprocessSeedPath(pathPattern, paramsMap, params...)
}

func (p *preprocessor) preprocessSeedPath(pathPattern string, paramsMap map[string]interface{}, params ...interface{}) (int, error) {

	relativePath, filePattern := path.Split(pathPattern)

	dirs, err := p.resolveDirectoryWildcards(relativePath)
	if err != nil {
		return 0, err
	}

	if err := p.ensureSourceDriverOpen(); err != nil {
		return 0, err
	}

	count := 0

	for _, dir := range dirs {
		filePaths, err := p.sourceDriver.ReadFiles(dir, filePattern)
		if err != nil {
			return 0, err
		}

		for _, filePath := range filePaths {
			ext := path.Ext(filePath)
			switch ext {
			case ".csv":
				if !strings.Contains(relativePath, "/tables/") {
					return count, fmt.Errorf("only tables can be seeded with CSV: %s", filePath)
				}
				log.Debug("preprocessing CSV seed %s", filePath)
				paramsMap[FILE_PATH] = filePath
				p.ctx.addInstructionWithParams(INSTR_CSV_FILE, paramsMap)
				count++
			case ".sql": // TODO ".sh", ".ddsl":
				log.Debug("preprocessing SQL seed %s", filePath)
				paramsMap[FILE_PATH] = filePath
				p.ctx.addInstructionWithParams(INSTR_SQL_FILE, paramsMap)
				count++
			case ".ddsl":
				log.Debug("preprocessing DDSL seed %s", filePath)
				commandBytes, err := ioutil.ReadFile(filePath)
				if err != nil {
					return count, err
				}

				command := string(commandBytes)
				cmds, _, _, err := parser.Parse(command)
				if err != nil {
					return count, err
				}
				p.ctx.pushNesting()
				p.ctx.addInstructionWithParams(INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath})

				c, err := preprocessBatch(p.ctx, cmds)

				p.ctx.popNesting()
				p.ctx.addInstruction(INSTR_DDSL_FILE_END)

				count += c
				if err != nil {
					return count, err
				}
			default:
				return count, fmt.Errorf("unsupported file %s", filePath)
			}
		}
	}

	return count, nil
}

func (p *preprocessor) seedFromShellCommand() (int, error) {
	if len(p.command.ExtArgs) != 1 {
		return 0, fmt.Errorf("the sql command requires one argument")
	}

	tokens, err := shellParser.Parse(p.command.ExtArgs[0])
	if err != nil {
		return 0, err
	}

	log.Debug("preprocessing seed shell command: %s", p.command.ExtArgs[0])

	command := tokens[0]
	args := []string{}
	if len(tokens) > 1 {
		args = tokens[1:]
	}

	params := map[string]interface{}{
		COMMAND: command,
		ARGS:    args,
	}
	p.ctx.addInstructionWithParams(INSTR_SH_SCRIPT, params)
	return 1, nil
}

func (p *preprocessor) getSeedNames(relativeFilePathPattern string, with, without []string) ([]string, error) {
	if err := p.ensureSourceDriverOpen(); err != nil {
		return nil, err
	}

	relativeDir, filePattern := path.Split(relativeFilePathPattern)
	dirs, err := p.resolveDirectoryWildcards(relativeDir)
	if err != nil {
		return nil, err
	}

	seedNames := []string{}
	for _, d := range dirs {
		filePaths, err := p.sourceDriver.ReadFiles(d, filePattern)
		if err != nil {
			return nil, err
		}

		namedSeed := with != nil || without != nil

		if with == nil {
			with = []string{}
		}
		if without == nil {
			without = []string{}
		}

		for _, filePath := range filePaths {
			i := strings.Index(path.Base(filePath), ".")
			seedName := path.Base(filePath)[:i]
			if sliceutil.Contains(seedNames, seedName) {
				continue
			}
			if namedSeed && (seedName == "database" || seedName == "schema" || seedName == "table") {
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
