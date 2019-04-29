package parser_test

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/alecthomas/repr"
	"testing"
	"fmt"
)

func TestParser(t *testing.T) {
	command := `
		create DaTaBaSe
		create roles
		create extensions
		create foreign keys
		create schema foo
		create tables in foo
		create views in foo
		create table bar in foo
		create view cat in foo
		create indexes on foo.bar
		create constraints on foo.cat
		drop database
		drop roles
		drop extensions
		drop foreign keys
		drop schema foo
		drop tables in foo
		drop views in foo
		drop table bar in foo
		drop view cat in foo
		drop indexes on foo.bar
		drop constraints on foo.cat
		migrate top
		migrate bottom
		migrate up 2
		migrate down 2
	`
	command += "sql `UPDATE foo.bar SET field = 3;\n" +
		"" +
		"DELETE FROM foo.bar WHERE field = 5;\n" +
		"`"

	trees, err := parser.Parse(command)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for _, t := range trees {
		repr.Println(t, repr.Indent("  "), repr.OmitEmpty(true))
	}

}