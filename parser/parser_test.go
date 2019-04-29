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
		create roles @v1.3
		create extensions
		create foreign keys @v1.3
		create schema foo
		create tables in foo @v1.3
		create views in foo
		create table bar in foo @v1.3
		create view cat in foo
		create indexes on foo.bar @v1.3
		create constraints on foo.cat

		drop database @v1.3
		drop roles
		drop extensions @v1.3
		drop foreign keys
		drop schema foo @v1.3
		drop tables in foo
		drop views in foo @v1.3
		drop table bar in foo
		drop view cat in foo @v1.3
		drop indexes on foo.bar
		drop constraints on foo.cat @v1.3

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
