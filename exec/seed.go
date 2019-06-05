package exec

import (
	"errors"
	"fmt"
	"github.com/forestgiant/sliceutil"
	"github.com/mattn/go-shellwords"
	"github.com/neighborly/ddsl/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	DatabaseRelativeDir = "seeds"
	SchemaRelativeDir   = "schemas/%s/seeds"
	CMD                 = "cmd"
)

var seedPatterns = map[string]string{
	TABLES:   `schemas/%s/tables/.*\.seed\..*`,
	TABLE:    `schemas/%s/tables/%s\.seed\..*`,
	DATABASE: `seeds/%s\..*`,
	SCHEMA:   `schemas/%s/seeds/%s\..*`,
}

var shellParser *shellwords.Parser

func init() {
	shellParser = shellwords.NewParser()
	shellParser.ParseEnv = true
	shellParser.ParseBacktick = true
}

func (ex *executor) executeSeed() (int, error) {
	switch ex.command.CommandDef.Name {
	case TABLE:
		return ex.executeSeedTable()
	case TABLES:
		return ex.executeSeedTables()
	case DATABASE:
		return ex.executeSeedDatabase()
	case SQL:
		return ex.executeSql()
	case SCHEMA:
		return ex.executeSeedSchema()
	case CMD:
		return ex.seedFromShellCommand()
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) executeSeedTables() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		params := map[string]string{
			"schemaName": schemaName,
		}
		c, err := ex.executeSeedKey(TABLES, params)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSeedTable() (int, error) {
	count := 0
	for _, n := range ex.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		params := map[string]string{
			"schemaName": schemaName,
			"tableName":  tableName,
		}
		c, err := ex.executeSeedKey(TABLE, params)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSeedDatabase() (int, error) {
	count := 0
	var seedNames []string
	var err error
	switch ex.command.Clause {
	case "with":
		seedNames, err = ex.getSeedNames(DatabaseRelativeDir, ex.command.ExtArgs, nil)
	case "without":
		seedNames, err = ex.getSeedNames(DatabaseRelativeDir, nil, ex.command.ExtArgs)
	default:
		seedNames, err = ex.getSeedNames(DatabaseRelativeDir, nil, nil)
	}
	if err != nil {
		return count, err
	}

	for _, seedName := range seedNames {
		params := map[string]string{
			"seedName": seedName,
		}
		c, err := ex.executeSeedKey(DATABASE, params)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil

}

func (ex *executor) executeSeedSchema() (int, error) {
	if len(ex.command.Args) != 1 {
		return 0, fmt.Errorf("schema name must be specified")
	}
	schemaName := ex.command.Args[0]
	relativeDir := fmt.Sprintf(SchemaRelativeDir, schemaName)
	count := 0
	var seedNames []string
	var err error
	switch ex.command.Clause {
	case "with":
		seedNames, err = ex.getSeedNames(relativeDir, ex.command.ExtArgs, nil)
	case "without":
		seedNames, err = ex.getSeedNames(relativeDir, nil, ex.command.ExtArgs)
	default:
		seedNames, err = ex.getSeedNames(relativeDir, nil, nil)
	}
	if err != nil {
		return count, err
	}

	for _, seedName := range seedNames {
		params := map[string]string{
			"schemaName": schemaName,
			"seedName": seedName,
		}
		c, err := ex.executeSeedKey(SCHEMA, params)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil

}

func (ex *executor) executeSeedKey(patternKey string, params map[string]string) (int, error) {
	fparams := []interface{}{}
	for _, p := range params {
		fparams = append(fparams, p)
	}
	pathPattern := fmt.Sprintf(seedPatterns[patternKey], fparams...)

	if err := ex.getSourceDriver(ex.command.Ref); err != nil {
		return 0, err
	}
	defer ex.sourceDriver.Close()

	relativePath, filePattern := getRelativePathAndFilePattern(pathPattern)
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
			if patternKey != TABLE && patternKey != TABLES {
				return count, fmt.Errorf("only tables can be seeded with CSV: %s", fr.FilePath)
			}
			action = "seeding with CSV"
		case ".sql": // TODO ".sh", ".ddsl":
			action = "seeding with SQL"
		default:
			return count, fmt.Errorf("unsupported file %s", fr.FilePath)
		}

		count++

		log.Log(levelOrDryRun(ex.ctx, log.LEVEL_INFO), "%s %s", action, fr.FilePath)
		if ex.ctx.DryRun {
			continue
		}

		switch ext {
		case ".sql":
			if err = ex.ctx.dbDriver.Exec(fr.Reader); err != nil {
				return count, err
			}
		case ".csv":
			filename := path.Base(fr.FilePath)
			tablename := strings.Split(filename, ".")[0]
			if err = ex.ctx.dbDriver.ImportCSV(fr.FilePath, params["schemaName"], tablename, ",", true); err != nil {
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

func (ex *executor) getSeedNames(relativeDir string, with, without []string) ([]string, error) {
	if err := ex.getSourceDriver(ex.command.Ref); err != nil {
		return nil, err
	}
	defer ex.sourceDriver.Close()

	frs, err := ex.sourceDriver.ReadFiles(relativeDir, ".*")
	if err != nil {
		return nil, err
	}

	if with == nil {
		with = []string{}
	}
	if without == nil {
		without = []string{}
	}

	seedNames := []string{}
	for _, fr := range frs {
		i := strings.LastIndex(path.Base(fr.FilePath), path.Ext(fr.FilePath))
		seedName := path.Base(fr.FilePath)[:i]
		if (len(with) > 0 && sliceutil.Contains(with, seedName)) ||
			(len(without) > 0 && !sliceutil.Contains(without, seedName)) ||
			(len(with) == 0 && len(without) == 0) {
			seedNames = append(seedNames, seedName)
		}
	}

	return seedNames, nil
}
