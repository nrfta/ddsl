package parser

import (
	"bufio"
	"errors"
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
	Database       *Tag                `"DATABASE" [@@]`
	Roles          *Tag                `| "ROLES" [@@]`
	Extensions     *Tag                `| "EXTENSIONS" [@@]`
	ForeignKeys    *Tag                `| ("FOREIGN" "KEYS") [@@]`
	Schema         *Name               `| "SCHEMA" @@`
	TablesInSchema *Name               `| "TABLES" "IN" @@`
	ViewsInSchema  *Name               `| "VIEWS" "IN" @@`
	Table          *TableOrViewSubject `| "TABLE" @@`
	View           *TableOrViewSubject `| "VIEW" @@`
	Indexes        *SchemaItem         `| "INDEXES" "ON" @@`
	Constraints    *SchemaItem         `| "CONSTRAINTS" "ON" @@`
}

// TableOrViewSubject contains the schema and table or view when it is the subject of the command.
type TableOrViewSubject struct {
	TableOrView string `@Ident`
	Schema      string `"IN" @Ident`
	Tag         *Tag   `[@@]`
}

// SchemaItem contains the schema and table or view when it is the object of the command.
type SchemaItem struct {
	Item        string `@SchemaItem`
	Tag         *Tag   `[@@]`
	TableOrView string
	Schema      string
}

type Name struct {
	Name string `@Ident`
	Tag  *Tag   `[@@]`
}

// Tag contains the git tag to run the DDSL command against.
type Tag struct {
	Tag string `@Tag`
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
		// TODO: make schema optional and just have term after the dot (public schema)
		`|(?P<SchemaItem>[a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		"|(?P<Sql>(?s)`(.|\\n)*`)" +
		`|(?P<Tag>@[a-zA-Z0-9_\-\.]*)` +
		`|(?P<Int>\d*)`

	ddsllLexer = lexer.Must(lexer.Regexp(re))

	ddsllParser = participle.MustBuild(
		&DDSL{},
		participle.Lexer(ddsllLexer),
		participle.CaseInsensitive("Keyword"),
	)
)

// Parse parses an input of one or more commands and returns a slice of parse trees.
func Parse(commands string) ([]*DDSL, error) {
	scanner := bufio.NewScanner(strings.NewReader(commands))
	var trees []*DDSL
	multiline := ""
	command := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if len(multiline) > 0 {
			multiline += "\n" + line
			// We're in multiline so another backtick signals end of command
			if strings.Contains(line, "`") {
				command = multiline
				multiline = ""
			}
		} else {
			countTicks := strings.Count(line, "`")
			switch {

			// Zero or 2 backticks on a line means command is a one-liner
			case countTicks == 0 || countTicks == 2:
				command = line

			// Single backtick means enter multiline mode
			case countTicks == 1:
				multiline = line

			default:
				return nil, errors.New("syntax error: too many backticks on a line")
			}
		}

		// Once command is complete, parse it
		if len(command) > 0 {
			tree, err := parse(command)
			if err != nil {
				return nil, err
			}
			trees = append(trees, tree)
			command = ""
		}
	}

	return trees, nil
}

func parse(command string) (*DDSL, error) {
	tree := &DDSL{}
	err := ddsllParser.ParseString(command, tree)
	if err == nil {
		// Some commands require a bit of touch up
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

// Split `schema.table_or_view` in to schema and table_or_view fields
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
