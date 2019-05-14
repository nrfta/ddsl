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

func executeCreateOrDrop(ex *executor) error {
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

	return errors.New("unknown command")
}

func executeViews(ex *executor, schemaName string, ref *parser.Ref) error {
	if err := executeKey(ex, ref, views, schemaName, ex.createOrDrop); err != nil {
		return err
	}

	if ex.createOrDrop == drop {
		return nil
	}

	return createIndexesAndConstraints(ex, schemaName, views, ref)
}

func executeTables(ex *executor, schemaName string, ref *parser.Ref) error {
	if err := executeKey(ex, ref, tables, schemaName, ex.createOrDrop); err != nil {
		return err
	}

	if ex.createOrDrop == drop {
		return nil
	}

	return createIndexesAndConstraints(ex, schemaName, tables, ref)
}

func createIndexesAndConstraints(ex *executor, schemaName string, itemType string, ref *parser.Ref) error {
	items, err := ex.namesOf(itemType, schemaName, ref)
	if err != nil {
		return err
	}

	for _, item := range items {

		if err := executeConstraints(ex, schemaName, item, ref); err != nil {
			return err
		}

		if err := executeIndexes(ex, schemaName, item, ref); err != nil {
			return err
		}
	}

	return nil
}

func executeConstraints(ex *executor, schemaName string, tableName string, ref *parser.Ref) error {
	return executeKey(ex, ref, constraints, schemaName, tableName, ex.createOrDrop)
}

func executeIndexes(ex *executor, schemaName string, tableOrViewName string, ref *parser.Ref) error {
	return executeKey(ex, ref, indexes, schemaName, tableOrViewName, ex.createOrDrop)
}

func executeTableOrView(ex *executor, tableOrView string, schemaName string, tableOrViewName string, ref *parser.Ref) error {
	if err := executeKey(ex, ref, tableOrView, schemaName, tableOrViewName, ex.createOrDrop); err != nil {
		return err
	}

	if ex.createOrDrop == drop {
		return nil
	}

	if err := executeConstraints(ex, schemaName, tableOrViewName, ref); err != nil {
		return err
	}

	return executeIndexes(ex, schemaName, tableOrViewName, ref)
}

func executeSchema(ex *executor, schemaName string, ref *parser.Ref) error {
	if err := executeKey(ex, ref, schema, schemaName, ex.createOrDrop); err != nil {
		return err
	}

	if ex.createOrDrop == create {

		// create tables and views as well

		if err := executeTables(ex, schemaName, ref); err != nil {
			return err
		}
		return executeViews(ex, schemaName, ref)
	}

	return nil
}

func executeDatabase(ex *executor, ref *parser.Ref) error {
	if err := executeTopLevel(ex, database, ref); err != nil {
		return err
	}

	return nil
}

func executeTopLevel(ex *executor, itemType string, ref *parser.Ref) error {
	return executeKey(ex, ref, itemType, ex.createOrDrop)
}

func executeSchemas(ex *executor, ref *parser.Ref) error {
	if ex.createOrDrop == drop {
		if err := executeTopLevel(ex, foreignKeys, ref); err != nil {
			return err
		}
	}

	schemaNames, err := ex.getSchemaNames(ref)
	if err != nil {
		return err
	}

	for _, schemaName := range schemaNames {
		if err := executeSchema(ex, schemaName, ref); err != nil {
			return err
		}
	}

	if ex.createOrDrop == create {
		if err := executeTopLevel(ex, foreignKeys, ref); err != nil {
			return err
		}
	}

	return nil
}

func executeKey(ex *executor, ref *parser.Ref, patternKey string, params ...interface{}) error {
	path := fmt.Sprintf(pathPatterns[patternKey], params...)
	if err := ex.execute(path, ref); err != nil {
		return err
	}

	return nil

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
