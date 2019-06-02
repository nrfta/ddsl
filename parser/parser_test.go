package parser

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type posTestSpec struct {
	command string
	expectedRoot string
	expectedPrimary string
	expectedClause string
	expectedArgs []string
	expectedExtendedArgs []string
}

type negTestSpec struct {
	command string
	message string
}

var positiveTests = []posTestSpec{
	posTestSpec{"create database","create", "database", "", []string{}, []string{}},
	posTestSpec{"create DATABASE","create", "database", "", []string{}, []string{}},
	posTestSpec{"CREATE DATABASE","create", "database", "", []string{}, []string{}},
	posTestSpec{"CREATE database","create", "database", "", []string{}, []string{}},
	posTestSpec{"cReAtE dAtAbAsE","create", "database", "", []string{}, []string{}},
	posTestSpec{"create roles","create", "roles", "",[]string{}, []string{}},
	posTestSpec{"create foreign-keys", "create","foreign-keys", "",[]string{}, []string{}},
	posTestSpec{"create schemas", "create","schemas", "",[]string{}, []string{}},
	posTestSpec{"create schemas except foo_schema", "create","schemas", "except",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create schemas except foo_schema,bar_schema", "create","schemas", "except",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create schema foo_schema", "create","schema", "",[]string{},[]string{"foo_schema"}},
	posTestSpec{"CREATE SCHEMA FOO_SCHEMA", "create","schema", "",[]string{},[]string{"FOO_SCHEMA"}},
	posTestSpec{"create schema foo_schema,bar_schema", "create","schema", "",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create extensions", "create","extensions", "",[]string{}, []string{}},
	posTestSpec{"create extensions in foo_schema", "create","extensions", "in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create extensions in foo_schema,bar_schema", "create","extensions", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create extensions except in foo_schema", "create","extensions", "except in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create extensions except in foo_schema,bar_schema", "create","extensions", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create tables", "create","tables", "",[]string{}, []string{}},
	posTestSpec{"create tables in foo_schema", "create","tables", "in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create tables in foo_schema,bar_schema", "create","tables", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create tables except in foo_schema", "create","tables", "except in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create tables except in foo_schema,bar_schema", "create","tables", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create views", "create","views", "",[]string{}, []string{}},
	posTestSpec{"create views in foo_schema", "create","views", "in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create views in foo_schema,bar_schema", "create","views", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create views except in foo_schema", "create","views", "except in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create views except in foo_schema,bar_schema", "create","views", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create types", "create","types", "",[]string{}, []string{}},
	posTestSpec{"create types in foo_schema", "create","types", "in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create types in foo_schema,bar_schema", "create","types", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create types except in foo_schema", "create","types", "except in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"create types except in foo_schema,bar_schema", "create","types", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"create table foo_schema.foo_table", "create","table", "", []string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"create table foo_schema.foo_table,bar_schema.bar_table", "create","table", "", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	posTestSpec{"create view foo_schema.foo_view", "create","view", "", []string{},[]string{"foo_schema.foo_view"}},
	posTestSpec{"create view foo_schema.foo_view,bar_schema.bar_view", "create","view", "", []string{},[]string{"foo_schema.foo_view","bar_schema.bar_view"}},
	posTestSpec{"create type foo_schema.foo_type", "create","type", "", []string{},[]string{"foo_schema.foo_type"}},
	posTestSpec{"create type foo_schema.foo_type,bar_schema.bar_type", "create","type", "", []string{},[]string{"foo_schema.foo_type","bar_schema.bar_type"}},
	posTestSpec{"create constraints on foo_schema.foo_table", "create","constraints", "on", []string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"create constraints on foo_schema.foo_table,bar_schema.bar_table", "create","constraints", "on", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	posTestSpec{"create indexes on foo_schema.foo_table", "create","indexes", "on", []string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"create indexes on foo_schema.foo_table,bar_schema.bar_table", "create","indexes", "on", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	posTestSpec{"migrate up 1", "migrate","up", "", []string{},[]string{"1"}},
	posTestSpec{"migrate down 1", "migrate","down", "", []string{},[]string{"1"}},
	posTestSpec{"migrate top", "migrate","top", "", []string{}, []string{}},
	posTestSpec{"migrate bottom", "migrate","bottom", "", []string{}, []string{}},
	posTestSpec{`seed cmd "python3 foo.py"`, "seed","cmd","", []string{},[]string{"python3 foo.py"}},
	posTestSpec{"seed database","create","database", "", []string{}, []string{}},
	posTestSpec{"seed database with foo_seed", "seed","database", "with", []string{},[]string{"foo_seed"}},
	posTestSpec{"seed database with foo_seed,bar_seed", "seed","database", "with", []string{},[]string{"foo_seed","bar_seed"}},
	posTestSpec{"seed database without foo_seed", "seed","database", "without", []string{},[]string{"foo_seed"}},
	posTestSpec{"seed database without foo_seed,bar_seed", "seed","database", "without", []string{},[]string{"foo_seed","bar_seed"}},
	posTestSpec{"seed schemas", "seed","schemas", "",[]string{}, []string{}},
	posTestSpec{"seed schemas except foo_schema", "seed","schemas", "except",[]string{},[]string{"foo_schema"}},
	posTestSpec{"seed schemas except foo_schema,bar_schema", "seed","schemas", "except",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"seed schema foo_schema", "seed","schema", "",[]string{"foo_schema"},[]string{}},
	posTestSpec{"seed schema foo_schema with foo_seed", "seed","schema", "with",[]string{"foo_schema"}, []string{"foo_seed"}},
	posTestSpec{"SEED SCHEMA FOO_SCHEMA WITH FOO_SEED", "seed","schema", "with",[]string{"FOO_SCHEMA"}, []string{"FOO_SEED"}},
	posTestSpec{"seed schema foo_schema with foo_seed,bar_seed", "seed","schema", "with",[]string{"foo_schema"}, []string{"foo_seed","bar_seed"}},
	posTestSpec{"seed schema foo_schema without foo_seed", "seed","schema", "without",[]string{"foo_schema"}, []string{"foo_seed"}},
	posTestSpec{"seed schema foo_schema without foo_seed,bar_seed", "seed","schema", "without",[]string{"foo_schema"}, []string{"foo_seed","bar_seed"}},
	posTestSpec{"seed table foo_schema.foo_table", "seed","table", "",[]string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"seed table foo_schema.foo_table,bar_schema.bar_table", "seed","table", "",[]string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	posTestSpec{"seed tables", "seed","tables", "",[]string{}, []string{}},
	posTestSpec{"seed tables in foo_schema", "seed","tables", "in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"seed tables in foo_schema,bar_schema", "seed","tables", "in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"seed tables except in foo_schema", "seed","tables", "except in",[]string{},[]string{"foo_schema"}},
	posTestSpec{"seed tables except in foo_schema,bar_schema", "seed","tables", "except in",[]string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on database", "grant","database", "", []string{}, []string{}},
	posTestSpec{"grant on database", "grant","database", "", []string{}, []string{}},
	posTestSpec{"grant on database", "grant","database", "", []string{}, []string{}},
	posTestSpec{"grant privileges on schemas", "grant","schemas", "", []string{}, []string{}},
	posTestSpec{"grant on schemas", "grant","schemas", "", []string{}, []string{}},
	posTestSpec{"grant privileges on schemas except foo_schema", "grant","schemas", "except", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on schemas except foo_schema,bar_schema", "grant","schemas", "except", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on schema foo_schema", "grant","schema", "", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on schema foo_schema,bar_schema", "grant","schema", "", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on tables", "grant","tables", "", []string{}, []string{}},
	posTestSpec{"grant on tables", "grant","tables", "", []string{}, []string{}},
	posTestSpec{"grant privileges on tables in foo_schema", "grant","tables", "in", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on tables in foo_schema,bar_schema", "grant","tables", "in", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on tables except in foo_schema", "grant","tables", "except in", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on tables except in foo_schema,bar_schema", "grant","tables", "except in", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on views", "grant","views", "", []string{}, []string{}},
	posTestSpec{"grant on views", "grant","views", "", []string{}, []string{}},
	posTestSpec{"grant privileges on views in foo_schema", "grant","views", "in", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on views in foo_schema,bar_schema", "grant","views", "in", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on views except in foo_schema", "grant","views", "except in", []string{},[]string{"foo_schema"}},
	posTestSpec{"grant privileges on views except in foo_schema,bar_schema", "grant","views", "except in", []string{},[]string{"foo_schema","bar_schema"}},
	posTestSpec{"grant privileges on table foo_schema.foo_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"grant on table foo_schema.foo_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table"}},
	posTestSpec{"grant on table foo_schema.foo_table,bar_schema.bar_table", "grant","table", "", []string{},[]string{"foo_schema.foo_table","bar_schema.bar_table"}},
	posTestSpec{"grant privileges on view foo_schema.foo_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view"}},
	posTestSpec{"grant on view foo_schema.foo_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view"}},
	posTestSpec{"grant on view foo_schema.foo_view,bar_schema.bar_view", "grant","view", "", []string{},[]string{"foo_schema.foo_view","bar_schema.bar_view"}},
	posTestSpec{"begin","begin","begin", "", []string{},[]string{}},
	posTestSpec{"begin transaction","begin","begin", "transaction", []string{},[]string{}},
	posTestSpec{"commit","commit","commit", "", []string{},[]string{}},
	posTestSpec{"commit transaction","commit","commit", "transaction", []string{},[]string{}},
	posTestSpec{"rollback","roolback","rollback", "", []string{},[]string{}},
	posTestSpec{"rollback transaction","roolback","rollback", "transaction", []string{},[]string{}},
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
