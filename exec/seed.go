package exec

import (
	"errors"
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/viper"
	"path"
	"strings"
)

var seedPatterns = map[string]string{
	tables:      `schemas/%s/tables/.*\.seed\..*`,
	table:       `schemas/%s/tables/%s\.seed\..*`,
}

func executeSeed(ex *executor) (int, error) {
	cmd := ex.parseTree.Seed

	switch {
	case cmd.Table != nil:
		return executeSeedTableOrView(ex, table, cmd.Table.Schema, cmd.Table.TableOrView, cmd.Table.Ref)
	case cmd.TablesInSchema != nil:
		return executeSeedTables(ex, cmd.TablesInSchema.Name, cmd.TablesInSchema.Ref)
	}

	return 0, errors.New("unknown command")
}

func executeSeedTables(ex *executor, schemaName string, ref *parser.Ref) (int, error) {
	params := map[string]string{
		"schemaName": schemaName,
	}
	return executeSeedKey(ex, ref, tables, params)
}

func executeSeedTableOrView(ex *executor, tableOrView string, schemaName string, tableName string, ref *parser.Ref) (int, error) {
	params := map[string]string{
		"schemaName": schemaName,
		"tableName": tableName,
	}
	return executeSeedKey(ex, ref, tableOrView, params)
}

func executeSeedKey(ex *executor, ref *parser.Ref, patternKey string, params map[string]string) (int, error) {
	fparams := []interface{}{}
	for _, s := range params {
		fparams = append(fparams, s)
	}
	path := fmt.Sprintf(seedPatterns[patternKey], fparams...)
	count, err := ex.executeSeed(path, ref, params)
	return count, err
}

func (ex *executor) executeSeed(pathPattern string, ref *parser.Ref, params map[string]string) (int, error) {
	if err := ex.getSourceDriver(ref); err != nil {
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

		logLevel := "INFO"
		if dryRun {
			logLevel = "DRY-RUN"
		}
		fmt.Printf("[%s] %s %s\n", logLevel, action, fr.FilePath)
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

