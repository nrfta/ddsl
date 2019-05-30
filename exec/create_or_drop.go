package exec

import (
	"errors"
	"fmt"
	"github.com/neighborly/ddsl/parser"
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
)

var pathPatterns = map[string]string{
	database:    `database\.%s.*`,
	extensions:  `extensions\.%s.*`,
	roles:       `roles\.%s.*`,
	schemas:     `scheams\.*`,
	foreignKeys: `foreign_keys\.%s.*`,
	schema:      `schemas/%s/schema\.%s.*`,
	tables:      `schemas/%s/tables/.*%s.*`,
	views:       `schemas/%s/views/.*%s.*`,
	table:       `schemas/%s/tables/%s\.%s.*`,
	view:        `schemas/%s/views/%s\.%s.*`,
	indexes:     `schemas/%s/indexes/%s\.%s.*`,
	constraints: `schemas/%s/constraints/%s\.%s.*`,
}

func executeCreateOrDrop(ex *executor) (int, error) {
	var cmd *parser.Command
	if ex.createOrDrop == create {
		cmd = ex.parseTree.Create
	} else {
		cmd = ex.parseTree.Drop
	}

	switch {
	case cmd.Database != nil:
		return executeDatabase(ex, cmd.Database.Ref)
	case cmd.Schemas != nil:
		return executeSchemas(ex, cmd.Schemas.Ref)
	case cmd.Extensions != nil:
		return executeTopLevel(ex, extensions, cmd.Extensions.Ref)
	case cmd.ForeignKeys != nil:
		return executeTopLevel(ex, foreignKeys, cmd.ForeignKeys.Ref)
	case cmd.Roles != nil:
		return executeTopLevel(ex, roles, cmd.Roles.Ref)
	case cmd.Schema != nil:
		return executeSchema(ex, cmd.Schema.Name, cmd.Schema.Ref)
	case cmd.Table != nil:
		return executeTableOrView(ex, table, cmd.Table.Schema, cmd.Table.TableOrView, cmd.Table.Ref)
	case cmd.View != nil:
		return executeTableOrView(ex, view, cmd.View.Schema, cmd.View.TableOrView, cmd.View.Ref)
	case cmd.Indexes != nil:
		return executeIndexes(ex, cmd.Indexes.Schema, cmd.Indexes.TableOrView, cmd.Indexes.Ref)
	case cmd.Constraints != nil:
		return executeConstraints(ex, cmd.Constraints.Schema, cmd.Constraints.TableOrView, cmd.Constraints.Ref)
	case cmd.TablesInSchema != nil:
		return executeTables(ex, cmd.TablesInSchema.Name, cmd.TablesInSchema.Ref)
	case cmd.ViewsInSchema != nil:
		return executeViews(ex, cmd.ViewsInSchema.Name, cmd.TablesInSchema.Ref)
	}

	return 0, errors.New("unknown command")
}

func executeViews(ex *executor, schemaName string, ref *parser.Ref) (int, error) {
	count, err := executeKey(ex, ref, views, schemaName, ex.createOrDrop)
	if err != nil {
		return count, err
	}

	if ex.createOrDrop == drop {
		return count, nil
	}

	if count > 0 {
		c, err := createIndexesAndConstraints(ex, schemaName, views, ref)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func executeTables(ex *executor, schemaName string, ref *parser.Ref) (int, error) {
	count, err := executeKey(ex, ref, tables, schemaName, ex.createOrDrop)
	if err != nil {
		return count, err
	}

	if ex.createOrDrop == drop {
		return count, nil
	}

	if count > 0 {
		c, err := createIndexesAndConstraints(ex, schemaName, tables, ref)
		count += c
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func createIndexesAndConstraints(ex *executor, schemaName string, itemType string, ref *parser.Ref) (int, error) {
	items, err := ex.namesOf(itemType, schemaName, ref)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {

		c, err := executeConstraints(ex, schemaName, item, ref)
		count += c
		if err != nil {
			return count, err
		}

		c, err = executeIndexes(ex, schemaName, item, ref)
		count += c
		return count, err
	}

	return count, nil
}

func executeConstraints(ex *executor, schemaName string, tableName string, ref *parser.Ref) (int, error) {
	return executeKey(ex, ref, constraints, schemaName, tableName, ex.createOrDrop)
}

func executeIndexes(ex *executor, schemaName string, tableOrViewName string, ref *parser.Ref) (int, error) {
	return executeKey(ex, ref, indexes, schemaName, tableOrViewName, ex.createOrDrop)
}

func executeTableOrView(ex *executor, tableOrView string, schemaName string, tableOrViewName string, ref *parser.Ref) (int, error) {
	count, err := executeKey(ex, ref, tableOrView, schemaName, tableOrViewName, ex.createOrDrop)
	if err != nil {
		return count, err
	}

	if ex.createOrDrop == drop {
		return count, nil
	}

	if count > 0 {

		c, err := executeConstraints(ex, schemaName, tableOrViewName, ref)
		count += c
		if err != nil {
			return count, err
		}

		c, err = executeIndexes(ex, schemaName, tableOrViewName, ref)
		count += c
		return count, err
	}

	return count, nil
}

func executeSchema(ex *executor, schemaName string, ref *parser.Ref) (int, error) {
	count, err := executeKey(ex, ref, schema, schemaName, ex.createOrDrop)
	if err != nil {
		return count, err
	}

	if ex.createOrDrop == create && count > 0 {

		// create tables and views as well

		c, err := executeTables(ex, schemaName, ref)
		count += c
		if err != nil {
			return count, err
		}
		c, err = executeViews(ex, schemaName, ref)
		count += c
		return count, err
	}

	return count, nil
}

func executeDatabase(ex *executor, ref *parser.Ref) (int, error) {
	count, err := executeTopLevel(ex, database, ref)
	if err != nil {
		return count, err
	}

	return count, nil
}

func executeTopLevel(ex *executor, itemType string, ref *parser.Ref) (int, error) {
	return executeKey(ex, ref, itemType, ex.createOrDrop)
}

func executeSchemas(ex *executor, ref *parser.Ref) (int, error) {
	count := 0
	if ex.createOrDrop == drop {
		c, err := executeTopLevel(ex, foreignKeys, ref)
		count += c
		if err != nil {
			return count, err
		}
	}

	schemaNames, err := ex.getSchemaNames(ref)
	if err != nil {
		return count, err
	}

	for _, schemaName := range schemaNames {
		c, err := executeSchema(ex, schemaName, ref)
		count += c
		if err != nil {
			return count, err
		}
	}

	if ex.createOrDrop == create {
		c, err := executeTopLevel(ex, foreignKeys, ref)
		count += c
		return count, err
	}

	return count, nil
}

func executeKey(ex *executor, ref *parser.Ref, patternKey string, params ...interface{}) (int, error) {
	path := fmt.Sprintf(pathPatterns[patternKey], params...)
	count, err := ex.execute(path, ref)
	if err != nil {
		return count, err
	}

	return count, nil

}

func (ex *executor) namesOf(itemType string, schemaName string, ref *parser.Ref) ([]string, error) {
	if err := ex.getSourceDriver(ref); err != nil {
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
