package exec

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

func (p *preprocessor) preprocessCreateOrDrop() (int, error) {
	switch p.command.CommandDef.Name {
	case DATABASE:
		return p.preprocessDatabase()
	case SCHEMAS:
		return p.preprocessSchemas()
	case EXTENSIONS:
		return p.preprocessExtensions()
	case KEYS:
		return p.preprocessTopLevel(FOREIGN_KEYS)
	case ROLES:
		return p.preprocessTopLevel(ROLES)
	case SCHEMA:
		return p.preprocessSchema()
	case TABLE:
		return p.preprocessSchemaItem(TABLE)
	case VIEW:
		return p.preprocessSchemaItem(VIEW)
	case FUNCTION:
		return p.preprocessSchemaItem(FUNCTION)
	case PROCEDURE:
		return p.preprocessSchemaItem(PROCEDURE)
	case INDEXES:
		return p.preprocessIndexes()
	case CONSTRAINTS:
		return p.preprocessConstraints()
	case TRIGGERS:
		return p.preprocessTriggers()
	case TABLES:
		return p.preprocessTables()
	case VIEWS:
		return p.preprocessViews()
	case FUNCTIONS:
		return p.preprocessFunctions()
	case PROCEDURES:
		return p.preprocessProcedures()
	case TYPE:
		return p.preprocessSchemaItem(TYPE)
	case TYPES:
		return p.preprocessTypes()
	}

	return 0, errors.New("unknown command")
}

func (p *preprocessor) preprocessTypes() (int, error) {

	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(TYPES, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessViews() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(VIEWS, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := p.preprocessIndexesConstraintsAndPrivileges(VIEWS, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessTables() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(TABLES, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := p.preprocessIndexesConstraintsAndPrivileges(TABLES, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessFunctions() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(FUNCTIONS, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := p.preprocessIndexesConstraintsAndPrivileges(FUNCTIONS, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessProcedures() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(PROCEDURES, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := p.preprocessIndexesConstraintsAndPrivileges(PROCEDURES, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessIndexesConstraintsAndPrivileges(itemType string, schemaName string) (int, error) {
	relativeDir := fmt.Sprintf("schemas/%s/%s", schemaName, itemType)
	items, err := p.getSubdirectories(relativeDir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {

		base := path.Base(item)

		c, err := p.preprocessConstraintsWork(schemaName, base)
		count += c
		if err != nil {
			return count, err
		}

		c, err = p.preprocessIndexesWork(schemaName, base)
		count += c
		if err != nil {
			return count, err
		}

		c, err = p.preprocessPrivilegesWork(schemaName, base)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessPrivilegesWork(schemaName, tableName string) (int, error) {
	return p.preprocessGrantOrRevokeKey(PRIVILEGES, schemaName, tableName)
}

func (p *preprocessor) preprocessConstraints() (int, error) {
	count := 0

	if p.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables is required")
	}

	for _, n := range p.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessConstraintsWork(schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}
func (p *preprocessor) preprocessConstraintsWork(schemaName, tableName string) (int, error) {
	return p.preprocessCreateOrDropKey(CONSTRAINTS, schemaName, tableName)
}

func (p *preprocessor) preprocessIndexes() (int, error) {
	count := 0

	if p.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables or views is required")
	}

	for _, n := range p.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessIndexesWork(schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}
func (p *preprocessor) preprocessIndexesWork(schemaName, tableOrViewName string) (int, error) {
	return p.preprocessCreateOrDropKey(INDEXES, schemaName, tableOrViewName)
}

func (p *preprocessor) preprocessTriggers() (int, error) {
	count := 0

	if p.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables or views is required")
	}

	for _, n := range p.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessCreateOrDropKey(TRIGGERS, schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSchemaItem(itemType string) (int, error) {
	count := 0

	if len(p.command.ExtArgs) == 0 {
		return count, fmt.Errorf("comma-delimited list of %ss must be provided", itemType)
	}

	for _, n := range p.command.ExtArgs {
		schemaName, tableOrViewName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessCreateOrDropKey(itemType, schemaName, tableOrViewName)
		count += c
		if err != nil {
			return count, err
		}

		if p.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			// when creating a table, also create its constraints
			if itemType == TABLE {
				c, err := p.preprocessConstraintsWork(schemaName, tableOrViewName)
				count += c
				if err != nil {
					return count, err
				}
			}

			// when creating a table or view, also create its indexes
			if itemType == TABLE || itemType == VIEW {
				c, err = p.preprocessIndexesWork(schemaName, tableOrViewName)
				count += c
				if err != nil {
					return count, err
				}
			}

			// when creating any schema item, also grant its privileges
			c, err = p.preprocessPrivilegesWork(schemaName, tableOrViewName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}
	return count, nil
}

func (p *preprocessor) preprocessExtensions() (int, error) {
	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "in":
		schemaNames, err = p.getSchemaNames(p.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := p.preprocessCreateOrDropKey(EXTENSIONS, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSchema() (int, error) {
	count := 0

	for _, schemaName := range p.command.ExtArgs {
		c, err := p.preprocessSchemaWork(schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessSchemaWork(schemaName string) (int, error) {
	count, err := p.preprocessCreateOrDropKey(SCHEMA, schemaName)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *preprocessor) preprocessDatabase() (int, error) {
	count, err := p.preprocessTopLevel(DATABASE)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *preprocessor) preprocessTopLevel(itemType string) (int, error) {
	return p.preprocessCreateOrDropKey(itemType)
}

func (p *preprocessor) preprocessSchemas() (int, error) {
	count := 0

	var schemaNames []string
	var err error
	switch p.command.Clause {
	case "except":
		schemaNames, err = p.getSchemaNames(nil, p.command.ExtArgs)
	default:
		schemaNames, err = p.getSchemaNames(nil, nil)
	}
	if err != nil {
		return count, err
	}

	for _, schemaName := range schemaNames {
		c, err := p.preprocessSchemaWork(schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessCreateOrDropKey(patternKey string, params ...interface{}) (int, error) {
	params = append(params, p.createOrDrop)
	path := fmt.Sprintf(pathPatterns[patternKey], params...)
	count, err := p.makeFileInstructions(path)
	if err != nil {
		return count, err
	}

	return count, nil

}

func parseSchemaItemName(item string) (schemaName string, tableOrViewName string, err error) {
	if len(item) == 0 {
		return "", "", fmt.Errorf("empty table or view name provided; check for trailing comma or space after comma in list arg")
	}
	nparts := strings.Split(item, ".")
	if len(nparts) != 2 {
		return "", "", fmt.Errorf("tables and views must be defined as <schema_name>.<table_or_view_name>")
	}

	return nparts[0], nparts[1], nil
}
