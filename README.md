# ddsl

Data Definition Support Language (DDSL, pronounced "diesel") provides a scripting language for DDL, migrations, and seeds. 

**Why a new language**

This project was born of the philosophy of "Database as Code". Database structures should be defined as code, versioned 
as code, released as code, and maintained as code. DDSL is DevOps for Database.

In addition, database code should be delivered as readily into unit tests, integration tests, and development environments
as it is into a production environment.

There are a few DevOps for Database tools in circulation already. What makes Diesel different is its full, unabashed
commitment to SQL. Database code should be stored and revisioned in its native language. Other DevOps for Database
tools require database structures to be defined in other languages such as YAML or XML.

Diesel helps you store your database DDL, migrations, and seeds in one revision control repository. It is opinionated
about the structure of the repository. That allows a set of simple commands to apply DDL structures, migrations, and seeds
agnostic of the database structure or RDS in use.

## Install

```$sh
go install github.com/nrfta/ddsl
```

## Usage

```$sh
# execute commands directly
`ddsl -s <source_repo> -d <database_rds_url> COMMAND`

# execute commands as a script
`ddsl -s <source_repo> -d <database_rds_url> -f /path/to/file.ddsl`
```

Alternatively, place a shebang at the top of the script file and make it executable. 
Set the environment variables before executing as described below.
```
#!/path/to/ddsl
COMMAND
COMMAND
COMMAND
```

The usage can be shortend by setting environment variables.

* `DDSL_SOURCE` - Source code repo URL for the database DDL and migrations
* `DDSL_DATABASE` - Database URL in format expected by RDS, properly URL encoded

The `--dry-run` switch will present what a command or script would do without making any changes.

## Command Syntax

Commands are not case sensitive, through database objects usually are. Commands may be separated by a semicolon and/or a newline. The semicolon is not required when executing a single command.

### Databases

Databases cannot be created within a transaction on certain RDSs such as Postgres. When creating a database from scratch,
the recommended order of operations is:

1. `create roles;`
2. `create database;` 
3. `create schemas; create extensions; create tables; create views; etc..`    

### CREATE and DROP
```
create database;
create roles;
create extensions;
create schemas;
create schema <schema_name>[,<schema_name> ...];
create tables [ (in | except in) <schema_name>[,<schema_name> ...] ];
create foreign-keys [ (in | except in) <schema_name>[,<schema_name> ...] ];
create views [ (in | except in) <schema_name>[,<schema_name> ...] ];
create types [ (in | except in) <schema_name>[,<schema_name> ...] ];
create functions [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];
create procedures [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];
create table <schema_name.table_name>[,<schema_name.table_name> ...];
create view <schema_name.view_name>[,<schema_name.view_name> ...];
create type <schema_name.type_name>[,<schema_name.type_name> ...];
create function <schema_name.function_name>[,<schema_name.function_name> ...];
create procedure <schema_name.procedure_name>[,<schema_name.procedure_name> ...];
create constraints on <schema_name.table_name>[,<schema_name.table_name> ...];
create indexes on <schema_name.table_or_view_name>[,<schema_name.table_or_view_name> ...];
create triggers [on] <schema_name.table_name>[,<schema_name.table_name> ...];
```

`drop` syntax is the same as `create`.


### LIST
List objects from the database. This command ignores the source files.
```
list schemas;
list tables [ (in | except in) <schema_name>[,<schema_name> ...] ];
list foreign-keys [ (in | except in) <schema_name>[,<schema_name> ...] ];
list views [ (in | except in) <schema_name>[,<schema_name> ...] ];
list types [ (in | except in) <schema_name>[,<schema_name> ...] ];
list functions [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];
list procedures [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];
list table <schema_name.table_name>[,<schema_name.table_name> ...];
list view <schema_name.view_name>[,<schema_name.view_name> ...];
list type <schema_name.type_name>[,<schema_name.type_name> ...];
list function <schema_name.function_name>[,<schema_name.function_name> ...];
list procedure <schema_name.procedure_name>[,<schema_name.procedure_name> ...];
```

### SQL
```
sql `
    UPDATE foo.bar SET field1 = 4 WHERE field2 = 0;
    DELETE FROM foo.bar WHERE field1 <> 4;
    `
```

### MIGRATE (NOT YET IMPLEMENTED)
```
migrate top;
migrate bottom;
migreate up 2;
migrate down 2;
```

