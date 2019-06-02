package repl

import (
	"github.com/c-bata/go-prompt"
	"github.com/neighborly/ddsl/exec"
)

var suggestions = []prompt.Suggest{
	// Config
	{"config", "Get or set configuration options"},
	{"get", "Get configuration options"},
	{"set", "Set configuration options"},
	{"source", "DDL source location"},
	{"database", "Database connection URL"},
/*
	// Create or Drop
	{createCmd.Use, createCmd.Short},
	{dropCmd.Use, dropCmd.Short},
	{constraintsCmd.Use, constraintsCmd.Short},
	{databaseCmd.Use, databaseCmd.Short},
	{extensionsCmd.Use, extensionsCmd.Short},
	{foreignKeysCmd.Use, foreignKeysCmd.Short},
	{indexesCmd.Use, indexesCmd.Short},
	{rolesCmd.Use, rolesCmd.Short},
	{schemaCmd.Use, schemaCmd.Short},
	{schemasCmd.Use, schemasCmd.Short},
	{tableCmd.Use, tableCmd.Short},
	{tablesCmd.Use, tablesCmd.Short},
	{typesCmd.Use, typesCmd.Short},
	{viewCmd.Use, viewCmd.Short},
	{viewsCmd.Use, viewsCmd.Short},

	// Migrate
	{bottomCmd.Use, bottomCmd.Short},
	{downCmd.Use, downCmd.Short},
	{topCmd.Use, topCmd.Short},
	{upCmd.Use, upCmd.Short},

	// Transactions
	{beginCmd.Use, beginCmd.Short},
	{commitCmd.Use, commitCmd.Short},
	{rollbackCmd.Use, rollbackCmd.Short},

	// Seed
	{seedCmd.Use, seedCmd.Short},
	{seedDatabaseCmd.Use, seedDatabaseCmd.Short},
	{seedSchemaCmd.Use, seedSchemaCmd.Short},
	{seedSqlCmd.Use, seedSqlCmd.Short},
	{seedTableCmd.Use, seedTableCmd.Short},
	{seedTablesCmd.Use, seedTablesCmd.Short},

	// SQL
	{sqlCmd.Use, sqlCmd.Short},
*/
}


func Run(ctx *exec.Context) (exitCode int, err error) {
	initializeCache(ctx)
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("ddsl repl"),
		prompt.OptionPrefix("ddsl> "),
	)
	p.Run()
	return 0, nil
}