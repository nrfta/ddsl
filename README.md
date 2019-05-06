# ddsl

Data-Definition-Specific Language (DDSL, pronounced "diesel") provides a scripting language for DDL and migrations. 

**Why a new language**

A relational database system (RDS) is not a source code repository. DDL needs to be stored and versioned separately
in order to manage it like any other code artifact. There are few tools to manage and apply DDL thus stored 
during release and upgrade activities.

DDSL helps you store your database DDL and migrations in one revision control repository. It is opinionated
about the structure of the repository. That allows a set of simple commands to apply DDL objects and migrations
agnostic of the database structure or RDS in use.

## Install

```$sh
go get github.com/neighborly/ddsl
```

## Usage

```$sh
ddsl -s <source_repo> -d <database_rds_url> -c COMMAND
ddsl -s <source_repo> -d <database_rds_url> -f /pat/to/file.ddsl
ddsl --version
```

The usage can be shortend by setting environment variables.

* `DDSL_SOURCE` - Source code repo URL for the database DDL and migrations
* `DDSL_DATABASE` - Database URL in format expected by RDS, properly URL encoded

## Command Syntax

All commands accept a final token of `@git_tag` which will run the command against that version of the DDL reposititory.

### CREATE
```
CREATE DATABASE
CREATE ROLES
CREATE EXTENSIONS
CREATE FOREIGN KEYS
CREATE SCHEMA foo 
CREATE TABLES IN foo 
CREATE VIEWS IN foo
CREATE TABLE bar IN foo
CREATE VIEW cat IN foo
CREATE INDEXES ON foo.bar @v1.1
CREATE CONSTRAINTS ON foo.cat @v1.2
```

### DROP
```
DROP DATABASE
DROP ROLES
DROP EXTENSIONS
DROP FOREIGN KEYS
DROP SCHEMA foo
DROP TABLES IN foo
DROP VIEWS IN foo
DROP TABLE bar IN foo
DROP VIEW cat IN foo
DROP INDEXES ON foo.bar @v1.1
DROP CONSTRAINTS ON foo.cat @v1.2
```

### SQL
```
SQL `
    UPDATE foo.bar SET field1 = 4 WHERE field2 = 0;
    DELETE FROM foo.bar WHERE field1 <> 4;
    `
```

### MIGRATE
```
MIGRATE TOP
MIGRATE BOTTOM
MIGRATE UP 2
MIGRATE DOWN 2
```

## Database Repo Structure

DDSL is opinionated about the structure of the database source repository.
The following structure is required.

```
📂 <any_parent_path>
  📂 <database_name>
    📄 database.create.<ext> 
    📄 database.drop.<ext>
    📄 extensions.create.<ext>
    📄 extensions.drop.<ext>
    📄 foreign_keys.create.<ext>  
    📄 foreign_keys.drop.<ext>
    📄 roles.create.<ext>
    📄 roles.drop.<ext>
    📂 schemas
      📂 <schema_name>
        📄 schema.create.<ext>
        📄 schema.drop.<ext>
        📂 constraints
          📄 <table_or_view_name>.create.<ext>
          📄 <table_or_view_name>.drop.<ext>
        📂 indexes
          📄 <table_or_view_name>.create.<ext>
          📄 <table_or_view_name>.drop.<ext>
        📂 tables
          📄 <table_name>.create.<ext>
          📄 <table_name>.drop.<ext>
        📂 views
          📄 <view_name>.create.<ext>
          📄 <view_name>.drop.<ext>
      📂 migrations
        📄 <version>_<title>.up.ddsl
        📄 <version>_<title>.down.ddsl
```

Migrations are written in DDSL because often migrations simply need to create a specific table
or index, or drop something. The DDL for that already exists in some version of the database 
code repository, so it is DRY to be able to access that code from within the migrations. You 
can also run SQL commands in the migrations using the DDSL `SQL` command.
