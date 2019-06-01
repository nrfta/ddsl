package parser

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"strings"
)

type Tree struct {
	CommandDefs map[string]*CommandDef
}

type Arg struct {
	Name      string
	ShortDesc string
}

type CommandDef struct {
	Name        string
	ShortDesc   string
	Props       map[string]string
	Level       int
	Parent      *CommandDef
	CommandDefs map[string]*CommandDef
	ArgDefs     map[string]*Arg
}

type Command struct {
	CommandDef *CommandDef
	RootDef *CommandDef
	Clause string
	Args []string
	ExtArgs []string
	Ref *string
	tokenIndex int
}

var ParseTree *Tree
var shellParser *shellwords.Parser

var commandSpec = `ddsl,Top level command,primary
  create,Top level create command,root
    database,Create or drop the database,primary
    roles,Create or drop roles,primary
    foreign-keys,Create or drop foreign keys,primary
    schemas,Create or drop all schemas,primary
      except,Comma-delimited list of schemas to exclude,optional
        -exclude_schemas,Comma-delimited list of schemas to exclude
    schema,Create or drop one or more schemas,primary
      -include_schemas,Comma-delimited list of schemas
    extensions,Create or drop all extensions in one or more schemas,primary
      in,Comma delimited list of schemas,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    tables,Create or drop all tables in one or more schemas,primary
      in,Comma delimited list of schemas,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    views,Create or drop all views in one or more schemas,primary
      in,Comma delimited list of schemas,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    table,Create or drop one or more tables,primary
      -include_tables,Comma-delimited list of tables
    view,Create or drop one or more views,primary
      -include_views,Comma-delimited list of views
    indexes,Create or drop all indexes on one or more tables or views,primary
      on,Comma delimited list of tables and views
        -include_tables_and_views,Comma-delimited list of tables and views
    constraints,Create or drop all constraints on one or more tables,primary
      on,Comma delimited list of tables
        -include_tables,Comma-delimited list of tables
    types,Create or drop all types in one or more schemas,primary
      in,Create or drop types,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    type,Create or drop one or more types,primary
      -include_types,Comma-delimited list of types
  migrate,Top level migrate command,root
    up,Migrate the database up in version,primary
      -number_of_versions,Number of versions to migrate
    down,Migrate the database down in version,primary
      -number_of_versions,Number of versions to migrate
    top,Migrate the database to the latest version,primary
    bottom,Migrate the database to the earliest version,primary
  seed,Top level seed command,root
    cmd,Seed the database by running a shell command,primary
      -command,Shell command to run
    database,Seed the database,primary
      with,Seed the database,optional
        -database_seeds,Comma-delimited list of seeds
      without,Seed the database,optional
        -database_seeds,Comma-delimited list of seeds
    schemas,Seed all schemas,primary
      except,Comma-delimited list of schemas to exclude,optional
        -exclude_schemas,Comma-delimited list of schemas
    schema,Seed a schema,primary,ext-args
      -single_schema,Single schema name
      with,Seed a schema,optional
        -schema_seeds,Comma-delimited list of seeds
      without,Seed a schema,optional
        -schema_seeds,Comma-delimited list of seeds
    table,Seed one or more tables,primary
      -include_tables,Comma-delimited list of tables
    tables,Seed all tables in given schemas,primary
      in,Comma delimited list of schemas,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    sql,Seed with SQL command of script,primary
      -command,SQL command to run
      -file,SQL file to seed database
    from,Seed with CSV file,primary
      -file,CSV file to seed database
  sql,Run an SQL command or script,primary
    -command,SQL command to run
    -file,SQL file to run
  grant,Top-level grant command,root
    privileges,Top-level grant or revoke privileges command,optional,non-exec
      on,Top-level grant or revoke privileges command,non-exec
        database,Grant or revoke privileges on the database,primary
        schemas,Grant or revoke privileges on all schemas,primary
          except,Comma-delimited list of schemas to exclude,optional
            -exclude_schemas,Comma-delimited list of schemas
        schema,Grant or revoke privileges on one or more schemas,primary
          -include_schemas,Comma-delimited list of schemas
        tables,Grant or revoke privileges on all tables in one or more schemas,primary
          in,Comma delimited list of schemas,optional
            -include_schemas,Comma-delimited list of schemas
          except,Comma-delimited list of schemas to exclude,optional
            in,Comma delimited list of schemas
              -exclude_schemas,Comma-delimited list of schemas
        views,Grant or revoke privileges on all views in one or more schemas,primary
          in,Comma delimited list of schemas,optional
            -include_schemas,Comma-delimited list of schemas
          except,Comma-delimited list of schemas to exclude,optional
            in,Comma delimited list of schemas
              -exclude_schemas,Comma-delimited list of schemas
        table,Grant or revoke privileges on one or more tables,primary
          -include_tables,Comma-delimited list of tables
        view,Grant or revoke privileges on one or more views,primary
          -include_views,Comma-delimited list of views
  begin,Begin a transaction,root,primary
    transaction,Begin a transaction,optional
  commit,Commit the current transaction,root,primary
    transaction,Commit the current transaction,optional
  rollback,Rollback the current transaction,root,primary
    transaction,Rollback the current transaction,optional`

