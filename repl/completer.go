package repl

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/nrfta/ddsl/parser"
	"sort"
	"strings"
)

type followingSuggestion struct {
	text string
	cmdDef *parser.CommandDef
}


func completer(d prompt.Document) []prompt.Suggest {
	command := d.TextBeforeCursor()

	if len(command) == 0 {
		return suggestFromCommandDefs(parser.ParseTree.CommandDefs)
	}

	cmd, remainder, _ := parser.TryParse(command)
	if cmd == nil {
		return []prompt.Suggest{}
	}

	if cmd.CommandDef.Name != d.GetWordBeforeCursor() {
		return []prompt.Suggest{}
	}

	if len(remainder) == 0 && len(cmd.CommandDef.ArgDefs) == 0 {
		return []prompt.Suggest{}
	}

	partial := ""
	if len(remainder) > 0 {
		partial = remainder[len(remainder)-1]
	}

	followingCmds := getFollowingSuggestions(cmd.CommandDef, "")
	suggestions = []prompt.Suggest{}
	if len(followingCmds) > 0 {
		for _, c := range followingCmds {
			cparts := strings.Split(c.text, ",")
			suggestions = append(suggestions, prompt.Suggest{ Text: cparts[0], Description: cparts[1]})
		}
		if len(partial) > 0 {
			return prompt.FilterHasPrefix(suggestions, partial, true)
		}
	}

	//
	if len(cmd.CommandDef.ArgDefs) > 0 {
		suggestedArgs, err := suggestArgs(cmd.RootDef, cmd.CommandDef)
		if err != nil {
			return []prompt.Suggest{}
		}
		for _, a := range suggestedArgs {
			suggestions = append(suggestions, prompt.Suggest{Text: a})
		}
	}

	if len(partial) > 0 {
		return prompt.FilterHasPrefix(suggestions, partial, true)
	}

	return suggestions
}

func suggestArgs(rootCmd, nextCmd *parser.CommandDef) ([]string, error) {
	result := []string{}
	for _, a := range nextCmd.ArgDefs {
		switch a.Name {
		case "include_schemas", "exclude_schemas":
			schemas, err := cache.getDatabaseSchemas()
			if err != nil {
				return nil, err
			}
			result = append(result, schemas...)
		case "include_tables", "exclude_tables":
			tables, err := cache.getDatabaseTables()
			if err != nil {
				return nil, err
			}
			result = append(result, tables...)
		case "include_views", "exclude_views":
			views, err := cache.getDatabaseViews()
			if err != nil {
				return nil, err
			}
			result = append(result, views...)
		case "include_types", "exclude_types":
			types, err := cache.getDatabaseTypes()
			if err != nil {
				return nil, err
			}
			result = append(result, types...)
		case "database_seeds":
			seeds, err := cache.getDatabaseSeeds()
			if err != nil {
				return nil, err
			}
			result = append(result, seeds...)
		case "schema_seeds":
			seeds, err := cache.getSchemaSeeds()
			if err != nil {
				return nil, err
			}
			result = append(result, seeds...)
		case "table_seeds":
			seeds, err := cache.getTableSeeds()
			if err != nil {
				return nil, err
			}
			result = append(result, seeds...)
		}
	}
	return result, nil
}

func suggestFromCommandDefs(commandDefs map[string]*parser.CommandDef) []prompt.Suggest {
	result := []prompt.Suggest{}
	for _, cmdDef := range commandDefs {
		result = append(result, prompt.Suggest{
			Text:        cmdDef.Name,
			Description: cmdDef.ShortDesc,
		})
	}

	sort.Slice(result, func(i,j int) bool { return result[i].Text > result[j].Text })

	return result
}



func getFollowingSuggestions(cmdDef *parser.CommandDef, prefix string) []*followingSuggestion {
	result := []*followingSuggestion{}
	for _, c := range cmdDef.CommandDefs {
		cmd := fmt.Sprintf("%s %s", prefix, c.Name)
		subCmds := getFollowingSuggestions(c, cmd)
		if len(subCmds) == 0 || !c.IsOptional() {
			result = append(result, &followingSuggestion{
				text: fmt.Sprintf("%s,%s",cmd, c.ShortDesc),
				cmdDef: c,
			})
		} else {
			result = append(result, subCmds...)
		}
	}

	return result
}

