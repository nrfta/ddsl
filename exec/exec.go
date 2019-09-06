package exec

import (
	"errors"
	"github.com/neighborly/ddsl/drivers/database/postgres"
	"github.com/neighborly/ddsl/drivers/source"
	"github.com/neighborly/ddsl/drivers/source/file"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"strings"
)

func init() {
	postgres.Register()
	file.Register()
}

type executor struct {
	ctx           *Context
	sourceDriver  source.Driver
	databaseName  string
	command       *parser.Command
	createOrDrop  string
	grantOrRevoke string
}

func ExecuteBatch(ctx *Context, cmds []*parser.Command) error {
	_, err := preprocessBatch(ctx, cmds)
	if err != nil {
		return err
	}

	p := &processor{ctx}
	return p.process()
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

func (ex *executor) resolveDirectoryWildcards(relativeDir string) ([]string, error) {
	return nil, nil
}

func levelOrDryRun(ctx *Context, level log.LogLevel) log.LogLevel {
	if ctx.DryRun {
		return log.LEVEL_DRY_RUN
	}
	return level
}
