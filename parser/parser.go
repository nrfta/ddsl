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
	Sql     *string  `| "SQL" @Sql`
}

// Command contains details of a create or drop command. Only one of the members will be `true` or populated.
type Command struct {
	Database       *Database    `"DATABASE" @@`
	Roles          *Roles       `| "ROLES" @@`
	Extensions     *Extensions  `| "EXTENSIONS" @@`
	ForeignKeys    *ForeignKeys `| ("FOREIGN" "KEYS") @@`
	Schema         *Name        `| "SCHEMA" @@`
	TablesInSchema *Name        `| "TABLES" "IN" @@`
	ViewsInSchema  *Name        `| "VIEWS" "IN" @@`
	Table          *SchemaItem  `| "TABLE" @@`
	View           *SchemaItem  `| "VIEW" @@`
	Indexes        *SchemaItem  `| "INDEXES" "ON" @@`
	Constraints    *SchemaItem  `| "CONSTRAINTS" "ON" @@`
}

// Database contains details for action on a database.
type Database struct {
	Ref *Ref `[@@]`
}

// Roles contains details for action on roles.
type Roles struct {
	Ref *Ref `[@@]`
}

// Extensions contains details for action on extensions.
type Extensions struct {
	Ref *Ref `[@@]`
}

// ForeignKeys contains details for action on a foreign keys.
type ForeignKeys struct {
	Ref *Ref `[@@]`
}

// SchemaItem contains the schema and table or view when it is the object of the command.
type SchemaItem struct {
	Item        string `@SchemaItem`
	Ref         *Ref   `[@@]`
	TableOrView string
	Schema      string
}

type Name struct {
	Name string `@Ident`
	Ref  *Ref   `[@@]`
}

// Tag contains the git tag to run the DDSL command against.
type Ref struct {
	Ref string `@Ref`
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
		`|(?P<Ref>@[a-zA-Z0-9_\-\.\/]*)` +
		`|(?P<Comment>--.*)` +
		`|(?P<MultiComment>(?s)/\*(.*|\n)\*/)` +
		`|(?P<Int>\d*)`

	ddslEbnf = `
		Keyword = ( "CREATE" | "DROP" | "DATABASE" | "ROLES" | "EXTENSIONS" | "FOREIGN" | "KEYS" | "SCHEMA" | "TABLES" | "TABLE" | "VIEWS" | "VIEW" | "INDEXES" | "CONSTRAINTS" | "IN" | "ON" | "MIGRATE" | "TOP" | "BOTTOM" | "UP" | "DOWN" | "SQL" ) .
	    Comment = "--" {  any_no_newline } .
		MultiComment = "/*" { any } "*/" .
		Ident = (alpha | "_") { "_" | alpha | digit } .
		SchemaItem = Ident [ "." Ident ] .
		Ref = "@" ( alpha | digit ) { alpha | digit | "_" | "." | "-" | "/" } .
		Int = { digit } .
		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		punct = "!"…"/" | ":"…"@" | "["…"_" | "{"…"~" .
		any = "\u0000"…"\uffff" .
		any_no_newline = ( "\u0000"…"\u0009" | "\u000e"…"\uffff" ) .
		newline = ( "\n" | "\r" ) . ` +
		"Sql = \"`\" { alpha | digit | punct } \"`\" ."

	//ddslLexer = lexer.Must(ebnf.New(ddslEbnf))

	ddslLexer = lexer.Must(lexer.Regexp(re))

	ddslParser = participle.MustBuild(
		&DDSL{},
		participle.Lexer(ddslLexer),
		participle.CaseInsensitive("Keyword"),
		participle.Elide("Comment", "MultiComment"),
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
	err := ddslParser.ParseString(command, tree)
	if err == nil {
		// Some commands require a bit of touch up
		if tree.Sql != nil {
			sql := *tree.Sql
			trimmedSql := strings.Trim(sql, "`")
			tree.Sql = &trimmedSql
		}
		if tree.Create != nil {
			if tree.Create.Table != nil {
				tree.Create.Table.populate()
			}
			if tree.Create.View != nil {
				tree.Create.View.populate()
			}
			if tree.Create.Indexes != nil {
				tree.Create.Indexes.populate()
			}
			if tree.Create.Constraints != nil {
				tree.Create.Constraints.populate()
			}
		}
		if tree.Drop != nil {
			if tree.Drop.Table != nil {
				tree.Drop.Table.populate()
			}
			if tree.Drop.View != nil {
				tree.Drop.View.populate()
			}
			if tree.Drop.Indexes != nil {
				tree.Drop.Indexes.populate()
			}
			if tree.Drop.Constraints != nil {
				tree.Drop.Constraints.populate()
			}
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
