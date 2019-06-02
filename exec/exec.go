package exec

import (
	"errors"
	"fmt"
	"github.com/forestgiant/sliceutil"
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"github.com/neighborly/ddsl/drivers/database/postgres"
	"github.com/neighborly/ddsl/drivers/source"
	"github.com/neighborly/ddsl/drivers/source/file"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"path"
	"strings"
)

func init() {
	postgres.Register()
	file.Register()
}

type executor struct {
	ctx          *Context
	sourceDriver source.Driver
	databaseName string
	command      *parser.Command
	createOrDrop string
}

const (
	create string = "create"
	drop   string = "drop"
)

func ExecuteBatch(ctx *Context, cmds []*parser.Command) error {
	dbDriver, err := dbdr.Open(ctx.DatbaseUrl)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	ctx.dbDriver = dbDriver

	err = ensureAuditTable(ctx)
	if err != nil {
		return err
	}

	if ctx.AutoTransaction {
		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "beginning transaction")
		if !ctx.DryRun {
			if err = ctx.dbDriver.Begin(); err != nil {
				return err
			}
		}
		ctx.inTransaction = true
	}

	count := 0

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		c, err := execute(ctx, cmd)
		count += c
		if err != nil {
			if ctx.AutoTransaction && ctx.inTransaction {
				log.Log(levelOrDryRun(ctx, log.LEVEL_WARN), "rolling back transaction")
				if !ctx.DryRun {
					ctx.dbDriver.Rollback()
				}
			}
			return err
		}
	}

	if ctx.AutoTransaction && ctx.inTransaction {
		if count > 0 {
			log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "committing transaction")
			if !ctx.DryRun {
				if err = dbDriver.Commit(); err != nil {
					return err
				}
			}
		} else {
			log.Log(levelOrDryRun(ctx, log.LEVEL_WARN), "rolling back transaction; no commands executed")
			if !ctx.DryRun {
				if err = dbDriver.Rollback(); err != nil {
					return err
				}
			}
		}
	}

	s := "s"
	if count == 1 {
		s = ""
	}
	log.Log(levelOrDryRun(ctx, log.LEVEL_INFO),"%d file%s processed", count, s)

	return nil
}

func execute(ctx *Context, cmd *parser.Command) (int, error) {
	log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "DDSL> %s", cmd.Text)

	cmdDef := cmd.CommandDef

	if cmdDef.Name == "begin" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot begin transaction in auto transaction context")
		}
		if ctx.inTransaction {
			return 0, fmt.Errorf("already in transaction")
		}
		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "beginning transaction")
		if !ctx.DryRun {
			err := ctx.dbDriver.Begin()
			if err != nil {
				return 0, err
			}
		}
		ctx.inTransaction = true
		return 0, nil
	}


	if cmdDef.Name == "commit" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot commit transaction in auto transaction context")
		}
		if !ctx.inTransaction {
			return 0, fmt.Errorf("not in transaction")
		}
		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "committing transaction")
		if !ctx.DryRun {
			err := ctx.dbDriver.Commit()
			if err != nil {
				return 0, err
			}
		}
		return 0, nil
	}

	if cmdDef.Name == "rollback" {
		if ctx.AutoTransaction {
			return 0, fmt.Errorf("cannot rollback transaction in auto transaction context")
		}
		if !ctx.inTransaction {
			return 0, fmt.Errorf("not in transaction")
		}
		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "rolling back transaction")
		if !ctx.DryRun {
			err := ctx.dbDriver.Rollback()
			if err != nil {
				return 0, err
			}
		}
		return 0, nil
	}

	count := 0
	ex := &executor{
		ctx:      ctx,
		command:  cmd,
	}

	// database commands cannot run in transaction
	isDatabaseCommand := cmdDef.Name == "database" && (cmdDef.Parent.Name == "create" || cmdDef.Parent.Name == "drop")
	if isDatabaseCommand {
		if ctx.inTransaction {
			return count, fmt.Errorf("database commands cannot be run within a transaction")
		}
		c, err := ex.executeCmd()
		count += c
		if err != nil {
			return count, err
		}
		return count, nil
	}

	c, err := ex.executeCmd()
	count += c
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ex *executor) executeCmd() (int, error) {
	topCmd := ex.command.RootDef
	var err error
	var count int
	switch topCmd.Name {
	case create:
		ex.createOrDrop = create
		count, err = ex.executeCreateOrDrop()
	case drop:
		ex.createOrDrop = drop
		count, err = ex.executeCreateOrDrop()
	case "seed":
		count, err = ex.executeSeed()
	case "migrate":
		count, err = ex.executeMigrate()
	case "sql":
		log.Log(levelOrDryRun(ex.ctx, log.LEVEL_INFO), "executing SQL statement")
		if ex.ctx.DryRun {
			return 1, nil
		}
		sql := ex.command.Args[0]
		err = ex.ctx.dbDriver.Exec(strings.NewReader(sql))
		count = 1
	default:
		return 0, errors.New("unknown command")
	}

	if err != nil {
		return count, err
	}
	if count > 0 {
		err := ex.audit()
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) getSourceDriver(ref *string) error {
	url := strings.TrimRight(ex.ctx.SourceRepo, "/")

	i := strings.LastIndex(url, "/")
	if i == -1 {
		return errors.New("database name must be last element of DDSL_SOURCE")
	}
	ex.databaseName = url[i+1:]

	if ref != nil {
		url += "#" + *ref
	}
	sourceDriver, err := source.Open(url)
	if err != nil {
		return err
	}

	ex.sourceDriver = sourceDriver

	return nil
}

func (ex *executor) execute(pathPattern string) (int, error) {
	if err := ex.getSourceDriver(ex.command.Ref); err != nil {
		return 0, err
	}
	defer ex.sourceDriver.Close()

	relativePath, filePattern := getRelativePathAndFilePattern(pathPattern)
	readers, err := ex.sourceDriver.ReadFiles(relativePath, filePattern)
	if err != nil {
		return 0, err
	}

	fileCount := len(readers)

	for _, fr := range readers {
		logLevel := log.LEVEL_INFO
		if ex.ctx.DryRun {
			logLevel = log.LEVEL_DRY_RUN
		}
		log.Log(logLevel, "executing %s", fr.FilePath)
		if ex.ctx.DryRun {
			continue
		}
		err = ex.ctx.dbDriver.Exec(fr.Reader)
		if err != nil {
			return fileCount, err
		}
	}

	return fileCount, nil
}

func (ex *executor) getSchemaNames(in, except []string) ([]string, error) {
	if err := ex.getSourceDriver(ex.command.Ref); err != nil {
		return nil, err
	}
	defer ex.sourceDriver.Close()

	dirReaders, err := ex.sourceDriver.ReadDirectories("schemas", ".*")
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
	for _, dr := range dirReaders {
		schemaName := path.Base(dr.DirectoryPath)
		if (len(in) > 0 && sliceutil.Contains(in, schemaName)) ||
			(len(except) > 0 && !sliceutil.Contains(except, schemaName)) ||
			(len(in) == 0 && len(except) == 0) {
			schemaNames = append(schemaNames, schemaName)
		}
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

func levelOrDryRun(ctx *Context, level log.LogLevel) log.LogLevel {
	if ctx.DryRun {
		return log.LEVEL_DRY_RUN
	}
	return level
}

