package exec

import (
	"errors"
	"fmt"
	"github.com/neighborly/ddsl/log"
	"github.com/spf13/viper"
	"path"
	"strings"
)

var seedPatterns = map[string]string{
	tables:      `schemas/%s/tables/.*\.seed\..*`,
	table:       `schemas/%s/tables/%s\.seed\..*`,
}

func (ex *executor) executeSeed() (int, error) {
	switch ex.command.CommandDef.Name {
	case table:
		return ex.executeSeedTable()
	case tables:
		return ex.executeSeedTables()
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) executeSeedTables() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.Args, nil)
	case "except in":
		schemaNames ,err = ex.getSchemaNames(nil, ex.command.Args)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		params := map[string]string {
			"schemaName": schemaName,
		}
		c, err := ex.executeSeedKey(tables, params)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSeedTable() (int, error) {
	count := 0
	for _, n := range ex.command.Args {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		params := map[string]string {
			"schemaName": schemaName,
			"tableName": tableName,
		}
		c, err := ex.executeSeedKey(table, params)
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
	path := fmt.Sprintf(seedPatterns[patternKey], fparams...)
	count, err := ex.executeSeedWork(path, params)
	return count, err
}

func (ex *executor) executeSeedWork(pathPattern string, params map[string]string) (int, error) {
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

	dryRun := viper.GetBool("dry_run")

	for _, fr := range readers {
		ext := path.Ext(fr.FilePath)
		action := ""
		switch ext {
		case ".csv":
			action = "seeding with CSV"
		case ".sql": // TODO ".sh", ".ddsl":
			action = "seeding with SQL"
		default:
			return count, fmt.Errorf("[ERROR] unsupported file %s", fr.FilePath)
		}

		count++

		logLevel := log.LEVEL_INFO
		if dryRun {
			logLevel = log.LEVEL_DRY_RUN
		}
		log.Log(logLevel,"%s %s\n", action, fr.FilePath)
		if dryRun {
			continue
		}

		switch ext {
		case ".sql":
			if err = ex.dbDriver.Exec(fr.Reader); err != nil {
				return count, err
			}
		case ".csv":
			filename := path.Base(fr.FilePath)
			tablename := strings.Split(filename, ".")[0]
			if err = ex.dbDriver.ImportCSV(fr.FilePath, params["schemaName"], tablename,",", true); err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

