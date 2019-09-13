package exec

import (
	"errors"
	"fmt"
	"path"
)

func (p *preprocessor) preprocessCreateOrDrop() (int, error) {
	switch p.command.CommandDef.Name {
	case DATABASE:
		return p.preprocessKey(DATABASE)
	case SCHEMAS:
		return p.preprocessSchemas(SCHEMA)
	case EXTENSIONS:
		return p.preprocessKey(EXTENSIONS)
	case FOREIGN_KEYS:
		return p.preprocessKey(FOREIGN_KEYS)
	case ROLES:
		return p.preprocessKey(ROLES)
	case SCHEMA:
		return p.preprocessSchema(SCHEMA)
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
		return p.preprocessSchemaItems(TABLES)
	case VIEWS:
		return p.preprocessSchemaItems(VIEWS)
	case FUNCTIONS:
		return p.preprocessSchemaItems(FUNCTIONS)
	case PROCEDURES:
		return p.preprocessSchemaItems(PROCEDURES)
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

func (p *preprocessor) preprocessIndexesConstraintsAndPrivileges(itemType string, schemaName string) (int, error) {
	relativeDir := fmt.Sprintf("schemas/%s/%s", schemaName, itemType)
	items, err := p.getSubdirectories(relativeDir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {

		base := path.Base(item)

		c, err := p.preprocessCreateOrDropKey(CONSTRAINTS, schemaName, base)
		count += c
		if err != nil {
			return count, err
		}

		c, err = p.preprocessCreateOrDropKey(INDEXES, schemaName, base)
		count += c
		if err != nil {
			return count, err
		}

		c, err = p.preprocessGrantOrRevokeKey(PRIVILEGES, schemaName, base)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
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

		c, err := p.preprocessCreateOrDropKey(CONSTRAINTS, schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (p *preprocessor) preprocessIndexes() (int, error) {
	count := 0

	if p.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables or views is required")
	}

	for _, n := range p.command.ExtArgs {
		schemaName, tableOrViewName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := p.preprocessCreateOrDropKey(INDEXES, schemaName, tableOrViewName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
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

func (p *preprocessor) preprocessCreateOrDropKey(patternKey string, params ...interface{}) (int, error) {
	params = append(params, p.createOrDrop)
	pathPattern, ok := pathPatterns[patternKey]
	if !ok {
		panic("unknown patternKey: " + patternKey)
	}
	path := fmt.Sprintf(pathPattern, params...)
	count, err := p.makeFileInstructions(path)
	if err != nil {
		return count, err
	}

	return count, nil
}
