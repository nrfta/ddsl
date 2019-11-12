package fluent

import (
	"fmt"
	"github.com/neighborly/ddsl/exec"
	"strings"
)

type Create struct {
	f *FluentDDSL
}

type Drop struct {
	f *FluentDDSL
}

func (o *Create) Database() error {
	o.f.addTokens(exec.DATABASE)
	return o.f.execute()
}
func (o *Drop) Database() error {
	o.f.addTokens(exec.DATABASE)
	return o.f.execute()
}

func (o *Create) Roles() error {
	o.f.addTokens(exec.ROLES)
	return o.f.execute()
}
func (o *Drop) Roles() error {
	o.f.addTokens(exec.ROLES)
	return o.f.execute()
}

func (o *Create) Extensions() error {
	o.f.addTokens(exec.EXTENSIONS)
	return o.f.execute()
}
func (o *Drop) Extensions() error {
	o.f.addTokens(exec.EXTENSIONS)
	return o.f.execute()
}

func (o *Create) ForeignKeys() error {
	o.f.addTokens(exec.FOREIGN_KEYS)
	return o.f.execute()
}
func (o *Drop) ForeignKeys() error {
	o.f.addTokens(exec.FOREIGN_KEYS)
	return o.f.execute()
}

func (o *Create) Schemas() error {
	o.f.addTokens(exec.SCHEMAS)
	return o.f.execute()
}
func (o *Drop) Schemas() error {
	o.f.addTokens(exec.SCHEMAS)
	return o.f.execute()
}

func (o *Create) SchemasExcept(schemaNames ...string) error {
	o.f.addTokens(exec.SCHEMAS, exec.EXCEPT)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) SchemasExcept(schemaNames ...string) error {
	o.f.addTokens(exec.SCHEMAS, exec.EXCEPT)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Schema(schemaNames ...string) error {
	o.f.addTokens(exec.SCHEMA)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) Schema(schemaNames ...string) error {
	o.f.addTokens(exec.SCHEMA)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Tables() error {
	o.f.addTokens(exec.TABLES)
	return o.f.execute()
}
func (o *Drop) Tables() error {
	o.f.addTokens(exec.TABLES)
	return o.f.execute()
}

func (o *Create) TablesIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) TablesIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) TablesExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) TablesExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.TABLES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Table(schemaName, tableName string) error {
	o.f.addTokens(exec.TABLE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableName))
	return o.f.execute()
}
func (o *Drop) Table(schemaName, tableName string) error {
	o.f.addTokens(exec.TABLE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableName))
	return o.f.execute()
}

func (o *Create) Views() error {
	o.f.addTokens(exec.VIEWS)
	return o.f.execute()
}
func (o *Drop) Views() error {
	o.f.addTokens(exec.VIEWS)
	return o.f.execute()
}

func (o *Create) ViewsIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) ViewsIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) ViewsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) ViewsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.VIEWS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) View(schemaName, viewName string) error {
	o.f.addTokens(exec.TABLE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, viewName))
	return o.f.execute()
}
func (o *Drop) View(schemaName, viewName string) error {
	o.f.addTokens(exec.TABLE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, viewName))
	return o.f.execute()
}

func (o *Create) Types() error {
	o.f.addTokens(exec.TYPES)
	return o.f.execute()
}
func (o *Drop) Types() error {
	o.f.addTokens(exec.TYPES)
	return o.f.execute()
}

func (o *Create) TypesIn(schemaNames ...string) error {
	o.f.addTokens(exec.TYPES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) TypesIn(schemaNames ...string) error {
	o.f.addTokens(exec.TYPES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) TypesExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.TYPES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) TypesExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.TYPES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Type(schemaName, typeName string) error {
	o.f.addTokens(exec.TYPE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, typeName))
	return o.f.execute()
}
func (o *Drop) Type(schemaName, typeName string) error {
	o.f.addTokens(exec.TYPE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, typeName))
	return o.f.execute()
}

func (o *Create) Procedures() error {
	o.f.addTokens(exec.PROCEDURES)
	return o.f.execute()
}
func (o *Drop) Procedures() error {
	o.f.addTokens(exec.PROCEDURES)
	return o.f.execute()
}

func (o *Create) ProceduresIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) ProceduresIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) ProceduresExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) ProceduresExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.PROCEDURES, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Procedure(schemaName, procedureName string) error {
	o.f.addTokens(exec.PROCEDURE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, procedureName))
	return o.f.execute()
}
func (o *Drop) Procedure(schemaName, procedureName string) error {
	o.f.addTokens(exec.PROCEDURE)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, procedureName))
	return o.f.execute()
}

func (o *Create) Functions() error {
	o.f.addTokens(exec.FUNCTIONS)
	return o.f.execute()
}
func (o *Drop) Functions() error {
	o.f.addTokens(exec.FUNCTIONS)
	return o.f.execute()
}

func (o *Create) FunctionsIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) FunctionsIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) FunctionsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}
func (o *Drop) FunctionsExceptIn(schemaNames ...string) error {
	o.f.addTokens(exec.FUNCTIONS, exec.EXCEPT, exec.IN)
	o.f.addTokens(strings.Join(schemaNames, ","))
	return o.f.execute()
}

func (o *Create) Function(schemaName, functionName string) error {
	o.f.addTokens(exec.FUNCTION)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, functionName))
	return o.f.execute()
}
func (o *Drop) Function(schemaName, functionName string) error {
	o.f.addTokens(exec.FUNCTION)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, functionName))
	return o.f.execute()
}

func (o *Create) ConstraintsOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.CONSTRAINTS, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}
func (o *Drop) ConstraintsOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.CONSTRAINTS, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}

func (o *Create) IndexesOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.INDEXES, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}
func (o *Drop) IndexesOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.INDEXES, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}

func (o *Create) TriggersOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.TRIGGERS, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}
func (o *Drop) TriggersOn(schemaName, tableOrViewName string) error {
	o.f.addTokens(exec.TRIGGERS, exec.ON)
	o.f.addTokens(fmt.Sprintf("%s.%s", schemaName, tableOrViewName))
	return o.f.execute()
}
