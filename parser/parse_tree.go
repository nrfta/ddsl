package parser

import (
	"bufio"
	"strings"
)

type tree struct {
	CommandDefs map[string]*CommandDef
}

var ParseTree *tree

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
    functions,Create or drop all functions in one or more schemas,primary
      in,Comma delimited list of schemas,optional
        -include_schemas,Comma-delimited list of schemas
      except,Comma-delimited list of schemas to exclude,optional
        in,Comma delimited list of schemas
          -exclude_schemas,Comma-delimited list of schemas
    procedures,Create or drop all procedures in one or more schemas,primary
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
    triggers,Create or drop all triggers on one or more tables,primary
      on,Comma delimited list of tables
        -include_tables,Comma-delimited list of tables
    function,Create or drop one or more functions,primary
      -include_functions,Comma-delimited list of functions
    procedure,Create or drop one or more procedures,primary
      -include_procedures,Comma-delimited list of procedures
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
        -database_seeds,Comma-delimited list of named seeds
      without,Seed the database,optional
        -database_seeds,Comma-delimited list of named seeds
    schemas,Seed all schemas,primary
      except,Comma-delimited list of schemas to exclude,optional
        -exclude_schemas,Comma-delimited list of schemas
    schema,Seed a schema,primary,ext-args
      -single_schema,Single schema name
      with,Seed a schema,optional
        -schema_seeds,Comma-delimited list of named seeds
      without,Seed a schema,optional
        -schema_seeds,Comma-delimited list of named seeds
    table,Seed one or more tables,primary,ext-args
      -include_tables,Comma-delimited list of tables
      with,Seed a table,optional
        -table_seeds,Comma-delimited list of named seeds
      without,Seed a table,optional
        -schema_seeds,Comma-delimited list of named seeds
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
  sql,Run a SQL command or script,root,primary
    -command,SQL command to run
    -file,SQL file to run
  grant,Top level grant command,root
    privileges,Top level grant or revoke privileges command,optional,non-exec
      on,Top level grant or revoke privileges command,non-exec
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

	dropDef := ddsl.CommandDefs["create"].clone(ddsl)
	dropDef.Name = "drop"
	dropDef.ShortDesc = "Top level drop command"
	ddsl.CommandDefs["drop"] = dropDef

	revokeDef := ddsl.CommandDefs["grant"].clone(ddsl)
	revokeDef.Name = "revoke"
	revokeDef.ShortDesc = "Top level revoke command"
	ddsl.CommandDefs["revoke"] = revokeDef

	ParseTree = &tree{ddsl.CommandDefs}
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
		ArgDefs:     map[string]*ArgDef{},
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

func procArg(line string) *ArgDef {
	items := strings.Split(line, ",")

	return &ArgDef{
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

func (c *CommandDef) clone(parent *CommandDef) *CommandDef {
	cl := &CommandDef{
		Name:        c.Name,
		ShortDesc:   c.ShortDesc,
		Props:       map[string]string{},
		Level:       c.Level,
		Parent:      parent,
		CommandDefs: map[string]*CommandDef{},
		ArgDefs:     map[string]*ArgDef{},
	}
	for _, cd := range c.CommandDefs {
		cl.CommandDefs[cd.Name] = cd.clone(cl)
	}
	for _, ad := range c.ArgDefs {
		cl.ArgDefs[ad.Name] = &ArgDef{
			Name:      ad.Name,
			ShortDesc: ad.ShortDesc,
		}
	}
	for name, value := range c.Props {
		cl.Props[name] = value
	}
	return cl
}
