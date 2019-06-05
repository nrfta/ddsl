package exec

import (
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"strings"
)

type Context struct {
	SourceRepo      string
	DatbaseUrl      string
	AutoTransaction bool
	DryRun          bool
	inTransaction   bool
	dbDriver        dbdr.Driver
	patterns        []string
}

type Name struct {
	Name *string
}

func NewContext(sourceRepo, databaseURL string, autoTx, dryRun bool) *Context {
	return &Context{
		SourceRepo:      sourceRepo,
		DatbaseUrl:      databaseURL,
		AutoTransaction: autoTx,
		DryRun:          dryRun,
		patterns:        []string{},
	}
}

func (c *Context) GetDatabaseSchemas() ([]string, error) {
	return c.getNames("SELECT schema_name AS name FROM information_schema.schemata;")
}
func (c *Context) GetDatabaseTables() ([]string, error) {
	return c.getNames(`SELECT schema_name + '.' + table_name AS name
                             FROM information_schema.tables
                             WHERE table_type = 'BASE TABLE'`)
}
func (c *Context) GetDatabaseViews() ([]string, error) {
	return c.getNames(`SELECT schema_name + '.' + table_name AS name
                             FROM information_schema.tables
                             WHERE table_type = 'VIEW'`)
}
func (c *Context) GetDatabaseTypes() ([]string, error) {
	return c.getNames(`SELECT udt_schema + '.' + udt.name AS name
                             FROM information_schema.columns
                             WHERE data_type = 'USER-DEFINED'`)
}

func (c *Context) getNames(query string) ([]string, error) {
	dbdr, err := c.dbDriver.Open(c.DatbaseUrl)
	if err != nil {
		return nil, err
	}

	rows, err := dbdr.Query(strings.NewReader(query))
	if err != nil {
		return nil, err
	}

	names := []string{}
	for rows.Next() {
		n := &Name{}
		err := rows.Scan(n)
		if err != nil {
			return nil, err
		}
		names = append(names, *n.Name)
	}

	return names, nil
}

func (c *Context) addPattern(pattern string) {
	c.patterns = append(c.patterns, pattern)
}

func (c *Context) getPatterns() string {
	if len(c.patterns) == 0 {
		return "none"
	}
	return strings.Join(c.patterns, "\n")
}