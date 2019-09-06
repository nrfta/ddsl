package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/neighborly/ddsl/drivers/database"
	"github.com/neighborly/ddsl/util"
	"io"
	"io/ioutil"
	nurl "net/url"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func init() {
	db := Postgres{}
	database.Register("postgres", &db)
	database.Register("postgresql", &db)
}

var (
	ErrNilConfig      = fmt.Errorf("no config")
	ErrNoDatabaseName = fmt.Errorf("no database name")
	ErrNoSchema       = fmt.Errorf("no schema")
	ErrNoUser         = fmt.Errorf("no user")
	ErrDatabaseDirty  = fmt.Errorf("database is dirty")
)

type Config struct {
	DatabaseName string
	SchemaName   string
	URL          string
	User         string
}

type Postgres struct {
	// Locking and unlocking need to use the same connection
	conn     *sql.Conn
	db       *sql.DB
	isLocked bool
	tx       *sql.Tx

	// Open and WithInstance need to guarantee that config is never nil
	config *Config
}
// compile-time interface compliance
var _ database.Driver = (*Postgres)(nil)

func Register() {
	// do nothing, but call to force compiler to accept import without use
}

func WithInstance(instance *sql.DB, config *Config) (database.Driver, error) {
	if config == nil {
		return nil, ErrNilConfig
	}

	if err := instance.Ping(); err != nil {
		return nil, err
	}

	query := `SELECT CURRENT_DATABASE()`
	var databaseName string
	if err := instance.QueryRow(query).Scan(&databaseName); err != nil {
		return nil, &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if len(databaseName) == 0 {
		return nil, ErrNoDatabaseName
	}

	config.DatabaseName = databaseName

	query = `SELECT CURRENT_SCHEMA()`
	var schemaName string
	if err := instance.QueryRow(query).Scan(&schemaName); err != nil {
		return nil, &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if len(schemaName) == 0 {
		return nil, ErrNoSchema
	}

	config.SchemaName = schemaName

	query = `SELECT CURRENT_USER`
	var user string
	if err := instance.QueryRow(query).Scan(&user); err != nil {
		return nil, &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if len(user) == 0 {
		return nil, ErrNoUser
	}

	config.User = user

	conn, err := instance.Conn(context.Background())

	if err != nil {
		return nil, err
	}

	px := &Postgres{
		conn:   conn,
		db:     instance,
		config: config,
	}

	return px, nil
}

func (p *Postgres) Open(url string) (database.Driver, error) {
	purl, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", migrate.FilterCustomQuery(purl).String())
	if err != nil {
		return nil, err
	}

	px, err := WithInstance(db, &Config{
		DatabaseName: purl.Path,
		URL:          url,
	})

	if err != nil {
		return nil, err
	}

	return px, nil
}

func (p *Postgres) Close() error {
	connErr := p.conn.Close()
	dbErr := p.db.Close()
	if connErr != nil || dbErr != nil {
		return fmt.Errorf("conn: %v, db: %v", connErr, dbErr)
	}
	return nil
}

// https://www.postgresql.org/docs/9.6/static/explicit-locking.html#ADVISORY-LOCKS
func (p *Postgres) Lock() error {
	if p.isLocked {
		return database.ErrLocked
	}

	aid, err := database.GenerateAdvisoryLockId(p.config.DatabaseName, p.config.SchemaName)
	if err != nil {
		return err
	}

	// This will either obtain the lock immediately and return true,
	// or return false if the lock cannot be acquired immediately.
	query := `SELECT pg_advisory_lock($1)`
	if _, err := p.conn.ExecContext(context.Background(), query, aid); err != nil {
		return &database.Error{OrigErr: err, Err: "try lock failed", Query: []byte(query)}
	}

	p.isLocked = true
	return nil
}

func (p *Postgres) Unlock() error {
	if !p.isLocked {
		return nil
	}

	aid, err := database.GenerateAdvisoryLockId(p.config.DatabaseName, p.config.SchemaName)
	if err != nil {
		return err
	}

	query := `SELECT pg_advisory_unlock($1)`
	if _, err := p.conn.ExecContext(context.Background(), query, aid); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	p.isLocked = false
	return nil
}

func (p *Postgres) User() string {
	return p.config.User
}

func (p *Postgres) DatabaseName() string {
	return p.config.DatabaseName
}

func (p *Postgres) Begin() error {
	if p.tx != nil {
		return &database.Error{Err: "connection is already in transaction"}
	}

	opts := &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	}
	tx, err := p.conn.BeginTx(context.Background(), opts)
	if err != nil {
		return &database.Error{OrigErr: err, Err: "error beginning transaction"}
	}

	_, err = p.conn.ExecContext(context.Background(), "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;")
	if err != nil {
		tx.Rollback()
		return &database.Error{OrigErr: err, Err: "error setting transaction isolation level"}
	}

	p.tx = tx

	return nil
}

func (p *Postgres) Rollback() error {
	if p.tx == nil {
		return &database.Error{Err: "connection is not in transaction"}
	}

	err := p.tx.Rollback()
	if err != nil {
		return &database.Error{OrigErr: err, Err: "error rolling back transaction"}
	}

	p.tx = nil

	return nil
}

func (p *Postgres) Commit() error {
	if p.tx == nil {
		return &database.Error{Err: "connection is not in transaction"}
	}

	err := p.tx.Commit()
	if err != nil {
		return &database.Error{OrigErr: err, Err: "error committing transaction"}
	}

	p.tx = nil

	return nil
}

func (p *Postgres) Exec(command io.Reader, params ...interface{}) error {
	cmdBytes, err := ioutil.ReadAll(command)
	if err != nil {
		return err
	}

	if params == nil {
		params = []interface{}{}
	}

	cmd := string(cmdBytes[:])
	if _, err = p.conn.ExecContext(context.Background(), cmd, params...); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			var line uint
			var col uint
			var lineColOK bool
			if pgErr.Position != "" {
				if pos, err := strconv.ParseUint(pgErr.Position, 10, 64); err == nil {
					line, col, lineColOK = computeLineFromPos(cmd, int(pos))
				}
			}
			message := pgErr.Message
			if lineColOK {
				message = fmt.Sprintf("%s (column %d)", message, col)
			}
			if pgErr.Detail != "" {
				message = fmt.Sprintf("%s, %s", message, pgErr.Detail)
			}
			return database.Error{OrigErr: err, Err: message, Query: cmdBytes, Line: line}
		}
		return database.Error{OrigErr: err, Err: "command failed", Query: cmdBytes}
	}

	return nil
}

