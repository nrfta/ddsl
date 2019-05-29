package exec

import (
	"errors"
	"fmt"
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"github.com/neighborly/ddsl/drivers/database/postgres"
	"github.com/neighborly/ddsl/drivers/source"
	"github.com/neighborly/ddsl/drivers/source/file"
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
	repo         string
	sourceDriver source.Driver
	dbDriver     dbdr.Driver
	dbURL        string
	databaseName string
	parseTree    *parser.DDSL
	createOrDrop string
}

const (
	create string = "create"
	drop   string = "drop"
)

func Execute(repo string, dbURL string, command string) error {
	trees, err := parser.Parse(command)
	if err != nil {
		return err
	}

	dbDriver, err := dbdr.Open(dbURL)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	// database commands cannot run in transaction
	cmds := getDatabaseCommands(trees)
	if len(cmds) > 0 {
		fmt.Println("[INFO] running database commands, transaction not possible")
	}
	for _, t := range cmds {
		ex := &executor{
			repo:      repo,
			dbDriver:  dbDriver,
			dbURL:     dbURL,
			parseTree: t,
		}
		if err := execute(ex); err != nil {
			return err
		}
	}

	cmds = []*parser.DDSL{}

	for _, t := range trees {
		if isDatabaseCommand(t) {
			continue
		}
		cmds = append(cmds, t)
	}

	if len(cmds) == 0 {
		return nil
	}

	dryRun := viper.GetBool("dry_run")

	if dryRun {
		fmt.Println("[DRY-RUN] would begin transaction")
	} else {
		fmt.Println("[INFO] beginning transaction")
		err = dbDriver.Begin()
		if err != nil {
			return err
		}
	}

	for _, t := range cmds {

		ex := &executor{
			repo:      repo,
			dbDriver:  dbDriver,
			parseTree: t,
		}
		if err := execute(ex); err != nil {
			if dryRun {
				fmt.Println("[DRY-RUN] would rollback transaction")
				return err
			}
			fmt.Println("[WARN] rolling back transaction")
			_ = dbDriver.Rollback()
			return err
		}
	}

	if dryRun {
		fmt.Println("[DRY-RUN] would commit transaction")
		return nil
	}

	fmt.Println("[INFO] committing transaction")
	return dbDriver.Commit()
}

func getDatabaseCommands(trees []*parser.DDSL) []*parser.DDSL {
	cmds := []*parser.DDSL{}

	for _, t := range trees {
		if isDatabaseCommand(t) {
			cmds = append(cmds, t)
		}
	}
	return cmds
}

func isDatabaseCommand(t *parser.DDSL) bool {
	return (t.Create != nil && t.Create.Database != nil) || (t.Drop != nil && t.Drop.Database != nil)
}

func execute(ex *executor) error {
	switch {
	case ex.parseTree.Create != nil:
		ex.createOrDrop = create
		return executeCreateOrDrop(ex)
	case ex.parseTree.Drop != nil:
		ex.createOrDrop = drop
		return executeCreateOrDrop(ex)
	case ex.parseTree.Migrate != nil:
		return executeMigrate(ex, ex.parseTree.Migrate)
	case ex.parseTree.Sql != nil:
		if viper.GetBool("dry_run") {
			fmt.Println("[DRY-RUN] would execute SQL statement")
			return nil
		}
		fmt.Println("[INFO] executing SQL statement")
		return ex.dbDriver.Exec(strings.NewReader(*ex.parseTree.Sql))
	}

	return errors.New("unknown command")
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

func (ex *executor) execute(pathPattern string, ref *parser.Ref) error {
	if err := ex.getSourceDriver(ref); err != nil {
		return err
	}
	defer ex.sourceDriver.Close()

	relativePath, filePattern := getRelativePathAndFilePattern(pathPattern)
	readers, err := ex.sourceDriver.ReadFiles(relativePath, filePattern)
	if err != nil {
		return err
	}

	dryRun := viper.GetBool("dry_run")

	for _, fr := range readers {
		if dryRun {
			fmt.Printf("[DRY-RUN] would execute %s\n", fr.FilePath)
			continue
		}
		fmt.Printf("[INFO] executing %s\n", fr.FilePath)
		err = ex.dbDriver.Exec(fr.Reader)
		if err != nil {
			return err
		}
	}

	return nil
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
