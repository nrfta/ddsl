package parser

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type posTestSpec struct {
	command              string
	expectedRoot         string
	expectedPrimary      string
	expectedClause       string
	expectedArgs         []string
	expectedExtendedArgs []string
}

type negTestSpec struct {
	command string
	message string
}

var positiveTests = []posTestSpec{
	{"create database", "create", "database", "", []string{}, []string{}},
	{"create DATABASE", "create", "database", "", []string{}, []string{}},
	{"CREATE DATABASE", "create", "database", "", []string{}, []string{}},
	{"CREATE database", "create", "database", "", []string{}, []string{}},
	{"cReAtE dAtAbAsE", "create", "database", "", []string{}, []string{}},
	{"create roles", "create", "roles", "", []string{}, []string{}},
	{"create schemas", "create", "schemas", "", []string{}, []string{}},
	{"create schemas except foo_schema", "create", "schemas", "except", []string{}, []string{"foo_schema"}},
	{"create schemas except foo_schema,bar_schema", "create", "schemas", "except", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create schema foo_schema", "create", "schema", "", []string{}, []string{"foo_schema"}},
	{"CREATE SCHEMA FOO_SCHEMA", "create", "schema", "", []string{}, []string{"FOO_SCHEMA"}},
	{"create schema foo_schema,bar_schema", "create", "schema", "", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create extensions", "create", "extensions", "", []string{}, []string{}},
	{"create extensions in foo_schema", "create", "extensions", "in", []string{}, []string{"foo_schema"}},
	{"create extensions in foo_schema,bar_schema", "create", "extensions", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create extensions except in foo_schema", "create", "extensions", "except in", []string{}, []string{"foo_schema"}},
	{"create extensions except in foo_schema,bar_schema", "create", "extensions", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create tables", "create", "tables", "", []string{}, []string{}},
	{"create tables in foo_schema", "create", "tables", "in", []string{}, []string{"foo_schema"}},
	{"create tables in foo_schema,bar_schema", "create", "tables", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create tables except in foo_schema", "create", "tables", "except in", []string{}, []string{"foo_schema"}},
	{"create tables except in foo_schema,bar_schema", "create", "tables", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create views", "create", "views", "", []string{}, []string{}},
	{"create views in foo_schema", "create", "views", "in", []string{}, []string{"foo_schema"}},
	{"create views in foo_schema,bar_schema", "create", "views", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create views except in foo_schema", "create", "views", "except in", []string{}, []string{"foo_schema"}},
	{"create views except in foo_schema,bar_schema", "create", "views", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create types", "create", "types", "", []string{}, []string{}},
	{"create types in foo_schema", "create", "types", "in", []string{}, []string{"foo_schema"}},
	{"create types in foo_schema,bar_schema", "create", "types", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create types except in foo_schema", "create", "types", "except in", []string{}, []string{"foo_schema"}},
	{"create types except in foo_schema,bar_schema", "create", "types", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create foreign-keys", "create", "foreign-keys", "", []string{}, []string{}},
	{"create foreign-keys in foo_schema", "create", "foreign-keys", "in", []string{}, []string{"foo_schema"}},
	{"create foreign-keys in foo_schema,bar_schema", "create", "foreign-keys", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create foreign-keys except in foo_schema", "create", "foreign-keys", "except in", []string{}, []string{"foo_schema"}},
	{"create foreign-keys except in foo_schema,bar_schema", "create", "foreign-keys", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"create foreign-keys on foo_schema.foo_table", "create", "foreign-keys", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"create foreign-keys on foo_schema.foo_table,bar_schema.bar_table", "create", "foreign-keys", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"create table foo_schema.foo_table", "create", "table", "", []string{}, []string{"foo_schema.foo_table"}},
	{"create table foo_schema.foo_table,bar_schema.bar_table", "create", "table", "", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"create view foo_schema.foo_view", "create", "view", "", []string{}, []string{"foo_schema.foo_view"}},
	{"create view foo_schema.foo_view,bar_schema.bar_view", "create", "view", "", []string{}, []string{"foo_schema.foo_view", "bar_schema.bar_view"}},
	{"create type foo_schema.foo_type", "create", "type", "", []string{}, []string{"foo_schema.foo_type"}},
	{"create type foo_schema.foo_type,bar_schema.bar_type", "create", "type", "", []string{}, []string{"foo_schema.foo_type", "bar_schema.bar_type"}},
	{"create constraints on foo_schema.foo_table", "create", "constraints", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"create constraints on foo_schema.foo_table,bar_schema.bar_table", "create", "constraints", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"create indexes on foo_schema.foo_table", "create", "indexes", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"create indexes on foo_schema.foo_table,bar_schema.bar_table", "create", "indexes", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"list foreign-keys", "list", "foreign-keys", "", []string{}, []string{}},
	{"list foreign-keys in foo_schema", "list", "foreign-keys", "in", []string{}, []string{"foo_schema"}},
	{"list foreign-keys in foo_schema,bar_schema", "list", "foreign-keys", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list foreign-keys except in foo_schema", "list", "foreign-keys", "except in", []string{}, []string{"foo_schema"}},
	{"list foreign-keys except in foo_schema,bar_schema", "list", "foreign-keys", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list schemas", "list", "schemas", "", []string{}, []string{}},
	{"list tables", "list", "tables", "", []string{}, []string{}},
	{"list tables in foo_schema", "list", "tables", "in", []string{}, []string{"foo_schema"}},
	{"list tables in foo_schema,bar_schema", "list", "tables", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list tables except in foo_schema", "list", "tables", "except in", []string{}, []string{"foo_schema"}},
	{"list tables except in foo_schema,bar_schema", "list", "tables", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list views", "list", "views", "", []string{}, []string{}},
	{"list views in foo_schema", "list", "views", "in", []string{}, []string{"foo_schema"}},
	{"list views in foo_schema,bar_schema", "list", "views", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list views except in foo_schema", "list", "views", "except in", []string{}, []string{"foo_schema"}},
	{"list views except in foo_schema,bar_schema", "list", "views", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list types", "list", "types", "", []string{}, []string{}},
	{"list types in foo_schema", "list", "types", "in", []string{}, []string{"foo_schema"}},
	{"list types in foo_schema,bar_schema", "list", "types", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list types except in foo_schema", "list", "types", "except in", []string{}, []string{"foo_schema"}},
	{"list types except in foo_schema,bar_schema", "list", "types", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list functions", "list", "functions", "", []string{}, []string{}},
	{"list functions in foo_schema", "list", "functions", "in", []string{}, []string{"foo_schema"}},
	{"list functions in foo_schema,bar_schema", "list", "functions", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list functions except in foo_schema", "list", "functions", "except in", []string{}, []string{"foo_schema"}},
	{"list functions except in foo_schema,bar_schema", "list", "functions", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list procedures", "list", "procedures", "", []string{}, []string{}},
	{"list procedures in foo_schema", "list", "procedures", "in", []string{}, []string{"foo_schema"}},
	{"list procedures in foo_schema,bar_schema", "list", "procedures", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list procedures except in foo_schema", "list", "procedures", "except in", []string{}, []string{"foo_schema"}},
	{"list procedures except in foo_schema,bar_schema", "list", "procedures", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"list constraints on foo_schema.foo_table", "list", "constraints", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"list constraints on foo_schema.foo_table,bar_schema.bar_table", "list", "constraints", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"list indexes on foo_schema.foo_table", "list", "indexes", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"list indexes on foo_schema.foo_table,bar_schema.bar_table", "list", "indexes", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"list triggers on foo_schema.foo_table", "list", "triggers", "on", []string{}, []string{"foo_schema.foo_table"}},
	{"list triggers on foo_schema.foo_table,bar_schema.bar_table", "list", "triggers", "on", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"migrate up 1", "migrate", "up", "", []string{}, []string{"1"}},
	{"migrate down 1", "migrate", "down", "", []string{}, []string{"1"}},
	{"migrate top", "migrate", "top", "", []string{}, []string{}},
	{"migrate bottom", "migrate", "bottom", "", []string{}, []string{}},
	{`seed cmd "python3 foo.py"`, "seed", "cmd", "", []string{}, []string{"python3 foo.py"}},
	{"seed database", "create", "database", "", []string{}, []string{}},
	{"seed database with foo_seed", "seed", "database", "with", []string{}, []string{"foo_seed"}},
	{"seed database with foo_seed,bar_seed", "seed", "database", "with", []string{}, []string{"foo_seed", "bar_seed"}},
	{"seed database without foo_seed", "seed", "database", "without", []string{}, []string{"foo_seed"}},
	{"seed database without foo_seed,bar_seed", "seed", "database", "without", []string{}, []string{"foo_seed", "bar_seed"}},
	{"seed schemas", "seed", "schemas", "", []string{}, []string{}},
	{"seed schemas except foo_schema", "seed", "schemas", "except", []string{}, []string{"foo_schema"}},
	{"seed schemas except foo_schema,bar_schema", "seed", "schemas", "except", []string{}, []string{"foo_schema", "bar_schema"}},
	{"seed schema foo_schema", "seed", "schema", "", []string{"foo_schema"}, []string{}},
	{"seed schema foo_schema with foo_seed", "seed", "schema", "with", []string{"foo_schema"}, []string{"foo_seed"}},
	{"SEED SCHEMA FOO_SCHEMA WITH FOO_SEED", "seed", "schema", "with", []string{"FOO_SCHEMA"}, []string{"FOO_SEED"}},
	{"seed schema foo_schema with foo_seed,bar_seed", "seed", "schema", "with", []string{"foo_schema"}, []string{"foo_seed", "bar_seed"}},
	{"seed schema foo_schema without foo_seed", "seed", "schema", "without", []string{"foo_schema"}, []string{"foo_seed"}},
	{"seed schema foo_schema,bar_schema without foo_seed,bar_seed", "seed", "schema", "without", []string{"foo_schema", "bar_schema"}, []string{"foo_seed", "bar_seed"}},
	{"seed table foo_schema.foo_table", "seed", "table", "", []string{"foo_schema.foo_table"}, []string{}},
	{"seed table foo_schema.foo_table,bar_schema.bar_table", "seed", "table", "", []string{"foo_schema.foo_table", "bar_schema.bar_table"}, []string{}},
	{"seed table foo_schema.foo_table,bar_schema.bar_table with foo_seed,bar_seed", "seed", "table", "with", []string{"foo_schema.foo_table", "bar_schema.bar_table"}, []string{"foo_seed", "bar_seed"}},
	{"seed tables", "seed", "tables", "", []string{}, []string{}},
	{"seed tables in foo_schema", "seed", "tables", "in", []string{}, []string{"foo_schema"}},
	{"seed tables in foo_schema,bar_schema", "seed", "tables", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"seed tables except in foo_schema", "seed", "tables", "except in", []string{}, []string{"foo_schema"}},
	{"seed tables except in foo_schema,bar_schema", "seed", "tables", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on database", "grant", "database", "", []string{}, []string{}},
	{"grant on database", "grant", "database", "", []string{}, []string{}},
	{"grant on database", "grant", "database", "", []string{}, []string{}},
	{"grant privileges on schemas", "grant", "schemas", "", []string{}, []string{}},
	{"grant on schemas", "grant", "schemas", "", []string{}, []string{}},
	{"grant privileges on schemas except foo_schema", "grant", "schemas", "except", []string{}, []string{"foo_schema"}},
	{"grant privileges on schemas except foo_schema,bar_schema", "grant", "schemas", "except", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on schema foo_schema", "grant", "schema", "", []string{}, []string{"foo_schema"}},
	{"grant privileges on schema foo_schema,bar_schema", "grant", "schema", "", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on tables", "grant", "tables", "", []string{}, []string{}},
	{"grant on tables", "grant", "tables", "", []string{}, []string{}},
	{"grant privileges on tables in foo_schema", "grant", "tables", "in", []string{}, []string{"foo_schema"}},
	{"grant privileges on tables in foo_schema,bar_schema", "grant", "tables", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on tables except in foo_schema", "grant", "tables", "except in", []string{}, []string{"foo_schema"}},
	{"grant privileges on tables except in foo_schema,bar_schema", "grant", "tables", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on views", "grant", "views", "", []string{}, []string{}},
	{"grant on views", "grant", "views", "", []string{}, []string{}},
	{"grant privileges on views in foo_schema", "grant", "views", "in", []string{}, []string{"foo_schema"}},
	{"grant privileges on views in foo_schema,bar_schema", "grant", "views", "in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on views except in foo_schema", "grant", "views", "except in", []string{}, []string{"foo_schema"}},
	{"grant privileges on views except in foo_schema,bar_schema", "grant", "views", "except in", []string{}, []string{"foo_schema", "bar_schema"}},
	{"grant privileges on table foo_schema.foo_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table"}},
	{"grant on table foo_schema.foo_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table"}},
	{"grant on table foo_schema.foo_table,bar_schema.bar_table", "grant", "table", "", []string{}, []string{"foo_schema.foo_table", "bar_schema.bar_table"}},
	{"grant privileges on view foo_schema.foo_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view"}},
	{"grant on view foo_schema.foo_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view"}},
	{"grant on view foo_schema.foo_view,bar_schema.bar_view", "grant", "view", "", []string{}, []string{"foo_schema.foo_view", "bar_schema.bar_view"}},
	{"begin", "begin", "begin", "", []string{}, []string{}},
	{"begin transaction", "begin", "begin", "transaction", []string{}, []string{}},
	{"commit", "commit", "commit", "", []string{}, []string{}},
	{"commit transaction", "commit", "commit", "transaction", []string{}, []string{}},
	{"rollback", "roolback", "rollback", "", []string{}, []string{}},
	{"rollback transaction", "roolback", "rollback", "transaction", []string{}, []string{}},
}

var _ = Describe("ParseTree", func() {

	Describe("Parse", func() {

		It("parses positive tests", func() {
			for _, s := range positiveTests {
				fmt.Fprintln(GinkgoWriter, s.command)
				cmds, _, _, err := Parse(s.command)
				cmd := cmds[0]
				Expect(err).To(BeNil())
				Expect(cmd.CommandDef.Name).To(Equal(s.expectedPrimary))
				Expect(cmd.Args).To(ConsistOf(s.expectedArgs))
				Expect(cmd.ExtArgs).To(ConsistOf(s.expectedExtendedArgs))
			}
		})
	})

	Describe("ShortDesc", func() {

		It("returns short desc", func() {
			desc := ShortDesc("create")
			Expect(desc).To(Equal("Top level create command"))

			desc = ShortDesc("drop")
			Expect(desc).To(Equal("Top level drop command"))

			desc = ShortDesc("sql")
			Expect(desc).To(Equal("Run a SQL command or script"))
		})
	})

})