### GRANT and REVOKE
```
grant [privileges] on database;
grant [privileges] on schemas [except <schema_name>[,<schema_name> ...] ];
grant [privileges] on schema <schema_name>[,<schema_name> ...];
grant [privileges] on tables [except <schema_name.table_name>[,<schema_name.table_name> ...] ];
grant [privileges] on views [except <schema_name.view_name>[,<schema_name.view_name> ...] ];
grant [privileges] on functions [except <schema_name.function_name>[,<schema_name.function_name> ...] ];
grant [privileges] on procedures [except <schema_name.procedure_name>[,<schema_name.procedure_name> ...] ];
grant [privileges] on table <schema_name.table_name>[,<schema_name.table_name> ...];
grant [privileges] on view <schema_name.view_name>[,<schema_name.view_name> ...];
grant [privileges] on function <schema_name.function_name>[,<schema_name.function_name> ...];
grant [privileges] on procedure <schema_name.procedure_name>[,<schema_name.procedure_name> ...];
```

`revoke` syntax is the same as `grant`.

### SEED
```
seed cmd "SHELL COMMAND";
seed cmd -f /path/to/script.sh;
seed database [ (with | without) <seed_name>[,<seed_name> ...] ];
seed schema <schema_name> [ (with | without) <seed_name>[,<seed_name> ...] ];
seed tables [ (in | except in) <schema_name>[,<schema_name> ...] ];
seed table <schema_name.table_name>[,<schema_name.table_name> ...] ];
seed sql `
    UPDATE foo.bar SET field1 = 4 WHERE field2 = 0;
    DELETE FROM foo.bar WHERE field1 <> 4;
    `;
```

## Database Repo Structure

DDSL is opinionated about the structure of the database source repository.
The following structure is required.

```
📂 <any_parent_path>
  📂 <database_name>
    📄 database.create.sql
    📄 database.drop.sql
    📄 database.grank.sql 
    📄 database.revoke.sql
    📄 extensions.create.sql
    📄 extensions.drop.sql
    📄 privileges.grant.sql
    📄 privileges.revoke.sql
    📄 roles.create.sql
    📄 roles.drop.sql
    📂 schemas
      📂 <schema_name>
        📄 schema.create.sql
        📄 schema.drop.sql
        📄 privileges.grant.sql
        📄 privileges.revoke.sql
        📂 tables
          📂 <table_name>
            📄 table.create.sql
            📄 table.drop.sql
            📄 foreign_keys.create.sql
            📄 foreign_keys.drop.sql
            📄 indexes.create.sql
            📄 indexes.drop.sql
            📄 constraints.create.sql
            📄 constraints.drop.sql
            📄 privileges.grant.sql
            📄 privileges.revoke.sql
            📄 triggers.create.sql
            📄 triggers.drop.sql
            📂 seeds
              📄 table.csv
              📄 <seed_name>.sql
              📄 <seed_name>.csv
              📄 <seed_name>.sh
        📂 views
          📂 <view_name>
            📄 view.create.sql
            📄 view.drop.sql
            📄 indexes.create.sql
            📄 indexes.drop.sql
            📄 constraints.create.sql
            📄 constraints.drop.sql
            📄 privileges.grant.sql
            📄 privileges.revoke.sql
        📂 functions
          📂 <function_name>
            📄 function.create.sql
            📄 function.drop.sql
            📄 privileges.grant.sql
            📄 privileges.revoke.sql
        📂 procedures
          📂 <procedure_name>
            📄 procedure.create.sql
            📄 procedure.drop.sql
            📄 privileges.grant.sql
            📄 privileges.revoke.sql
        📂 types
          📄 <type_name>.create.sql
          📄 <type_name>.drop.sql
        📂 seeds
          📄 schema.ddsl
          📄 <seed_name>.ddsl
          📄 <seed_name>.sql
          📄 <seed_name>.sh
    📂 seeds
      📄 database.ddsl
      📄 <seed_name>.ddsl
      📄 <seed_name>.sql
      📄 <seed_name>.sh
    📂 migrations
      📄 <version>_<title>.up.ddsl
      📄 <version>_<title>.down.ddsl
```

## About Seeding

There are two types of seeds in Diesel. `Runtime Seeds` are meant to seed the database with required runtime data. These seeds are placed within the `database.ddsl`, `schema.ddsl`, and `table.ddsl` files. The file should contain the Diesel commands required to bring the database to a minimum usable state, such as seeding lookup tables, etc. Issuing a Diesel `seed` command without the `with` clause will execute the specified `Runtime Seed`. For example:

```
seed database;
seed schemas;
seed schema foo;
seed tables;
seed tables in foo;
seed table foo.bar;
```

The other type of seed in Diesel is the `Named Seed`. These are meant to seed the database for specific scenarios, for example for unit testing or integration testing. These named seeds may also be called from a `Runtime Seed` Diesel file. The `Named Seed` commands always have a `with` clause to specify the name of the desired seed. For example:

```
seed database with unit_test_1;
seed schemas with integration_test_1;
seed schema foo with performance_test_1;
seed tables with performance_test_1;
seed tables in foo with debug_data;
seed table foo.bar with debug_data;
```

