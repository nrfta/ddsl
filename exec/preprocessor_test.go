package exec

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ppTestSpec struct {
	command      string
	expectInstrs []*instruction
}

var sourceDir string

var posTests []ppTestSpec

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	src := os.Getenv("DDSL_SOURCE")
	sourceDir = filepath.Clean(filepath.Join(wd, strings.TrimLeft(src, "file://")))

	posTests = []ppTestSpec{
		{"create database", []*instruction{{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("database.create.sql")}}}},
		{"create roles", []*instruction{{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("roles.create.sql")}}}},
		{"create foreign-keys", []*instruction{{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("foreign_keys.create.sql")}}}},
		{"create extensions", []*instruction{{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("extensions.create.sql")}},}},
		{"create schemas", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/schema.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/schema.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/schema.create.sql")}},
		}},
		{"create schemas except baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/schema.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/schema.create.sql")}},
		}},
		{"create schemas except bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/schema.create.sql")}},
		}},
		{"create schema foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/schema.create.sql")}},
		}},
		{"create schema foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/schema.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/schema.create.sql")}},
		}},
		{"create tables", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/tables/foo_table/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create tables in foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create tables in foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create tables except in bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/baz_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create views", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/views/foo_view/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create views in foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create views in foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create views except in bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/baz_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create functions", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/foo_function/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/functions/foo_function/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/privileges.grant.sql")}},
		}},
		{"create functions in foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/privileges.grant.sql")}},
		}},
		{"create functions in foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/functions/foo_function/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/privileges.grant.sql")}},
		}},
		{"create functions except in bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/function.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/bar_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/baz_function/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/functions/foo_function/privileges.grant.sql")}},
		}},
		{"create procedures", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/foo_procedure/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/procedures/foo_procedure/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/privileges.grant.sql")}},
		}},
		{"create procedures in foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/privileges.grant.sql")}},
		}},
		{"create procedures in foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/procedures/foo_procedure/privileges.grant.sql")}},

			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/privileges.grant.sql")}},
		}},
		{"create procedures except in bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/procedure.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/bar_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/baz_procedure/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/procedures/foo_procedure/privileges.grant.sql")}},
		}},
		{"create types", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create types in foo_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create types in foo_schema,bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create types except in bar_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create types except in bar_schema,baz_schema", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create table foo_schema.foo_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create table foo_schema.foo_table,bar_schema.foo_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/table.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/tables/foo_table/privileges.grant.sql")}},
		}},
		{"create view foo_schema.foo_view", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create view foo_schema.foo_view,bar_schema.foo_view", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/privileges.grant.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/view.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/views/foo_view/privileges.grant.sql")}},
		}},
		{"create type foo_schema.foo_type", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
		}},
		{"create type foo_schema.foo_type,bar_schema.foo_type", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
		}},
		{"create constraints on foo_schema.foo_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
		}},
		{"create constraints on foo_schema.foo_table,foo_schema.bar_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/constraints.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/constraints.create.sql")}},
		}},
		{"create indexes on foo_schema.foo_view", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/foo_view/indexes.create.sql")}},
		}},
		{"create indexes on foo_schema.foo_table,foo_schema.bar_view", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/indexes.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/views/bar_view/indexes.create.sql")}},
		}},
		{"create triggers on foo_schema.foo_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/triggers.create.sql")}},
		}},
		{"create triggers on foo_schema.foo_table,foo_schema.bar_table", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/foo_table/triggers.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/tables/bar_table/triggers.create.sql")}},
		}},
		{`seed cmd "python3 script.py"`, []*instruction{
			{INSTR_SH_SCRIPT, map[string]interface{}{
				COMMAND: "python3",
				ARGS:    []string{"script.py"},
			}},
		}},
		{"seed database", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("seeds/database.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed schemas"}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in bar_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in baz_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in foo_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed database with foo_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/foo_seed.sql"),
				SEED_NAME: "foo_seed",
			}},
		}},
		{"seed database with foo_seed,bar_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/bar_seed.sql"),
				SEED_NAME: "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/foo_seed.sql"),
				SEED_NAME: "foo_seed",
			}},
		}},
		{"seed database without foo_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/bar_seed.sql"),
				SEED_NAME: "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/baz_seed.sql"),
				SEED_NAME: "baz_seed",
			}},
		}},
		{"seed database without foo_seed,bar_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH: filePath("seeds/baz_seed.sql"),
				SEED_NAME: "baz_seed",
			}},
		}},
		{"seed schemas", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in bar_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in baz_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in foo_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed schemas except foo_schema", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in bar_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in baz_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed schemas except foo_schema", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in bar_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in baz_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed schemas except foo_schema,bar_schema", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in baz_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed schema foo_schema", []*instruction{
			{INSTR_DDSL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/seeds/schema.ddsl")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "seed tables in foo_schema"}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_DDSL_FILE_END, map[string]interface{}{}},
		}},
		{"seed schema foo_schema with foo_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/foo_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "foo_seed",
			}},
		}},
		{"seed schema foo_schema with foo_seed,bar_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/bar_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/foo_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "foo_seed",
			}},
		}},
		{"seed schema foo_schema without foo_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/bar_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/baz_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "baz_seed",
			}},
		}},
		{"seed schema foo_schema,bar_schema without foo_seed,bar_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/seeds/baz_seed.sql"),
				SCHEMA_NAME: "bar_schema",
				SEED_NAME:   "baz_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/baz_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "baz_seed",
			}},
		}},
		{"seed schema foo_schema,bar_schema without foo_seed,bar_seed", []*instruction{
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/seeds/baz_seed.sql"),
				SCHEMA_NAME: "bar_schema",
				SEED_NAME:   "baz_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/seeds/baz_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				SEED_NAME:   "baz_seed",
			}},
		}},
		{"seed table foo_schema.foo_table", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed table foo_schema.foo_table,bar_schema.bar_table", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed table foo_schema.foo_table,bar_schema.bar_table with foo_seed,bar_seed", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/bar_seed.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
				SEED_NAME:   "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/foo_seed.sql"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
				SEED_NAME:   "foo_seed",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/bar_seed.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
				SEED_NAME:   "bar_seed",
			}},
			{INSTR_SQL_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/foo_seed.sql"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
				SEED_NAME:   "foo_seed",
			}},
		}},
		{"seed tables", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed tables in foo_schema", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed tables in foo_schema,bar_schema", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/foo_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "foo_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed tables except in foo_schema", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/bar_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "bar_schema",
				TABLE_NAME:  "foo_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"seed tables except in foo_schema,bar_schema", []*instruction{
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/bar_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "bar_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/baz_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "baz_table",
			}},
			{INSTR_CSV_FILE, map[string]interface{}{
				FILE_PATH:   filePath("schemas/baz_schema/tables/foo_table/seeds/table.csv"),
				SCHEMA_NAME: "baz_schema",
				TABLE_NAME:  "foo_table",
			}},
		}},
		{"begin transaction; create types; commit", []*instruction{
			{INSTR_BEGIN, map[string]interface{}{}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "create types"}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "commit"}},
			{INSTR_COMMIT, map[string]interface{}{}},
		}},
		{"begin; create types; commit transaction", []*instruction{
			{INSTR_BEGIN, map[string]interface{}{}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "create types"}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "commit transaction"}},
			{INSTR_COMMIT, map[string]interface{}{}},
		}},
		{"begin transaction; create types; rollback", []*instruction{
			{INSTR_BEGIN, map[string]interface{}{}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "create types"}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "rollback"}},
			{INSTR_ROLLBACK, map[string]interface{}{}},
		}},
		{"begin transaction; create types; rollback transaction", []*instruction{
			{INSTR_BEGIN, map[string]interface{}{}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "create types"}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/bar_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/baz_schema/types/foo_type.create.sql")}},
			{INSTR_SQL_FILE, map[string]interface{}{FILE_PATH: filePath("schemas/foo_schema/types/foo_type.create.sql")}},
			{INSTR_DDSL, map[string]interface{}{COMMAND: "rollback transaction"}},
			{INSTR_ROLLBACK, map[string]interface{}{}},
		}},
		/*
			ppTestSpec{"migrate up 1", "migrate", "up", "", []string{}, []string{"1"}},
			ppTestSpec{"migrate down 1", "migrate", "down", "", []string{}, []string{"1"}},
			ppTestSpec{"migrate top", "migrate", "top", "", []string{}, []string{}},
			ppTestSpec{"migrate bottom", "migrate", "bottom", "", []string{}, []string{}},
			ppTestSpec{"grant privileges on database", "grant", "database", "", []string{}, []string{}},
			ppTestSpec{"grant on database", "grant", "database", "", []string{}, []string{}},
			ppTestSpec{"grant on database", "grant", "database", "", []string{}, []string{}},
			ppTestSpec{"grant privileges on schemas", "grant", "schemas", "", []string{}, []string{}},
			ppTestSpec{"grant on schemas", "grant", "schemas", "", []string{}, []string{}},
			ppTestSpec{"grant privileges on schemas except foo_schema", "grant", "schemas", "except", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on schemas except foo_schema,bar_schema", "grant", "schemas", "except", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on schema foo_schema", "grant", "schema", "", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on schema foo_schema,bar_schema", "grant", "schema", "", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on tables", "grant", "tables", "", []string{}, []string{}},
			ppTestSpec{"grant on tables", "grant", "tables", "", []string{}, []string{}},
			ppTestSpec{"grant privileges on tables in foo_schema", "grant", "tables", "in", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on tables in foo_schema,bar_schema", "grant", "tables", "in", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on tables except in foo_schema", "grant", "tables", "except in", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on tables except in foo_schema,bar_schema", "grant", "tables", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on views", "grant", "views", "", []string{}, []string{}},
			ppTestSpec{"grant on views", "grant", "views", "", []string{}, []string{}},
			ppTestSpec{"grant privileges on views in foo_schema", "grant", "views", "in", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on views in foo_schema,bar_schema", "grant", "views", "in", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on views except in foo_schema", "grant", "views", "except in", []string{}, []string{"foo_schema"}},
			ppTestSpec{"grant privileges on views except in foo_schema,bar_schema", "grant", "views", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
			ppTestSpec{"grant privileges on table foo_schema.foo_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table"}},
			ppTestSpec{"grant on table foo_schema.foo_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table"}},
			ppTestSpec{"grant on table foo_schema.foo_table,bar_schema.bar_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
			ppTestSpec{"grant privileges on view foo_schema.foo_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view"}},
			ppTestSpec{"grant on view foo_schema.foo_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view"}},
			ppTestSpec{"grant on view foo_schema.foo_view,bar_schema.bar_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view", "bar_schema.bar_view"}},
			ppTestSpec{"rollback", "roolback", "rollback", "", []string{}, []string{}},
			ppTestSpec{"rollback transaction", "roolback", "rollback", "transaction", []string{}, []string{}},
		*/
	}
}

var _ = ginkgo.Describe("preprocessor.go", func() {

	ginkgo.Describe("positive tests", func() {

		ginkgo.It("positive tests succeed", func() {

			for _, pt := range posTests {
				ctx := &Context{SourceRepo: "file://" + sourceDir}
				cmds, _, _, err := parser.Parse(pt.command)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				c, err := preprocessBatch(ctx, cmds)
				if err != nil {
					ginkgo.Fail(fmt.Sprintf("preprocess '%s'; %s", pt.command, err.Error()))
				}
				if c == 0 || len(ctx.instructions) == 0 {
					ginkgo.Fail("no instructions")
				}
				Expect(strings.HasPrefix(pt.command, ctx.instructions[0].params[COMMAND].(string))).To(BeTrue())
				instrs := ctx.instructions[1:]
				Expect(len(instrs)).To(Equal(len(pt.expectInstrs)))
				for i, instr := range instrs {
					Expect(instr).To(Equal(pt.expectInstrs[i]), fmt.Sprintf("%d: %s", i, pt.command))
				}
			}

		})

	})
})

func filePath(relativePath string) string {
	return path.Join(sourceDir, relativePath)
}
