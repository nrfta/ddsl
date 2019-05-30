package exec

import (
	"errors"
	"fmt"
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"github.com/neighborly/ddsl/drivers/database/postgres"
	"github.com/neighborly/ddsl/drivers/source"
	"github.com/neighborly/ddsl/drivers/source/file"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/viper"
	"path"
	"strings"
)

func init() {
	postgres.Register()
	file.Register()
}

type executor struct {
	ctx *Context
	sourceDriver source.Driver
	dbDriver     dbdr.Driver
	databaseName string
	cmdDef *parser.CommandDef
	createOrDrop string
}

const (
	create string = "create"
	drop   string = "drop"
)

func Execute(ctx *Context, command string) error {
	cmdDef, err := parser.Parse(command)
	if err != nil {
		return err
	}

	dbDriver, err := dbdr.Open(ctx.DatbaseUrl)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	count := 0

	// database commands cannot run in transaction
	isDatabaseCommand := cmdDef.Name == "database" && (cmdDef.Parent.Name == "create" || cmdDef.Parent.Name == "drop")
	if isDatabaseCommand {
		if ctx.inTransaction {
			return fmt.Errorf("database commands cannot be run within a transaction")
		}
		ex := &executor{
			ctx: ctx,
			dbDriver:  dbDriver,
			cmdDef: cmdDef,
		}
		c, err := ex.execute()
		count += c
		if err != nil {
			return err
		}
		return nil
	}

	dryRun := viper.GetBool("dry_run")

	logLevel := "INFO"
	if dryRun {
		logLevel = "DRY-RUN"
	}

	if ctx.AutoTransaction {
		log.log(logLevel, "beginning transaction")
		if !dryRun {
			err = dbDriver.Begin()
			if err != nil {
				return err
			}
		}
	}

	ex := &executor{
		ctx: ctx,
		dbDriver: dbDriver,
		cmdDef: cmdDef,
	}

	c, err := ex.execute()
	count += c
	if err != nil {
		logLevel = "WARN"
		if dryRun {
			logLevel = "DRY-RUN"
		}
		log.log(logLevel, "rolling back transaction\n", logLevel)
		if !dryRun {
			_ = dbDriver.Rollback()
		}
		return err
	}

	if ctx.AutoTransaction {
		log.log(logLevel, "committing transaction")
		if err = dbDriver.Commit(); err != nil {
			return err
		}
	}

	if count == 0 {
		log.Warn("*** command did nothing; no files matched ***")
	} else {
		log.Info("*** %d files processed ***\n", count)
	}

	return nil
}

func (ex *executor) execute() (int, error) {
	topCmd := ex.cmdDef.ParentAtLevel(1)
	switch topCmd.Name {
	case create:
		ex.createOrDrop = create
		return ex.executeCreateOrDrop()
	case drop:
		ex.createOrDrop = drop
		return ex.executeCreateOrDrop()
	case "migrate":
		return ex.executeMigrate()
	case "sql":
		dryRun := viper.GetBool("dry_run")
		logLevel := "INFO"
		if dryRun {
			logLevel = "DRY-RUN"
		}
		log.log(logLevel, "executing SQL statement\n", logLevel)
		if dryRun {
			return 1, nil
		}
		return 1, ex.dbDriver.Exec(strings.NewReader(*ex.parseTree.Sql))
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) getSourceDriver(ref *parser.Ref) error {
	url := strings.TrimRight(ex.repo, "/")

	i := strings.LastIndex(url, "/")
	if i == -1 {
		return errors.New("database name must be last element of DDSL_SOURCE")
	}
	ex.databaseName = url[i+1:]

	if ref != nil {
		url += "#" + strings.TrimLeft(ref.Ref, "@")
	}
	sourceDriver, err := source.Open(url)
	if err != nil {
		return err
	}

	ex.sourceDriver = sourceDriver

	return nil
}

func (ex *executor) execute(pathPattern string, ref *parser.Ref) (int, error) {
	if err := ex.getSourceDriver(ref); err != nil {
		return 0, err
	}
	defer ex.sourceDriver.Close()

	relativePath, filePattern := getRelativePathAndFilePattern(pathPattern)
	readers, err := ex.sourceDriver.ReadFiles(relativePath, filePattern)
	if err != nil {
		return 0, err
	}

	fileCount := len(readers)

	dryRun := viper.GetBool("dry_run")

	for _, fr := range readers {
		logLevel := "INFO"
		if dryRun {
			logLevel = "DRY-RUN"
		}
		log.log(logLevel, "executing %s\n", logLevel, fr.FilePath)
		if dryRun {
			continue
		}
		err = ex.dbDriver.Exec(fr.Reader)
		if err != nil {
			return fileCount, err
		}
	}

	return fileCount, nil
}

func (ex *executor) getSchemaNames(ref *parser.Ref) ([]string, error) {
	if err := ex.getSourceDriver(ref); err != nil {
		return nil, err
	}
	defer ex.sourceDriver.Close()

	dirReaders, err := ex.sourceDriver.ReadDirectories("schemas", ".*")
	if err != nil {
		return nil, err
	}

	schemaNames := []string{}
	for _, dr := range dirReaders {
		schemaNames = append(schemaNames, path.Base(dr.DirectoryPath))
	}

	return schemaNames, nil
}

func getRelativePathAndFilePattern(path string) (relativePath string, filePattern string) {
	p := strings.TrimRight(path, "/")

	i := strings.LastIndex(p, "/")
	if i == -1 {
		return "", p
	}

	return p[:i], p[i+1:]
}
