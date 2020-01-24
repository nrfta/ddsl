// Package testing has the database tests.
// All database drivers must pass the Test function.
// This lives in it's own package so it stays a test dependency.
package testing

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/nrfta/ddsl/drivers/database"
)

// Test runs tests against database implementations.
func Test(t *testing.T, d database.Driver, command []byte, params ...interface{}) {
	if command == nil {
		t.Fatal("test must provide command")
	}

	TestLockAndUnlock(t, d)
	TestExec(t, d, bytes.NewReader(command), params)
	TestTransaction(t, d)
}

func TestLockAndUnlock(t *testing.T, d database.Driver) {
	// add a timeout, in case there is a deadlock
	done := make(chan struct{})
	errs := make(chan error)

	go func() {
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-done:
				return
			case <-timeout:
				errs <- fmt.Errorf("Timeout after 15 seconds. Looks like a deadlock in Lock/UnLock.\n%#v", d)
				return
			}
		}
	}()

	// run the locking test ...
	go func() {
		if err := d.Lock(); err != nil {
			errs <- err
			return
		}

		// try to acquire lock again
		if err := d.Lock(); err == nil {
			errs <- errors.New("lock: expected err not to be nil")
			return
		}

		// unlock
		if err := d.Unlock(); err != nil {
			errs <- err
			return
		}

		// try to lock
		if err := d.Lock(); err != nil {
			errs <- err
			return
		}
		if err := d.Unlock(); err != nil {
			errs <- err
			return
		}
		// notify everyone
		close(done)
	}()

	// wait for done or any error
	for {
		select {
		case <-done:
			return
		case err := <-errs:
			t.Fatal(err)
		}
	}
}

func TestExec(t *testing.T, d database.Driver, command io.Reader, params []interface{}) {
	if command == nil {
		t.Fatal("command can't be nil")
	}

	if err := d.Exec(command, params); err != nil {
		t.Fatal(err)
	}
}

func TestTransaction(t *testing.T, d database.Driver) {
	if err := d.Rollback(); err == nil {
		t.Fatal("rollback without a transaction")
	}

	if err := d.Commit(); err == nil {
		t.Fatal("commit without a transaction")
	}

	if err := d.Begin(); err != nil {
		t.Fatal(err)
	}

	if err := d.Begin(); err == nil {
		t.Fatal("already in transaction")
	}

	if err := d.Rollback(); err != nil {
		t.Fatal(err)
	}

	if err := d.Begin(); err != nil {
		t.Fatal(err)
	}

	if err := d.Commit(); err != nil {
		t.Fatal(err)
	}
}
