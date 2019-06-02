# ddsl

Data Definition Support Language (DDSL, pronounced "diesel") provides a scripting language for DDL, migrations, and seeds. 

**Why a new language**

This project was born of the philosophy of "Database as Code". Database structures should be defined as code, versioned 
as code, released as code, and maintained as code. DDSL is DevOps for Database.

In addition, database code should be delivered as readily into unit tests, integration tests, and development environments
as it is into a production environment.

There are a few DevOps for Database tools in circulation already. What makes DDSL different is its full, unabashed
commitment to SQL. Database code should be stored and revisioned in its native language. Other DevOps for Database
tools require database structures to be defined in other languages such as YAML or XML.

DDSL helps you store your database DDL, migrations, and seeds in one revision control repository. It is opinionated
about the structure of the repository. That allows a set of simple commands to apply DDL structures, migrations, and seeds
agnostic of the database structure or RDS in use.

## Install

```$sh
go get github.com/neighborly/ddsl
```

## Usage

```$sh
# execute commands directly
ddsl -s <source_repo> -d <database_rds_url> COMMAND
# execute commands as a script
ddsl -s <source_repo> -d <database_rds_url> -f /pat/to/file.ddsl
# open a REPL shell
ddsl -s <source_repo> -d <database_rds_url>
```

The usage can be shortend by setting environment variables.

* `DDSL_SOURCE` - Source code repo URL for the database DDL and migrations
* `DDSL_DATABASE` - Database URL in format expected by RDS, properly URL encoded

The `--dry-run` switch will present what a command or script would do without making any changes.

## Command Syntax

Commands are not case sensitive. They may be separated by a semicolon and/or a newline.

### Databases

Databases cannot be created within a transaction on certain RDSs such as Postgres. When creating a database from scratch,
the recommended order of operations is:

1. `create roles`
2. `create database` 
3. `create schemas; create extensions; create tables; create views; etc..`    

### CREATE and DROP
```
create database
create roles
create foreign-keys
create schemas
create schema <schema_name>[,<schema_name> ...]
create extensions [ (in | except in) <schema_name>[,<schema_name> ...] ]
create tables [ (in | except in) <schema_name>[,<schema_name> ...] ]
create views [ (in | except in) <schema_name>[,<schema_name> ...] ]
create types [ (in | except in) <schema_name>[,<schema_name> ...] ]
create table <schema_name>.<table_name>[,<schema_name.table_name> ...]
create view <schema_name>.<view_name>[,<schema_name.view_name> ...]
create type <schema_name>.<type_name>[,<schema_name.type_name> ...]
create constraints on <schema_name>.<table_name>[,<schema_name.table_name> ...]
create indexes on <schema_name>.<table_or_view_name>[,<schema_name.table_or_view_name> ...]
```

`drop` syntax is the same as `create`.

### SQL
```
sql `
    UPDATE foo.bar SET field1 = 4 WHERE field2 = 0;
    DELETE FROM foo.bar WHERE field1 <> 4;
    `
```

### MIGRATE
```
migrate top
migrate bottom
migreate up 2
migrate down 2
```

### GRANT and REVOKE
```
grant [privileges] on database
grant [privileges] on schemas [except <schema_name>[,<schema_name> ...] ]
grant [privileges] on schema <schema_name>[,<schema_name> ...]
grant [privileges] on tables [except <schema_name.table_name>[,<schema_name.table_name> ...] ]
grant [privileges] on views [except <schema_name.view_name>[,<schema_name.view_name> ...] ]
grant [privileges] on table <schema_name.table_name>[,<schema_name.table_name> ...]
grant [privileges] on view <schema_name.view_name>[,<schema_name.view_name> ...]
```

`revoke` syntax is the same as `grant`.

### SEED
```
seed cmd "SHELL COMMAND"
seed cmd -f /path/to/script.sh
seed database [ (with | without) <seed_name>[,<seed_name> ...] ]
seed schema <schema_name> [ (with | without) <seed_name>[,<seed_name> ...] ]
seed tables [ (in | except in) <schema_name>[,<schema_name> ...] ]
seed table <schema_name.table_name>[,<schema_name.table_name> ...] ]
seed sql `
    UPDATE foo.bar SET field1 = 4 WHERE field2 = 0;
    DELETE FROM foo.bar WHERE field1 <> 4;
    `
```

## Database Repo Structure

DDSL is opinionated about the structure of the database source repository.
The following structure is required.

```
📂 <any_parent_path>
  📂 <database_name>
    📄 database.create.<ext> 
    📄 database.drop.<ext>
    📄 foreign_keys.create.<ext>  
    📄 foreign_keys.drop.<ext>
    📄 roles.create.<ext>
    📄 roles.drop.<ext>
    📂 schemas
      📂 <schema_name>
        📄 extensions.create.<ext>
        📄 extensions.drop.<ext>
        📄 schema.create.<ext>
        📄 schema.drop.<ext>
        📂 constraints
          📄 <table_or_view_name>.create.<ext>
          📄 <table_or_view_name>.drop.<ext>
        📂 indexes
          📄 <table_or_view_name>.create.<ext>
          📄 <table_or_view_name>.drop.<ext>
        📂 seeds
          📄 <seed_name>.ddsl
          📄 <seed_name>.sql
          📄 <seed_name>.sh
        📂 tables
          📄 <table_name>.create.<ext>
          📄 <table_name>.drop.<ext>
          📄 <table_name>.grant.<ext>
          📄 <table_name>.revoke.<ext>
          📄 <table_name>.seed.sql # or .csv or .sh
        📂 types
          📄 <type_name>.create.<ext>
          📄 <type_name>.drop.<ext>
        📂 views
          📄 <view_name>.create.<ext>
          📄 <view_name>.drop.<ext>
          📄 <table_name>.grant.<ext>
          📄 <table_name>.revoke.<ext>
    📂 seeds
      📄 <seed_name>.ddsl
      📄 <seed_name>.sql
      📄 <seed_name>.sh
    📂 migrations
      📄 <version>_<title>.up.ddsl
      📄 <version>_<title>.down.ddsl
```

