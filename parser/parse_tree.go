package parser

import (
	"bufio"
	"fmt"
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
	Optional    bool
	Primary     bool
	NonExec     bool
	Level       int
	Parent      *CommandDef
	CommandDefs map[string]*CommandDef
	ArgDefs     map[string]*Arg
}

type Command struct {
	CommandDef *CommandDef
	Args []string
}

var ParseTree *Tree

var commandSpec = `ddsl,Top level command,primary
  create,Top level create command,non-exec
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
        -include_tables,Comma-delimited list of tables
        -include_views,Comma-delimited list of views
    constraints,Create or drop all constraints on one or more talbes,primary
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
  migrate,Top level migrate command,non-exec
    up,Migrate the database up in version,primary
      -number_of_versions,Number of versions to migrate
    down,Migrate the database down in version,primary
      -number_of_versions,Number of versions to migrate
    top,Migrate the database to the latest version,primary
    bottom,Migrate the database to the earliest version,primary
  seed,Top level seed command,non-exec
    cmd,Seed the database by running a shell command,primary
      -command,Shell command to run
    database,Seed the database,primary
      with,Seed the database,optional
        -database_seeds,Comma-delimited list of seeds
      without,Seed the database,optional
        -database_seeds,Comma-delimited list of seeds
    schema,Seed a schema,primary
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
        in,Comma delimited list of schemas,optional
          -exclude_schemas,Comma-delimited list of schemas
    sql,Seed with SQL command of script,primary
      -command,SQL command to run
      -file,SQL file to run
  sql,Run an SQL command or script,primary
    -command,SQL command to run
    -file,SQL file to run
  grant,Top-level grant command,non-exec
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
  begin,Begin a transaction,primary
    transaction,Begin a transaction,optional
  commit,Commit the current transactionm,primary
    transaction,Commit the current transaction,optional
  rollback,Rollback the current transaction,primary
    transaction,Rollback the current transaction,optional`

func init() {
	initialize()
}

// TryParse parses the given partial command and returns the deepest associated `Command`.
// This is used for repl and commandline completions.
func TryParse(command string) (*Command, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("no command was provided")
	}

	keys := strings.Split(command, " ")
	cmdDefs:= ParseTree.CommandDefs
	var cmdDef *CommandDef
	for i,k := range keys {
		var ok bool
		_, ok = cmdDefs[strings.ToLower(k)]
		if !ok {
			if len(cmdDef.ArgDefs) > 0 {
				cmd := &Command{
					CommandDef: cmdDef,
					Args: keys[i:],
				}
				return cmd, nil
			}
			return nil, fmt.Errorf("syntax error at '%s'", k)
		}
		cmdDef = cmdDefs[strings.ToLower(k)]
		cmdDefs = cmdDef.CommandDefs
		if cmdDef.Primary {
			cmd := &Command{
				CommandDef: cmdDef,
				Args: keys[i:],
			}
			return cmd, nil
		}
	}

	cmd := &Command{
		CommandDef: cmdDef,
		Args: []string{},
	}
	return cmd, nil
}

func Parse(command string) (*Command, error) {
	cmd, err := TryParse(command)
	if !cmd.CommandDef.Primary {
		return nil, fmt.Errorf("primary command token not found")
	}

	return cmd, err
}

// ShortDesc returns the `ShortDesc` field of a command. ShortDesc panics
// if the command is zero length or contains an unrecognized command.
func ShortDesc(command string) string {
	cmdDef, err := Parse(command)
	if err != nil {
		panic(err)
	}

	return cmdDef.ShortDesc
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

func initialize() {
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
		NonExec:     true,
		Parent:      ddsl,
		CommandDefs: ddsl.CommandDefs["create"].CommandDefs,
		ArgDefs:     ddsl.CommandDefs["create"].ArgDefs,
	}

	ddsl.CommandDefs["revoke"] = &CommandDef{
		Name:        "revoke",
		ShortDesc:   "Top level revoke command",
		Level:       1,
		NonExec:     true,
		Parent:      ddsl,
		CommandDefs: ddsl.CommandDefs["grant"].CommandDefs,
		ArgDefs:     ddsl.CommandDefs["grant"].ArgDefs,
	}

	ParseTree = &Tree{ddsl.CommandDefs}
}

func procCommand(line string, parentCmd *CommandDef) *CommandDef {
	items := strings.Split(line, ",")

	cmdDef := &CommandDef{
		Name:        items[0],
		ShortDesc:   items[1],
		Level:       parentCmd.Level + 1,
		Parent:      parentCmd,
		CommandDefs: map[string]*CommandDef{},
		ArgDefs:     map[string]*Arg{},
	}

	for i := 2; i < len(items); i++ {
		cmdDef.Optional = items[i] == "optional"
		cmdDef.Primary = items[i] == "primary"
		cmdDef.NonExec = items[i] == "non-exec"
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
