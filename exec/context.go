package exec

import (
	dbdr "github.com/nrfta/ddsl/drivers/database"
	"strings"
)

type Context struct {
	SourceRepo      string
	DatbaseUrl      string
	AutoTransaction bool
	DryRun          bool
	OutputFormat    string
	inTransaction   bool
	dbDriver        dbdr.Driver
	patterns        []string
	instructions    []*instruction
	nesting         int
	nonList         bool
}

type Name struct {
	Name *string
}

func NewContext(sourceRepo, databaseURL string, autoTx, dryRun bool, output_format string) *Context {
	return &Context{
		SourceRepo:      sourceRepo,
		DatbaseUrl:      databaseURL,
		AutoTransaction: autoTx,
		DryRun:          dryRun,
		OutputFormat:    output_format,
		patterns:        []string{},
		instructions:    []*instruction{},
		nesting:         0,
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
		return "(none)"
	}
	return strings.Join(c.patterns, "\n")
}

func (c *Context) clearPatterns() {
	c.patterns = []string{}
}

func (c *Context) pushNesting() {
	c.nesting++
}

func (c *Context) popNesting() {
	if c.nesting == 0 {
		panic("level is already at 0")
	}
	c.nesting--
}

func (c *Context) getNestingForLogging() string {
	s := ""
	for l := 0; l < c.nesting; l++ {
		s += "+"
	}
	return s
}

func (c *Context) resetNesting() {
	c.nesting = 0
}

func (c *Context) clearInstructions() {
	c.instructions = []*instruction{}
}

func (c *Context) addInstructionWithParams(instrType InstructionType, params map[string]interface{}) {
	c.instructions = append(c.instructions, &instruction{instrType, params})
	if instrType != INSTR_LIST && instrType != INSTR_DDSL {
		c.nonList = true
	}
}

func (c *Context) addInstruction(instrType InstructionType) {
	c.addInstructionWithParams(instrType, make(map[string]interface{}))
}

func (c *Context) isListCommand() bool {
	return !c.nonList
}
