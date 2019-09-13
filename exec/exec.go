package exec

import (
	"github.com/neighborly/ddsl/drivers/database/postgres"
	"github.com/neighborly/ddsl/drivers/source/file"
	"github.com/neighborly/ddsl/parser"
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