func init() {
	initialize()
}

// TryParse parses the given partial command and returns the deepest associated `Command`.
// This is used for repl and commandline completions.
func TryParse(command string) (cmd *Command, remainder []string, err error) {
	if len(command) == 0 {
		return nil, nil, fmt.Errorf("no command was provided")
	}

	tokens, err  := shellParser.Parse(command)
	if err != nil {
		return nil, nil, err
	}

	cmdDefs:= ParseTree.CommandDefs
	args := []string{}
	remainder = []string{}
	err = nil
	var cmdDef *CommandDef
	for i, token := range tokens {
		next, ok := cmdDefs[strings.ToLower(token)]
		if ok {
			tokenIndex := i
			if next.HasExtArgs() {
				for a := 0; a < len(next.ArgDefs); a++ {
					if tokenIndex+1 < len(tokens) {
						tokenIndex++
						args = append(args, tokens[tokenIndex])
					} else {
						break
					}
				}
			}
			if len(tokens[tokenIndex:]) > 1 {
				remainder = tokens[tokenIndex+1:]
			} else {
				remainder = []string{}
			}
			if next.IsPrimary() {
				cmd = makeCommand(next, args, tokenIndex-i)
				return
			}
			cmdDef = next
			cmdDefs = next.CommandDefs
		} else {
			if len(cmdDef.ArgDefs) > 0 {
				// token is not a command, so assume it's an arg,
				// do not advance down the parse tree
				args = append(args, token)
			} else {
				next, _ = cmdDef.skipOptionalTo(token)
				if next == nil {
					err = fmt.Errorf("syntax error at '%s'", token)
					return
				}
				// advance down the parse tree
				cmdDef = next
				cmdDefs = next. CommandDefs
			}
		}
	}

	err = fmt.Errorf("syntax error in '%s'", command)
	return
}

func Parse(command string) (*Command, error) {
	cmd, remainder, err := TryParse(command)
	if !cmd.CommandDef.IsPrimary() {
		return nil, fmt.Errorf("primary command token not found")
	}

	if err != nil {
		return cmd, err
	}

	clause, extArgs, err := cmd.parseRemainder(remainder)
	if err != nil {
		return cmd, err
	}
	cmd.Clause = clause
	cmd.ExtArgs = extArgs

	return cmd, err
}

func (c *Command) parseRemainder(tokens []string) (clause string, extArgs []string, err error) {
	clause = ""
	extArgs = []string{}
	err = nil
	if len(tokens) == 0 {
		return
	}

	clauseSl := []string{}
	cmdDef := c.CommandDef
	for _, token := range tokens {
		token = strings.ToLower(token)
		next, ok := cmdDef.CommandDefs[token]
		if ok {
			clauseSl = append(clauseSl, token)
		} else {
			if len(cmdDef.ArgDefs) > 0 {
				// assume the rest is args
				clause = strings.Join(clauseSl, " ")
				extArgs = strings.Split(token,",")
				return
			}
			var skipped []string
			next, skipped = cmdDef.skipOptionalTo(token)
			if next == nil {
				err = fmt.Errorf("syntax error at '%s'", token)
				return
			}
			clauseSl = append(clauseSl, skipped...)
		}
		cmdDef = next
	}
	clause = strings.Join(clauseSl, " ")
	return
}

func (c *CommandDef) skipOptionalTo(token string) (*CommandDef, []string) {
	if len(c.CommandDefs) == 0 {
		return nil, []string{}
	}
	return c._skipOptionalToWork(token, []string{})
}

func (c *CommandDef) _skipOptionalToWork(token string, skipped []string) (*CommandDef, []string) {
	for _, next := range c.CommandDefs {
		if next.Name == strings.ToLower(token) {
			return next, skipped
		}
		if next.IsOptional() {
			skipped = append(skipped, next.Name)
			return next._skipOptionalToWork(token, skipped)
		}
	}
	return nil, skipped
}

