package fluent

import (
	"github.com/neighborly/ddsl/exec"
	"strings"
)

type Grant struct {
	f *FluentDDSL
}

type Revoke struct {
	f *FluentDDSL
}

type PrivlegesOn struct {
	f *FluentDDSL
}

func (o *Grant) PrivilegesOn() *PrivlegesOn {
	o.f.addTokens(exec.PRIVILEGES, exec.ON)
	return &PrivlegesOn{o.f}
}
func (o *Revoke) PrivilegesOn() *PrivlegesOn {
	o.f.addTokens(exec.PRIVILEGES, exec.ON)
	return &PrivlegesOn{o.f}
}

func (o *PrivlegesOn) Database() error {
	o.f.addTokens(exec.DATABASE)
	return o.f.execute()
}

func (o *PrivlegesOn) Schemas() error {
	o.f.addTokens(exec.SCHEMAS)
	return o.f.execute()
}

func (o *PrivlegesOn) Schema(schemaNames ...string) error {
	o.f.addTokens(exec.SCHEMA)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) Tables() error {
	o.f.addTokens(exec.TABLES)
	return o.f.execute()
}

func (o *PrivlegesOn) TablesIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) TablesExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) Views() error {
	o.f.addTokens(exec.VIEWS)
	return o.f.execute()
}

func (o *PrivlegesOn) ViewsIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) ViewsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) Procedures() error {
	o.f.addTokens(exec.PROCEDURES)
	return o.f.execute()
}

func (o *PrivlegesOn) ProceduresIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) ProceduresExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) Functions() error {
	o.f.addTokens(exec.FUNCTIONS)
	return o.f.execute()
}

func (o *PrivlegesOn) FunctionsIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *PrivlegesOn) FunctionsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

