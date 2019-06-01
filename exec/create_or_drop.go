package exec

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

const (
	database    string = "database"
	extensions  string = "extensions"
	roles       string = "roles"
	schemas     string = "schemas"
	foreignKeys string = "foreign_keys"
	schema      string = "schema"
	tables      string = "tables"
	views       string = "views"
	table       string = "table"
	view        string = "view"
	indexes     string = "indexes"
	constraints string = "constraints"
	typeCmd     string = "type"
	types       string = "types"
)

var pathPatterns = map[string]string{
	database:    `database\.%s\..*`,
	extensions:  `extensions\.%s\..*`,
	roles:       `roles\.%s\..*`,
	schemas:     `scheams\.*`,
	foreignKeys: `foreign_keys\.%s\..*`,
	schema:      `schemas/%s/schema\.%s\..*`,
	tables:      `schemas/%s/tables/.*\.%s\..*`,
	views:       `schemas/%s/views/.*\.%s\..*`,
	types:       `schemas/%s/types/.*\.%s\..*`,
	table:       `schemas/%s/tables/%s\.%s\..*`,
	view:        `schemas/%s/views/%s\.%s..*`,
	indexes:     `schemas/%s/indexes/%s\.%s\..*`,
	constraints: `schemas/%s/constraints/%s\.%s\..*`,
	typeCmd:     `schemas/%s/types/%s\.%s\..*`,
}

func (ex *executor) executeCreateOrDrop() (int, error) {
	switch ex.command.CommandDef.Name {
	case database:
		return ex.executeDatabase()
	case schemas:
		return ex.executeSchemas()
	case extensions:
		return ex.executeTopLevel(extensions)
	case foreignKeys:
		return ex.executeTopLevel(foreignKeys)
	case roles:
		return ex.executeTopLevel(roles)
	case schema:
		return ex.executeSchema()
	case table:
		return ex.executeSchemaItem(table)
	case view:
		return ex.executeSchemaItem(view)
	case indexes:
		return ex.executeIndexes()
	case constraints:
		return ex.executeConstraints()
	case tables:
		return ex.executeTables()
	case views:
		return ex.executeViews()
	case typeCmd:
		return ex.executeSchemaItem(typeCmd)
	case types:
		return ex.executeTypes()
	}

	return 0, errors.New("unknown command")
}

func (ex *executor) executeTypes() (int, error) {
	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return 0, err
	}
	var schemaNames []string
	switch cmd {
	case "in":
		schemaNames, err = ex.getSchemaNames(opts, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, opts)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(types, schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}


func (ex *executor) executeViews() (int, error) {
	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return 0, err
	}
	var schemaNames []string
	switch cmd {
	case "in":
		schemaNames, err = ex.getSchemaNames(opts, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, opts)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(views, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == drop {
			continue
		}

		if c > 0 {
			c, err := ex.createIndexesAndConstraints(views, schemaName)
			count += c
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func (ex *executor) executeTables() (int, error) {
	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return 0, err
	}
	var schemaNames []string
	switch cmd {
	case "in":
		schemaNames, err = ex.getSchemaNames(opts, nil)
	case "except in":
		schemaNames, err = ex.getSchemaNames(nil, opts)
	default:
		schemaNames, err = ex.getSchemaNames(nil, nil)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for _, schemaName := range schemaNames {
		c, err := ex.executeCreateOrDropKey(tables, schemaName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == drop {
			continue
		}

		if c > 0 {
			c, err := ex.createIndexesAndConstraints(tables, schemaName)
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
		return count, err
	}

	return count, nil
}

func (ex *executor) executeConstraints() (int, error) {
	count := 0

	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return count, err
	}

	if cmd != "on" {
		return count, fmt.Errorf("comma-delimited list of tables is required")
	}

	for _, n := range opts {
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
	return ex.executeCreateOrDropKey(constraints, schemaName, tableName)
}

func (ex *executor) executeIndexes() (int, error) {
	count := 0

	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return count, err
	}

	if cmd != "on" {
		return count, fmt.Errorf("comma-delimited list of tables or views is required")
	}

	for _, n := range opts {
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
	return ex.executeCreateOrDropKey(indexes, schemaName, tableOrViewName)
}

func (ex *executor) executeSchemaItem(tableOrView string) (int, error) {
	count := 0

	_, opts, err := ex.command.ParseArgs()
	if err != nil {
		return count, err
	}

	if len(opts) == 0 {
		return count, fmt.Errorf("comma-delimited list of %ss must be provided", tableOrView)
	}

	for _, n := range opts {
		schemaName, tableOrViewName, err := parseSchemaItemName(n)
		if err != nil {
			return count, err
		}

		c, err := ex.executeCreateOrDropKey(tableOrView, schemaName, tableOrViewName)
		count += c
		if err != nil {
			return count, err
		}

		if ex.createOrDrop == drop {
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

func (ex *executor) executeSchema() (int, error) {
	count := 0

	_, opts, err := ex.command.ParseArgs()
	if err != nil {
		return count, err
	}

	for _, schemaName := range opts {
		c, err := ex.executeSchemaWork(schemaName)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ex *executor) executeSchemaWork(schemaName string) (int, error) {
	count, err := ex.executeCreateOrDropKey(schema, schemaName)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ex *executor) executeDatabase() (int, error) {
	count, err := ex.executeTopLevel(database)
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

	cmd, opts, err := ex.command.ParseArgs()
	if err != nil {
		return count, err
	}
	var schemaNames []string
	switch cmd {
	case "except":
		schemaNames, err = ex.getSchemaNames(nil, opts)
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

	relativePath, filePattern := getRelativePathAndFilePattern(fmt.Sprintf(pathPatterns[itemType], schemaName, create))
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
	nparts := strings.Split(item, ".")
	if len(nparts) != 2 {
		return "", "", fmt.Errorf("tables and views must be defined as <schema_name>.<table_or_view_name>")
	}

	return nparts[0], nparts[1], nil
}