// ShortDesc returns the `ShortDesc` field of a command. ShortDesc panics
// if the command is zero length or contains an unrecognized command.
func ShortDesc(command string) string {
	cmd, err := Parse(command)
	if err != nil {
		panic(err)
	}

	return cmd.CommandDef.ShortDesc
}

func (c *CommandDef) ParentAtLevel(level int) *CommandDef {
	if c.Level == level {
		return c
	}

	if c.Level > level {
		p := c
		for p.Level > level {
			p = p.Parent
		}

		return p
	}

	panic("level must be lower than current command")
}

func (c *CommandDef) IsOptional() bool {
	return c.hasProp("optional")
}
func (c *CommandDef) IsRoot() bool {
	return c.hasProp("root")
}
func (c *CommandDef) IsNonExec() bool {
	return c.hasProp("non-exec")
}
func (c *CommandDef) IsPrimary() bool {
	return c.hasProp("primary")
}
func (c *CommandDef) HasExtArgs() bool {
	return c.hasProp("ext-args")
}

func (c *CommandDef) hasProp(name string) bool {
	_, ok := c.Props[name]
	return ok
}

func initialize() {
	shellParser = shellwords.NewParser()
	shellParser.ParseEnv = true
	shellParser.ParseBacktick = true

	scanner := bufio.NewScanner(strings.NewReader(commandSpec))
	levels := map[int]*CommandDef{}

	scanner.Scan()
	line := scanner.Text()
	ddsl := procCommand(line, nil)
	levels[0] = ddsl

	for scanner.Scan() {
		line = scanner.Text()
		level := indentLevel(line)
		parentCmd := levels[level-1]
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			arg := procArg(line)
			parentCmd.ArgDefs[arg.Name] = arg
		} else {
			subCmd := procCommand(line, parentCmd)
			parentCmd.CommandDefs[subCmd.Name] = subCmd
			levels[level] = subCmd
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	ddsl.CommandDefs["drop"] = &CommandDef{
		Name:        "drop",
		ShortDesc:   "Top level drop command",
		Level:       1,
		Props:       map[string]string{"non-exec":"true"},
		Parent:      ddsl,
		CommandDefs: ddsl.CommandDefs["create"].CommandDefs,
		ArgDefs:     ddsl.CommandDefs["create"].ArgDefs,
	}

	ddsl.CommandDefs["revoke"] = &CommandDef{
		Name:        "revoke",
		ShortDesc:   "Top level revoke command",
		Level:       1,
		Props:       map[string]string{"non-exec":"true"},
		Parent:      ddsl,
		CommandDefs: ddsl.CommandDefs["grant"].CommandDefs,
		ArgDefs:     ddsl.CommandDefs["grant"].ArgDefs,
	}

	ParseTree = &Tree{ddsl.CommandDefs}
}

func procCommand(line string, parentCmd *CommandDef) *CommandDef {
	items := strings.Split(line, ",")

	level := 0
	if parentCmd != nil {
		level = parentCmd.Level + 1
	}

	cmdDef := &CommandDef{
		Name:        items[0],
		ShortDesc:   items[1],
		Props:       map[string]string{},
		Level:       level,
		Parent:      parentCmd,
		CommandDefs: map[string]*CommandDef{},
		ArgDefs:     map[string]*Arg{},
	}

	for i := 2; i < len(items); i++ {
		p := strings.Split(items[i], "=")
		name := p[0]
		value := "true"
		if len(p) == 2 {
			value = p[1]
		}
		cmdDef.Props[name] = value
	}

	return cmdDef
}

func procArg(line string) *Arg {
	items := strings.Split(line, ",")

	return &Arg{
		Name:      items[0],
		ShortDesc: items[1],
	}
}

func indentLevel(line string) int {
	indent := 0
	for strings.HasPrefix(line, "  ") {
		indent++
		line = strings.Replace(line, "  ", "", 1)
	}

	return indent
}

func makeCommand(cmdDef *CommandDef, args []string, tokenIndex int) *Command {
	lastArg := ""
	if len(args) > 0 {
		lastArg = args[len(args)-1]
	}

	var ref *string
	if len(lastArg) > 0 && strings.HasPrefix(lastArg, "@") {
		r := lastArg[1:]
		ref = &r
		args = args[:len(args)-1]
	}
	return &Command{
		CommandDef: cmdDef,
		RootDef: cmdDef.getRoot(),
		Args: args,
		Ref: ref,
		tokenIndex: tokenIndex,
	}
}

func (c *CommandDef) getRoot() *CommandDef {
	if c.IsRoot() {
		return c
	}
	return c.Parent.getRoot()
}