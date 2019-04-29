package parser

import (
	"bufio"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"strings"
)

// DDSL is the top level struct of the parser. Only one of the members will be populated.
type DDSL struct {
	Create  *Command `"CREATE" @@`
	Drop    *Command `| "DROP" @@`
	Migrate *Migrate `| "MIGRATE" @@`
	Sql     string   `| "SQL" @Sql`
}

// Command contains details of a create or drop command. Only one of the members will be `true` or populated.
type Command struct {
	Database       bool                `@"DATABASE"`
	Roles          bool                `| @"ROLES"`
	Extensions     bool                `| @"EXTENSIONS"`
	ForeignKeys    bool                `| @("FOREIGN" "KEYS")`
	Schema         string              `| "SCHEMA" @Ident`
	TablesInSchema string              `| "TABLES" "IN" @Ident`
	ViewsInSchema  string              `| "VIEWS" "IN" @Ident`
	Table          *TableOrViewSubject `| "TABLE" @@`
	View           *TableOrViewSubject `| "VIEW" @@`
	Indexes        *SchemaItem         `| "INDEXES" "ON" @@`
	Constraints    *SchemaItem         `| "CONSTRAINTS" "ON" @@`
}

// TableOrViewSubject contains the schema and table or view when it is the subject of the command.
type TableOrViewSubject struct {
	TableOrView string `@Ident`
	Schema      string `"IN" @Ident`
}

// SchemaItem contains the schema and table or view when it is the object of the command.
type SchemaItem struct {
	Item        string `@SchemaItem`
	TableOrView string
	Schema      string
}

// Migrate contains the migration command. Only one of the members will be `true` or nonzero.
type Migrate struct {
	Top    bool `@"TOP"`
	Bottom bool `| @"BOTTOM"`
	Up     int  `| "UP" @Int`
	Down   int  `| "DOWN" @Int`
}

var (
	re = `(\s+)` +
		`|(?P<Keyword>(?i)CREATE|DROP|DATABASE|ROLES|EXTENSIONS|FOREIGN|KEYS|SCHEMA|TABLES|TABLE|VIEWS|VIEW|INDEXES|CONSTRAINTS|IN|ON|MIGRATE|TOP|BOTTOM|UP|DOWN|SQL)` +
		"|(?P<Sql>(?s)`(.|\\n)*`)" +
		`|(?P<SchemaItem>[a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<Int>\d*)`

	ddsllLexer = lexer.Must(lexer.Regexp(re))

	ddsllParser = participle.MustBuild(
		&DDSL{},
		participle.Lexer(ddsllLexer),
		participle.CaseInsensitive("Keyword"),
		// participle.Elide("Comment"),
		// Need to solve left recursion detection first, if possible.
		// participle.UseLookahead(),
	)
)

// Parse parses an input of one or more commands and returns a slice of parse trees.
func Parse(command string) ([]*DDSL, error) {
	scanner := bufio.NewScanner(strings.NewReader(command))
	trees := []*DDSL{}
	inMultiline := false
	multiline := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if inMultiline {
			multiline += "\n" + line
			if strings.Contains(line, "`") {
				inMultiline = false
				tree, err := parse(multiline)
				if err != nil {
					return nil, err
				}
				trees = append(trees, tree)
			}
		} else if strings.Contains(line, "`") {
			multiline = line
			inMultiline = true
		} else {
			tree, err := parse(line)
			if err != nil {
				return nil, err
			}
			trees = append(trees, tree)
		}
	}

	return trees, nil
}

func parse(command string) (*DDSL, error) {
	tree := &DDSL{}
	err := ddsllParser.ParseString(command, tree)
	if err == nil {
		if len(tree.Sql) > 0 {
			tree.Sql = strings.Trim(tree.Sql, "`")
		}
		if tree.Create != nil && tree.Create.Indexes != nil {
			tree.Create.Indexes.populate()
		}
		if tree.Create != nil && tree.Create.Constraints != nil {
			tree.Create.Constraints.populate()
		}
		if tree.Drop != nil && tree.Drop.Indexes != nil {
			tree.Drop.Indexes.populate()
		}
		if tree.Drop != nil && tree.Drop.Constraints != nil {
			tree.Drop.Constraints.populate()
		}
	}
	return tree, err

}

func (si *SchemaItem) populate() {
	if len(si.Item) > 0 {
		parts := strings.Split(si.Item, ".")
		if len(parts) == 0 {
			si.TableOrView = si.Item
		} else {
			si.Schema = parts[0]
			si.TableOrView = parts[1]
		}
	}
}
