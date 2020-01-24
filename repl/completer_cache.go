package repl

import "github.com/nrfta/ddsl/exec"

type completerCache struct {
	repo      string
	url       string
	context   *exec.Context
	dbSchemas []string
	dbTables  []string
	dbViews   []string
	dbTypes   []string

}
var cache *completerCache


func initializeCache(ctx *exec.Context) {
	cache = &completerCache{
		context:   ctx,
	}
}

func invalidateCache() {
	initializeCache(cache.context)
}

func (c *completerCache) getDatabaseSchemas() ([]string, error) {
	if c.dbSchemas == nil {
		schemas, err := c.context.GetDatabaseSchemas()
		if err != nil {
			return nil, err
		}
		c.dbSchemas = schemas
	}
	return c.dbSchemas, nil
}

func (c *completerCache) getDatabaseTables() ([]string, error) {
	if c.dbTables == nil {
		tables, err := c.context.GetDatabaseTables()
		if err != nil {
			return nil, err
		}
		c.dbTables = tables
	}
	return c.dbTables, nil
}

func (c *completerCache) getDatabaseViews() ([]string, error) {
	if c.dbViews == nil {
		views, err := c.context.GetDatabaseViews()
		if err != nil {
			return nil, err
		}
		c.dbViews = views
	}
	return c.dbViews, nil
}

func (c *completerCache) getDatabaseTypes() ([]string, error) {
	if c.dbTypes == nil {
		types, err := c.context.GetDatabaseTypes()
		if err != nil {
			return nil, err
		}
		c.dbTypes = types
	}
	return c.dbTypes, nil
}

func (c *completerCache) getDatabaseSeeds() ([]string, error) {
	return []string{}, nil
}

func (c *completerCache) getSchemaSeeds() ([]string, error) {
	return []string{}, nil
}

func (c *completerCache) getTableSeeds() ([]string, error) {
	return []string{}, nil
}



