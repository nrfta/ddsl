package exec

import (
	"errors"
	"fmt"
)

func (p *preprocessor) preprocessGrantOrRevoke() (int, error) {
	switch p.command.CommandDef.Name {
	case DATABASE:
		return p.preprocessKey(DATABASE_PRIVS)
	case SCHEMAS:
		return p.preprocessSchemas(SCHEMA_PRIVS)
	case SCHEMA:
		return p.preprocessSchema(SCHEMA_PRIVS)
	case TABLE:
		return p.preprocessSchemaItem(TABLE_PRIVS)
	case VIEW:
		return p.preprocessSchemaItem(VIEW_PRIVS)
	case FUNCTION:
		return p.preprocessSchemaItem(FUNCTION_PRIVS)
	case PROCEDURE:
		return p.preprocessSchemaItem(PROCEDURE_PRIVS)
	case TABLES:
		return p.preprocessSchemaItems(TABLES_PRIVS)
	case VIEWS:
		return p.preprocessSchemaItems(VIEWS_PRIVS)
	case FUNCTIONS:
		return p.preprocessSchemaItems(FUNCTIONS_PRIVS)
	case PROCEDURES:
		return p.preprocessSchemaItems(PROCEDURES_PRIVS)
	}

	return 0, errors.New("unknown command")
}

func (p *preprocessor) preprocessGrantOrRevokeKey(patternKey string, params ...interface{}) (int, error) {
	params = append(params, p.grantOrRevoke)
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