func (p *Postgres) Query(command io.Reader, params ...interface{}) (*sql.Rows, error) {
	cmdBytes, err := ioutil.ReadAll(command)
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = []interface{}{}
	}

	cmd := string(cmdBytes[:])
	var rows *sql.Rows
	if rows, err = p.conn.QueryContext(context.Background(), cmd, params...); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			var line uint
			var col uint
			var lineColOK bool
			if pgErr.Position != "" {
				if pos, err := strconv.ParseUint(pgErr.Position, 10, 64); err == nil {
					line, col, lineColOK = computeLineFromPos(cmd, int(pos))
				}
			}
			message := pgErr.Message
			if lineColOK {
				message = fmt.Sprintf("%s (column %d)", message, col)
			}
			if pgErr.Detail != "" {
				message = fmt.Sprintf("%s, %s", message, pgErr.Detail)
			}
			return nil, database.Error{OrigErr: err, Err: message, Query: cmdBytes, Line: line}
		}
		return nil, database.Error{OrigErr: err, Err: "command failed", Query: cmdBytes}
	}

	return rows, nil
}

func (p *Postgres) ImportCSV(filePath, schemaName, tableName, delimiter string, header bool) (output string, err error) {
	sql := fmt.Sprintf("\\COPY %s.%s FROM '%s' WITH DELIMITER '%s' CSV", schemaName, tableName, filePath, delimiter)
	if header {
		sql += " HEADER;"
	} else {
		sql += ";"
	}
	out, err := util.OSExec("psql", p.config.URL, "-q", "-c", sql)
	if err != nil {
		return out, err
	}

	return out, nil
}

func computeLineFromPos(s string, pos int) (line uint, col uint, ok bool) {
	// replace crlf with lf
	s = strings.Replace(s, "\r\n", "\n", -1)
	// pg docs: pos uses index 1 for the first character, and positions are measured in characters not bytes
	runes := []rune(s)
	if pos > len(runes) {
		return 0, 0, false
	}
	sel := runes[:pos]
	line = uint(runesCount(sel, newLine) + 1)
	col = uint(pos - 1 - runesLastIndex(sel, newLine))
	return line, col, true
}

const newLine = '\n'

func runesCount(input []rune, target rune) int {
	var count int
	for _, r := range input {
		if r == target {
			count++
		}
	}
	return count
}

func runesLastIndex(input []rune, target rune) int {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == target {
			return i
		}
	}
	return -1
}
