package exec

import (
	"github.com/nrfta/ddsl/drivers/database/postgres"
	"github.com/nrfta/ddsl/drivers/source/file"
	"github.com/nrfta/ddsl/parser"
)

func init() {
	postgres.Register()
	file.Register()
}

func ExecuteBatch(ctx *Context, cmds []*parser.Command) error {
	_, err := preprocessBatch(ctx, cmds)
	if err != nil {
		return err
	}

	p := &processor{ctx}
	return p.process()
}
