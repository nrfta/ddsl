package parser

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testSpec struct {
	command string
	expectedRoot string
	expectedPrimary string
	expectedClause string
	expectedArgs []string
	expectedExtendedArgs []string
}

var positiveTests = []testSpec{
	testSpec{"rollback","roolback","rollback", "", []string{},[]string{}},
	testSpec{"rollback transaction","roolback","rollback", "transaction", []string{},[]string{}},
	testSpec{"seed schema foo_schema", "seed","schema", "",[]string{"foo_schema"},[]string{}},
	testSpec{"seed schema foo_schema with foo_seed", "seed","schema", "with",[]string{"foo_schema"}, []string{"foo_seed"}},
	testSpec{"create database","create", "database", "", []string{}, []string{}},
	testSpec{"create roles","create", "roles", "",[]string{}, []string{}},
	testSpec{"create foreign-keys", "create","foreign-keys", "",[]string{}, []string{}},
	testSpec{"create schemas", "create","schemas", "",[]string{}, []string{}},
	testSpec{"create schemas except foo_schema", "create","schemas", "except",[]string{},[]string{"foo_schema"}},
	testSpec{"create schemas except foo_schema,bar_schema", "create","schemas", "except",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create schema foo_schema", "create","schema", "",[]string{},[]string{"foo_schema"}},
	testSpec{"create schema foo_schema,bar_schema", "create","schema", "",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create extensions", "create","extensions", "",[]string{}, []string{}},
	testSpec{"create extensions in foo_schema", "create","extensions", "in",[]string{},[]string{"foo_schema"}},
	testSpec{"create extensions in foo_schema,bar_schema", "create","extensions", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create extensions except in foo_schema", "create","extensions", "except in",[]string{},[]string{"foo_schema"}},
	testSpec{"create extensions except in foo_schema,bar_schema", "create","extensions", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create tables", "create","tables", "",[]string{}, []string{}},
	testSpec{"create tables in foo_schema", "create","tables", "in",[]string{},[]string{"foo_schema"}},
	testSpec{"create tables in foo_schema,bar_schema", "create","tables", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create tables except in foo_schema", "create","tables", "except in",[]string{},[]string{"foo_schema"}},
	testSpec{"create tables except in foo_schema,bar_schema", "create","tables", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create views", "create","views", "",[]string{}, []string{}},
	testSpec{"create views in foo_schema", "create","views", "in",[]string{},[]string{"foo_schema"}},
	testSpec{"create views in foo_schema,bar_schema", "create","views", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create views except in foo_schema", "create","views", "except in",[]string{},[]string{"foo_schema"}},
	testSpec{"create views except in foo_schema,bar_schema", "create","views", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create types", "create","types", "",[]string{}, []string{}},
	testSpec{"create types in foo_schema", "create","types", "in",[]string{},[]string{"foo_schema"}},
	testSpec{"create types in foo_schema,bar_schema", "create","types", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create types except in foo_schema", "create","types", "except in",[]string{},[]string{"foo_schema"}},
	testSpec{"create types except in foo_schema,bar_schema", "create","types", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"create table foo_schema.foo_table", "create","table", "", []string{},[]string{"foo_schema.foo_table"}},
	testSpec{"create table foo_schema.foo_table,bar_schema.bar_table", "create","table", "", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	testSpec{"create view foo_schema.foo_view", "create","view", "", []string{},[]string{"foo_schema.foo_view"}},
	testSpec{"create view foo_schema.foo_view,bar_schema.bar_view", "create","view", "", []string{},[]string{"foo_schema.foo_view","bar_schema.bar_view"}},
	testSpec{"create type foo_schema.foo_type", "create","type", "", []string{},[]string{"foo_schema.foo_type"}},
	testSpec{"create type foo_schema.foo_type,bar_schema.bar_type", "create","type", "", []string{},[]string{"foo_schema.foo_type","bar_schema.bar_type"}},
	testSpec{"create constraints on foo_schema.foo_table", "create","constraints", "on", []string{},[]string{"foo_schema.foo_table"}},
	testSpec{"create constraints on foo_schema.foo_table,bar_schema.bar_table", "create","constraints", "on", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	testSpec{"create indexes on foo_schema.foo_table", "create","indexes", "on", []string{},[]string{"foo_schema.foo_table"}},
	testSpec{"create indexes on foo_schema.foo_table,bar_schema.bar_table", "create","indexes", "on", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	testSpec{"migrate up 1", "migrate","up", "", []string{},[]string{"1"}},
	testSpec{"migrate down 1", "migrate","down", "", []string{},[]string{"1"}},
	testSpec{"migrate top", "migrate","top", "", []string{}, []string{}},
	testSpec{"migrate bottom", "migrate","bottom", "", []string{}, []string{}},
	testSpec{`seed cmd "python3 foo.py"`, "seed","cmd","", []string{},[]string{"python3 foo.py"}},
	testSpec{"seed database","create","database", "", []string{}, []string{}},
	testSpec{"seed database with foo_seed", "seed","database", "with", []string{},[]string{"foo_seed"}},
	testSpec{"seed database with foo_seed,bar_seed", "seed","database", "with", []string{},[]string{"foo_seed","bar_seed"}},
	testSpec{"seed database without foo_seed", "seed","database", "without", []string{},[]string{"foo_seed"}},
	testSpec{"seed database without foo_seed,bar_seed", "seed","database", "without", []string{},[]string{"foo_seed","bar_seed"}},
	testSpec{"seed schemas", "seed","schemas", "",[]string{}, []string{}},
	testSpec{"seed schemas except foo_schema", "seed","schemas", "except",[]string{},[]string{"foo_schema"}},
	testSpec{"seed schemas except foo_schema,bar_schema", "seed","schemas", "except",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"seed schema foo_schema", "seed","schema", "",[]string{"foo_schema"},[]string{}},
	testSpec{"seed schema foo_schema with foo_seed", "seed","schema", "with",[]string{"foo_schema"}, []string{"foo_seed"}},
	testSpec{"seed schema foo_schema with foo_seed,bar_seed", "seed","schema", "with",[]string{"foo_schema"}, []string{"foo_seed","bar_seed"}},
	testSpec{"seed schema foo_schema without foo_seed", "seed","schema", "without",[]string{"foo_schema"}, []string{"foo_seed"}},
	testSpec{"seed schema foo_schema without foo_seed,bar_seed", "seed","schema", "without",[]string{"foo_schema"}, []string{"foo_seed","bar_seed"}},
	testSpec{"seed table foo_schema.foo_table", "seed","table", "",[]string{},[]string{"foo_schema.foo_table"}},
	testSpec{"seed table foo_schema.foo_table,bar_schema.bar_table", "seed","table", "",[]string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	testSpec{"seed tables", "seed","tables", "",[]string{}, []string{}},
	testSpec{"seed tables in foo_schema", "seed","tables", "in",[]string{},[]string{"foo_schema"}},
	testSpec{"seed tables in foo_schema,bar_schema", "seed","tables", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"seed tables except in foo_schema", "seed","tables", "except in",[]string{},[]string{"foo_schema"}},
	testSpec{"seed tables except in foo_schema,bar_schema", "seed","tables", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on database", "grant","database", "", []string{}, []string{}},
	testSpec{"grant on database", "grant","database", "", []string{}, []string{}},
	testSpec{"grant on database", "grant","database", "", []string{}, []string{}},
	testSpec{"grant privileges on schemas", "grant","schemas", "", []string{}, []string{}},
	testSpec{"grant on schemas", "grant","schemas", "", []string{}, []string{}},
	testSpec{"grant privileges on schemas except foo_schema", "grant","schemas", "except", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on schemas except foo_schema,bar_schema", "grant","schemas", "except", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on schema foo_schema", "grant","schema", "", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on schema foo_schema,bar_schema", "grant","schema", "", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on tables", "grant","tables", "", []string{}, []string{}},
	testSpec{"grant on tables", "grant","tables", "", []string{}, []string{}},
	testSpec{"grant privileges on tables in foo_schema", "grant","tables", "in", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on tables in foo_schema,bar_schema", "grant","tables", "in", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on tables except in foo_schema", "grant","tables", "except in", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on tables except in foo_schema,bar_schema", "grant","tables", "except in", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on views", "grant","views", "", []string{}, []string{}},
	testSpec{"grant on views", "grant","views", "", []string{}, []string{}},
	testSpec{"grant privileges on views in foo_schema", "grant","views", "in", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on views in foo_schema,bar_schema", "grant","views", "in", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on views except in foo_schema", "grant","views", "except in", []string{},[]string{"foo_schema"}},
	testSpec{"grant privileges on views except in foo_schema,bar_schema", "grant","views", "except in", []string{},[]string{"foo_schema","bar_schema"}},
	testSpec{"grant privileges on table foo_schema.foo_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table"}},
	testSpec{"grant on table foo_schema.foo_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table"}},
	testSpec{"grant on table foo_schema.foo_table,bar_schema.bar_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	testSpec{"grant privileges on view foo_schema.foo_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view"}},
	testSpec{"grant on view foo_schema.foo_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view"}},
	testSpec{"grant on view foo_schema.foo_view,bar_schema.bar_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view","bar_schema.bar_view"}},
	testSpec{"begin","begin","begin", "", []string{},[]string{}},
	testSpec{"begin transaction","begin","begin", "transaction", []string{},[]string{}},
	testSpec{"commit","commit","commit", "", []string{},[]string{}},
	testSpec{"commit transaction","commit","commit", "transaction", []string{},[]string{}},
	testSpec{"rollback","roolback","rollback", "", []string{},[]string{}},
	testSpec{"rollback transaction","roolback","rollback", "transaction", []string{},[]string{}},
}

var _ = Describe("ParseTree", func() {

	Describe("Parse", func() {

		It("parses positive tests", func() {
			for _, s := range positiveTests {
				fmt.Fprintln(GinkgoWriter, s.command)
				cmd, err := Parse(s.command)
				Expect(err).To(BeNil())
				Expect(cmd.CommandDef.Name).To(Equal(s.expectedPrimary))
				Expect(cmd.Args).To(ConsistOf(s.expectedArgs))
				Expect(cmd.ExtArgs).To(ConsistOf(s.expectedExtendedArgs))
			}
		})
	})

})
