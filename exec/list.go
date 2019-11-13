package exec

import (
	"encoding/json"
	"errors"
	"fmt"
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

const (
	OUTPUT_TEXT = "text"
	OUTPUT_CSV  = "csv"
	OUTPUT_JSON = "json"
)

func (p *preprocessor) preprocessList() (int, error) {
	dbDriver, err := dbdr.Open(p.ctx.DatbaseUrl)
	if err != nil {
		return 0, err
	}
	defer dbDriver.Close()

	p.ctx.dbDriver = dbDriver

	switch p.command.CommandDef.Name {
	case SCHEMAS:
		p.makeListInstruction(SCHEMAS, map[string]interface{}{})
		return 1, nil
	case FOREIGN_KEYS:
		p.makeListInstruction(FOREIGN_KEYS, map[string]interface{}{})
		return 1, nil
	case SCHEMA_ITEMS:
		return p.preprocessListSchemaItems(SCHEMA_ITEMS)
	case TABLES:
		return p.preprocessListSchemaItems(TABLES)
	case VIEWS:
		return p.preprocessListSchemaItems(VIEWS)
	case FUNCTIONS:
		return p.preprocessListSchemaItems(FUNCTIONS)
	case PROCEDURES:
		return p.preprocessListSchemaItems(PROCEDURES)
	case TYPES:
		return p.preprocessListSchemaItems(TYPES)
	}

	return 0, errors.New("unknown command")
}

func (p *processor) listOutput(header []string, data [][]string) error {
	switch p.ctx.OutputFormat {
	case OUTPUT_TEXT:
		return listOutputText(header, data)
	case OUTPUT_CSV:
		return listOutputCSV(header, data)
	case OUTPUT_JSON:
		return listOutputJSON(header, data)
	default:
		return fmt.Errorf("unknown output format '%s'", p.ctx.OutputFormat)
	}
}

func listOutputText(header []string, data [][]string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for i, _ := range data {
		table.Append(data[i])
	}
	table.Render()
	return nil
}

func listOutputCSV(header []string, data [][]string) error {
	var csv strings.Builder

	csv.WriteString(joinCSVValues(header))

	for i, _ := range data {
		csv.WriteString(joinCSVValues(data[i]))
	}
	os.Stdout.WriteString(csv.String())
	return nil
}

func listOutputJSON(header []string, data [][]string) error {
	if len(data) == 0 {
		os.Stdout.WriteString("[]\n")
		return nil
	}

	dataIndex := 0
	jsonArray := []map[string]interface{}{}

	for dataIndex < len(data) {
		item := map[string]interface{}{}
		for i, h := range header {
			item[h] = data[dataIndex][i]
		}
		jsonArray = append(jsonArray, item)
		dataIndex++
	}

	b, err := json.MarshalIndent(jsonArray, "", "  ")
	if err != nil {
		return err
	}

	os.Stdout.WriteString(string(b))
	return nil
}

func joinCSVValues(values []string) string {
	return fmt.Sprintf(`"%s"`, strings.Join(values, `","`)) + "\n"

}