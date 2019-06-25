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
	CREATE  string = "create"
	DROP    string = "drop"
	MIGRATE string = "migrate"
	SEED    string = "seed"
	SQL     string = "sql"
	GRANT   string = "grant"
	REVOKE  string = "revoke"
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

	_, err = executeBatch(ctx, cmds)
	if err != nil {
		if ctx.AutoTransaction && ctx.inTransaction {
			log.Log(levelOrDryRun(ctx, log.LEVEL_WARN), "rolling back transaction")
			if !ctx.DryRun {
				ctx.dbDriver.Rollback()
			}
		}
		return err
	}

	if ctx.AutoTransaction && ctx.inTransaction {
		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "committing transaction")
		if !ctx.DryRun {
			if err = dbDriver.Commit(); err != nil {
				return err
			}
		}
	}

	return nil
}

func executeBatch(ctx *Context, cmds []*parser.Command) (int, error) {
	count := 0

	for _, cmd := range cmds {
		// blank lines and comments
		if cmd == nil {
			continue
		}

		ctx.clearPatterns()

		c, err := executeCmd(ctx, cmd)
		if err != nil {
			return count, err
		}

		s := "s"
		if c == 1 {
			s = ""
		}
		if c == 0 {
			return 0, fmt.Errorf("no matching files found; patterns tried:\n%s", ctx.getPatterns())
		}

		log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "%d file%s processed", c, s)
		log.Debug("path patterns processed:\n%s", ctx.getPatterns())

		count += c
	}

	return count, nil
}

func executeCmd(ctx *Context, cmd *parser.Command) (int, error) {
	log.Log(levelOrDryRun(ctx, log.LEVEL_INFO), "%sDDSL> %s", ctx.getNestingForLogging(), cmd.Text)

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
		ctx:     ctx,
		command: cmd,
	}
	defer func() {
		if ex.sourceDriver != nil {
			ex.sourceDriver.Close()
		}
	}()

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
	case CREATE:
		ex.createOrDrop = CREATE
		count, err = ex.executeCreateOrDrop()
	case DROP:
		ex.createOrDrop = DROP
		count, err = ex.executeCreateOrDrop()
	case SEED:
		count, err = ex.executeSeed()
	case MIGRATE:
		count, err = ex.executeMigrate()
	case SQL:
		count, err = ex.executeSql()
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

func (ex *executor) ensureSourceDriverOpen() error {
	if ex.sourceDriver != nil {
		return nil
	}

	url := strings.TrimRight(ex.ctx.SourceRepo, "/")

	i := strings.LastIndex(url, "/")
	if i == -1 {
		return errors.New("database name must be last element of DDSL_SOURCE")
	}
	ex.databaseName = url[i+1:]

	ref := ex.command.Ref
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
	if err := ex.ensureSourceDriverOpen(); err != nil {
		return 0, err
	}

	relativePath, filePattern := path.Split(pathPattern)
	dirs, err := ex.resolveDirectoryWildcards(relativePath)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, d := range dirs {
		ex.ctx.addPattern(path.Join(d, filePattern))

		readers, err := ex.sourceDriver.ReadFiles(d, filePattern)
		if err != nil {
			return 0, err
		}

		count += len(readers)

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
				return count, err
			}
		}
	}
	return count, nil
}

func (ex *executor) executeSql() (int, error) {
	if len(ex.command.ExtArgs) != 1 {
		return 0, fmt.Errorf("the sql command requires one argument")
	}

	log.Log(levelOrDryRun(ex.ctx, log.LEVEL_INFO), "executing SQL statement")
	if ex.ctx.DryRun {
		return 1, nil
	}
	sql := ex.command.ExtArgs[0]
	err := ex.ctx.dbDriver.Exec(strings.NewReader(sql))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (ex *executor) getSchemaNames(in, except []string) ([]string, error) {
	if err := ex.ensureSourceDriverOpen(); err != nil {
		return nil, err
	}

	dirs, err := ex.getSubdirectories(SchemasRelativeDir)
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

	return schemaNames, nil
}

func (ex *executor) getSubdirectories(relativeDir string) ([]string, error) {
	dirs, err := ex.resolveDirectoryWildcards(relativeDir)
	if err != nil {
		return nil, err
	}

	dirNames := []string{}
	for _, d := range dirs {

		if err := ex.ensureSourceDriverOpen(); err != nil {
			return nil, err
		}

		dirReaders, err := ex.sourceDriver.ReadDirectories(d, ".*")
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

func (ex *executor) resolveDirectoryWildcards(relativeDir string) ([]string, error) {
	if !strings.Contains(relativeDir, "?") {
		return []string{relativeDir}, nil
	}

	i := strings.Index(relativeDir, "?")
	dirs, err := ex.getSubdirectories(relativeDir[:i-1])
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, d := range dirs {
		base := path.Base(d)
		names, err := ex.resolveDirectoryWildcards(strings.Replace(relativeDir, "?", base, 1))
		if err != nil {
			return nil, err
		}
		result = append(result, names...)
	}

	return result, nil
}

func levelOrDryRun(ctx *Context, level log.LogLevel) log.LogLevel {
	if ctx.DryRun {
		return log.LEVEL_DRY_RUN
	}
	return level
}
