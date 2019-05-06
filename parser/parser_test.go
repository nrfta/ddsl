package parser_test

import (
	"fmt"
	"github.com/alecthomas/repr"
	"github.com/neighborly/ddsl/parser"
	"testing"
)

func TestParser(t *testing.T) {
	command := `
		create DaTaBaSe
		create roles @tags/v1.3
		create extensions
		create foreign keys @33fda2f3
		create schema foo
		create tables in foo @vtags/v1.3
		create views in foo
		create table bar in foo @vtags/1.3
		create view cat in foo
		create indexes on foo.bar @vtags/1.3
		create constraints on foo.cat

		drop database @tags/v1.3
		drop roles
		drop extensions @33fda2f3
		drop foreign keys
		drop schema foo @tags/v1.3
		drop tables in foo
		drop views in foo @tags/v1.3
		drop table bar in foo
		drop view cat in foo @tags/v1.3
		drop indexes on foo.bar
		drop constraints on foo.cat @tags/v1.3

		migrate top
		migrate bottom
		migrate up 2
		migrate down 2

	`
	command += "sql `UPDATE foo.bar SET field = 3;`\n" +
		"" +
		"sql `DELETE FROM foo.bar WHERE field = 5;\n" +
		"INSERT INTO foo.bar (field) VALUES (5);`" +
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
