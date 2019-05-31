package repl

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/neighborly/ddsl/parser"
	"sort"
	"strings"
)

type followingSuggestion struct {
	text string
	cmdDef *parser.CommandDef
}


func completer(d prompt.Document) []prompt.Suggest {
	command := d.TextBeforeCursor()
	args := strings.Split(command, " ")

	if len(command) == 0 {
		return suggestFromCommandDefs(parser.ParseTree.CommandDefs)
	}

	if len(args) <= 2 {
		cmdDef, err := parser.Parse(command)
		if err != nil {
			return []prompt.Suggest{}
		}

		return suggestFromCommandDefs(cmdDef.CommandDefs)
	}

	rootCmd, err := parser.Parse(args[0])
	if err != nil {
		return []prompt.Suggest{}
	}

	cmdDef, err := parser.Parse(strings.Join(args[:2], " "))
	if err != nil {
		return []prompt.Suggest{}
	}

	i := 2
	nextCmdDef := cmdDef
	partial := ""
	for i < len(args) {
		w := args[i]
		c, ok := nextCmdDef.CommandDefs[w]
		i++
		// unrecognized and not the last arg
		if !ok && i < len(args) {
			return []prompt.Suggest{}
		}
		if !ok && i == len(args) {
			partial = w
		}
		nextCmdDef = c
	}

	followingCmds := getFollowingSuggestions(nextCmdDef, "")
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
	if len(nextCmdDef.ArgDefs) > 0 {
		suggestedArgs, err := suggestArgs(rootCmd, nextCmdDef)
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
		if len(subCmds) == 0 || !c.Optional {
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

