package exec

import (
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/source"
	"github.com/neighborly/ddsl/parser"
)

func ExecuteTree(sourceDriver source.Driver, dbDriver database.Driver, ddsl *parser.DDSL) error {
	// TODO: execute the tree

	return nil
}
