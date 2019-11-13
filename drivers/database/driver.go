// Package database provides the Database interface.
// All database drivers must implement this interface, register themselves,
// optionally provide a `WithInstance` function and pass the tests
// in package database/testing.
package database

import (
	"database/sql"
	"fmt"
	"io"
	nurl "net/url"
	"sync"
)

const (
	SchemaItemTypeTable     = "TABLE"
	SchemaItemTypeView      = "VIEW"
	SchemaItemTypeFunction  = "FUNCTION"
	SchemaItemTypeProcedure = "PROCEDURE"
	SchemaItemTypeType      = "TYPE"
	SQLQuerySchemas         = "SELECT schema_name FROM information_schema.schemata;"
	SQLQueryTables          = `
		SELECT 'TABLE' AS item_type, table_schema AS schema_name, table_name AS item_name FROM information_schema.tables 
		WHERE table_schema = '%s' AND table_type IN ('TABLE', 'BASE TABLE')
	`
	SQLQueryViews = `
		SELECT 'VIEW' AS item_type, table_schema AS schema_name, table_name AS item_name FROM information_schema.tables 
		WHERE table_schema = '%s' AND table_type = 'VIEW'
	`
	SQLQueryFunctions = `
		SELECT 'FUNCTION' AS item_type, specific_schema AS schema_name, routine_name AS item_name 
		FROM information_schema.routines
		WHERE specific_schema = '%s' AND routine_type = 'FUNCTION'
	`
	SQLQueryProcedures = `
		SELECT 'PROCEDURE' AS item_type, specific_schema AS schema_name, routine_name AS item_name 
		FROM information_schema.routines
		WHERE specific_schema = '%s' AND routine_type = 'PROCEDURE'
	`
	SQLQueryTypes = `
		SELECT 'TYPE' AS item_type, user_defined_type_schema AS schema_name, user_defined_type_name AS item_name 
		FROM information_schema.user_defined_types
		WHERE user_defined_type_schema = '%s'
	`
	SQLQueryForeignKeys = `
		SELECT
			tc.table_schema AS parent_schema_name, tc.table_name AS parent_item_name, 
			kcu.column_name AS parent_column_name,
			ccu.table_schema AS child_schema_name, ccu.table_name AS child_item_name,
			ccu.column_name AS child_column_name
		FROM
			information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage 
				AS kcu ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.constraint_column_usage 
				AS ccu ON ccu.constraint_name = tc.constraint_name
		WHERE constraint_type = 'FOREIGN KEY';
	`
)

var (
	// ErrLocked should be returned if a lock cannot be required on the database
	// when requested.
	ErrLocked = fmt.Errorf("can't acquire lock")

	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Driver is the interface every database driver must implement.
//
// How to implement a database driver?
//   1. Implement this interface.
//   2. Optionally, add a function named `WithInstance`.
//      This function should accept an existing DB instance and a Config{} struct
//      and return a driver instance.
//   3. Add a test that calls database/testing.go:Test()
//   4. Add own tests for Open(), WithInstance() (when provided) and Close().
//      All other functions are tested by tests in database/testing.
//      Saves you some time and makes sure all database drivers behave the same way.
//   5. Call Register in init().
//   6. Create a migrate/cli/build_<driver-name>.go file
//   7. Add driver name in 'DATABASE' variable in Makefile
//
// Guidelines:
//   * Don't try to correct user input. Don't assume things.
//     When in doubt, return an error and explain the situation to the user.
//   * All configuration input must come from the URL string in func Open()
//     or the Config{} struct in WithInstance. Don't os.Getenv().
type Driver interface {
	// Open returns a new driver instance configured with parameters
	// coming from the URL string. Migrate will call this function
	// only once per instance.
	Open(url string) (Driver, error)

	// Close closes the underlying database instance managed by the driver.
	// Migrate will call this function only once per instance.
	Close() error

	// Lock should acquire a database lock to control concurrency if required by
	// the application.
	Lock() error

	// Unlock should release the lock. Applications should call this when they
	// have completed interacting with the driver.
	Unlock() error

	// Begin should begin a transaction in the database. It should return an error if there is
	// already an active transaction in the database.
	Begin() error

	// Rollback should rollback a transaction in the database. It should return an error
	// if there is currently no active transaction.
	Rollback() error

	// Commit should commit the active transaction in the database. It should return an error
	// if there is currently no active transaction.
	Commit() error

	// Execute should execute the given command against the database.
	Exec(command io.Reader, params ...interface{}) error

	// Query should query the database and return results
	Query(command io.Reader, params ...interface{}) (*sql.Rows, error)

	// ImportCSV imports a csv file into the database.
	ImportCSV(filePath, schemaName, tableName, delimiter string, header bool) (output string, err error)

	// User returns the database user.
	User() string

	// DatabaseName returns the name of the database.
	DatabaseName() string

	// Schemas returns the names of the database schemas.
	Schemas() ([]string, error)

	// Tables returns the tables in the database schema.
	Tables(schema string) ([]*SchemaItemInfo, error)

	// Views returns the views in the database schema.
	Views(schema string) ([]*SchemaItemInfo, error)

	// Functions returns the functions in the database schema.
	Functions(schema string) ([]*SchemaItemInfo, error)

	// Procedures returns the procedures in the database schema.
	Procedures(schema string) ([]*SchemaItemInfo, error)

	// Types returns the types in the database schema.
	Types(schema string) ([]*SchemaItemInfo, error)

	// SchemaItems returns all items in the database schema
	SchemaItems(schema string) ([]*SchemaItemInfo, error)

	// Extensions returns the names of the database extensions
	Extensions() ([]string, error)

	// ForeignKeys returns the names of the foreign keys
	ForeignKeys() ([]*ForeignKeyInfo, error)

	// Roles returns the names of the database roles
	Roles() ([]string, error)
}

type SchemaItemInfo struct {
	// ItemType is one of TABLE, VIEW, FUNCTION, PROCEDURE, TYPE
	ItemType string

	// SchemaName is the name of the schema
	SchemaName string

	// ItemName is the name of the schema item
	ItemName string
}

type ForeignKeyInfo struct {
	ParentSchemaName string
	ParentTableName  string
	ParentColumnName string
	ChildSchemaName  string
	ChildTableName   string
	ChildColumnName  string
}

// Open returns a new driver instance.
func Open(url string) (Driver, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse URL. Did you escape all reserved URL characters? "+
			"See: https://github.com/golang-migrate/migrate#database-urls Error: %v", err)
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("database driver: invalid URL scheme")
	}

	driversMu.RLock()
	d, ok := drivers[u.Scheme]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("database driver: unknown driver %v (forgotten import?)", u.Scheme)
	}

	return d.Open(url)
}

// Register globally registers a driver.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// List lists the registered drivers
func List() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	names := make([]string, 0, len(drivers))
	for n := range drivers {
		names = append(names, n)
	}
	return names
}
