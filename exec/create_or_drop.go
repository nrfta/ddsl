package exec

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

const (
	DATABASE     string = "database"
	EXTENSIONS   string = "extensions"
	ROLES        string = "roles"
	SCHEMAS      string = "schemas"
	FOREIGN_KEYS string = "foreign_keys"
	SCHEMA       string = "schema"
	TABLES       string = "tables"
	VIEWS        string = "views"
	TABLE        string = "table"
	VIEW         string = "view"
	INDEXES      string = "indexes"
	CONSTRAINTS  string = "constraints"
	TYPE         string = "type"
	TYPES        string = "types"
)

var pathPatterns = map[string]string{
	DATABASE:     `database\.%s\..*`,
	ROLES:        `roles\.%s\..*`,
	SCHEMAS:      `scheams\.*`,
	FOREIGN_KEYS: `foreign_keys\.%s\..*`,
	SCHEMA:       `schemas/%s/schema\.%s\..*`,
	EXTENSIONS:   `schemas/%s/extensions\.%s\..*`,
	TABLES:       `schemas/%s/tables/.*\.%s\..*`,
	VIEWS:        `schemas/%s/views/.*\.%s\..*`,
	TYPES:        `schemas/%s/types/.*\.%s\..*`,
	TABLE:        `schemas/%s/tables/%s\.%s\..*`,
	VIEW:         `schemas/%s/views/%s\.%s..*`,
	INDEXES:      `schemas/%s/indexes/%s\.%s\..*`,
	CONSTRAINTS:  `schemas/%s/constraints/%s\.%s\..*`,
	TYPE:         `schemas/%s/types/%s\.%s\..*`,
}

func (ex *executor) executeCreateOrDrop() (int, error) {
	switch ex.command.CommandDef.Name {
	case DATABASE:
		return ex.executeDatabase()
	case SCHEMAS:
		return ex.executeSchemas()
	case EXTENSIONS:
		return ex.executeExtensions()
	case FOREIGN_KEYS:
		return ex.executeTopLevel(FOREIGN_KEYS)
	case ROLES:
		return ex.executeTopLevel(ROLES)
	case SCHEMA:
		return ex.executeSchema()
	case TABLE:
		return ex.executeSchemaItem(TABLE)
	case VIEW:
		return ex.executeSchemaItem(VIEW)
	case INDEXES:
		return ex.executeIndexes()
	case CONSTRAINTS:
		return ex.executeConstraints()
	case TABLES:
		return ex.executeTables()
	case VIEWS:
		return ex.executeViews()
	case TYPE:
		return ex.executeSchemaItem(TYPE)
	case TYPES:
		return ex.executeTypes()
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) executeTypes() (int, error) {

	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(TYPES, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}


func (ex *executor) executeViews() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(VIEWS, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := ex.createIndexesAndConstraints(VIEWS, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (ex *executor) executeTables() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(TABLES, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == DROP {
			continue
		}

		if c > 0 {
			c, err := ex.createIndexesAndConstraints(TABLES, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (ex *executor) createIndexesAndConstraints(itemType string, schemaName string) (int, error) {
	items, err := ex.namesOf(itemType, schemaName)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {

		c, err := ex.executeConstraintsWork(schemaName, item)
		count += c
		if err != nil {
			return count, err
		}

		c, err = ex.executeIndexesWork(schemaName, item)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeConstraints() (int, error) {
	count := 0

	if ex.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables is required")
	}

	for _, n := range ex.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := ex.executeConstraintsWork(schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}
func (ex *executor) executeConstraintsWork(schemaName, tableName string) (int, error) {
	return ex.executeCreateOrDropKey(CONSTRAINTS, schemaName, tableName)
}

func (ex *executor) executeIndexes() (int, error) {
	count := 0

	if ex.command.Clause != "on" {
		return count, fmt.Errorf("comma-delimited list of tables or views is required")
	}

	for _, n := range ex.command.ExtArgs {
		schemaName, tableName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := ex.executeIndexesWork(schemaName, tableName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}
func (ex *executor) executeIndexesWork(schemaName, tableOrViewName string) (int, error) {
	return ex.executeCreateOrDropKey(INDEXES, schemaName, tableOrViewName)
}

func (ex *executor) executeSchemaItem(tableOrView string) (int, error) {
	count := 0

	if len(ex.command.ExtArgs) == 0 {
		return count, fmt.Errorf("comma-delimited list of %ss must be provided", tableOrView)
	}

	for _, n := range ex.command.ExtArgs {
		schemaName, tableOrViewName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := ex.executeCreateOrDropKey(tableOrView, schemaName, tableOrViewName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == DROP {
			continue
		}

		// when creating a table, also create its constraints and indexes
		if c > 0 {
			c, err := ex.executeConstraintsWork(schemaName, tableOrViewName)
			count += c
			if err != nil {
				return count, err
			}

			c, err = ex.executeIndexesWork(schemaName, tableOrViewName)
			count += c
			return count, err
		}
	}
	return count, nil
}

func (ex *executor) executeExtensions() (int, error) {
	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "in":
		schemaNames, err = ex.getSchemaNames(ex.command.ExtArgs, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(EXTENSIONS, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSchema() (int, error) {
	count := 0

	for _, schemaName := range ex.command.ExtArgs {
		c, err := ex.executeSchemaWork(schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSchemaWork(schemaName string) (int, error) {
	count, err := ex.executeCreateOrDropKey(SCHEMA, schemaName)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ex *executor) executeDatabase() (int, error) {
	count, err := ex.executeTopLevel(DATABASE)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ex *executor) executeTopLevel(itemType string) (int, error) {
	return ex.executeCreateOrDropKey(itemType)
}

func (ex *executor) executeSchemas() (int, error) {
	count := 0

	var schemaNames []string
	var err error
	switch ex.command.Clause {
	case "except":
		schemaNames, err = ex.getSchemaNames(nil, ex.command.ExtArgs)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return count, err
	}

	for _, schemaName := range schemaNames {
		c, err := ex.executeSchemaWork(schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeCreateOrDropKey(patternKey string, params ...interface{}) (int, error) {
	params = append(params, ex.createOrDrop)
	path := fmt.Sprintf(pathPatterns[patternKey], params...)
	count, err := ex.execute(path)
	if err != nil {
		return count, err
	}

	return count, nil

}

func (ex *executor) namesOf(itemType string, schemaName string) ([]string, error) {
	if err := ex.getSourceDriver(ex.command.Ref); err != nil {
		return nil, err
	}
	defer ex.sourceDriver.Close()

	relativePath, filePattern := getRelativePathAndFilePattern(fmt.Sprintf(pathPatterns[itemType], schemaName, CREATE))
	readers, err := ex.sourceDriver.ReadFiles(relativePath, filePattern)
	if err != nil {
		return nil, err
	}

	items := []string{}
	for _, fr := range readers {
		base := path.Base(fr.FilePath)
		i := strings.Index(base, ".")
		items = append(items, base[:i])
	}

	return items, nil
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