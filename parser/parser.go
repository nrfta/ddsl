package parser

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"strings"
)

type Command struct {
	Text       string
	CommandDef *CommandDef
	RootDef    *CommandDef
	Clause     string
	Args       []string
	ExtArgs    []string
	Ref        *string
}

var shellParser *shellwords.Parser

func init() {
	shellParser = shellwords.NewParser()
	shellParser.ParseEnv = true
	shellParser.ParseBacktick = true
}

func Parse(text string) (cmds []*Command, hasTx bool, hasDB bool, err error) {
	cmds = []*Command{}
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank and comment line
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// ignore everything including and after #
		i := strings.Index(line, "#")
		if i > -1 {
			line = strings.TrimSpace(line[:i])
		}

		cs := strings.Split(line, ";")
		for _, line = range cs {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			cmd, err := parse(line)
			if err != nil {
				return nil, false, false, err
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
				rootName := cmd.RootDef.Name
				if rootName == "begin" {
					hasTx = true
				}
				if (rootName == "create" || rootName == "drop") && cmd.CommandDef.Name == "database" {
					hasDB = true
				}
			}
		}
	}

	if hasTx && hasDB {
		return cmds, hasTx, hasDB, fmt.Errorf("cannot create or drop database in transaction")
	}

	return
}

func parse(command string) (*Command, error) {
	cmd, remainder, err := TryParse(command)
	if !cmd.CommandDef.IsPrimary() {
		return cmd, fmt.Errorf("primary command token not found")
	}

	if err != nil {
		return cmd, err
	}

	clause, extArgs, err := cmd.parseRemainder(remainder)
	if err != nil {
		return cmd, err
	}
	cmd.Clause = clause
	cmd.ExtArgs = extArgs
	cmd.Text = command

	return cmd, err
}

// TryParse parses the given partial command and returns the deepest associated `Command`.
// This is used for repl and commandline completions.
func TryParse(command string) (cmd *Command, remainder []string, err error) {
	if len(command) == 0 {
		return nil, nil, fmt.Errorf("no command was provided")
	}

	tokens, err := shellParser.Parse(command)
	if err != nil {
		return nil, nil, err
	}

	cmdDefs := ParseTree.CommandDefs
	args := []string{}
	remainder = []string{}
	err = fmt.Errorf("syntax error in '%s'", command)
	var cmdDef *CommandDef

	for i, token := range tokens {
		next, ok := cmdDefs[strings.ToLower(token)]
		if ok {
			tokenIndex := i
			if next.HasExtArgs() {
				for a := 0; a < len(next.ArgDefs); a++ {
					if tokenIndex+1 < len(tokens) {
						tokenIndex++
						args = append(args, tokens[tokenIndex])
					} else {
						break
					}
				}
			}
			if len(tokens[tokenIndex:]) > 1 {
				remainder = tokens[tokenIndex+1:]
			} else {
				remainder = []string{}
			}
			cmdDef = next
			cmdDefs = next.CommandDefs
			if cmdDef.IsPrimary() {
				err = nil
				break
			}
		} else {
			if len(cmdDef.ArgDefs) > 0 {
				// token is not a command, so assume it's an arg,
				// do not advance down the parse tree
				args = append(args, token)
			} else {
				next, _ = cmdDef.skipOptionalTo(token)
				if next == nil {
					err = fmt.Errorf("syntax error at '%s'", token)
					remainder = tokens[i:]
					break
				}
				// advance down the parse tree
				cmdDef = next
				cmdDefs = next.CommandDefs
			}
		}
	}

	if cmdDef != nil {
		cmd = makeCommand(cmdDef, args)
	}
	return
}

func (c *Command) parseRemainder(tokens []string) (clause string, extArgs []string, err error) {
	clause = ""
	extArgs = []string{}
	err = nil
	if len(tokens) == 0 {
		return
	}

	clauseSl := []string{}
	cmdDef := c.CommandDef
	for _, tokenOrig := range tokens {
		token := strings.ToLower(tokenOrig)
		next, ok := cmdDef.CommandDefs[token]
		if ok {
			clauseSl = append(clauseSl, token)
		} else {
			if len(cmdDef.ArgDefs) > 0 {
				// assume the rest is args
				clause = strings.Join(clauseSl, " ")
				extArgs = strings.Split(tokenOrig, ",")
				return
			}
			var skipped []string
			next, skipped = cmdDef.skipOptionalTo(token)
			if next == nil {
				err = fmt.Errorf("syntax error at '%s'", token)
				return
			}
			clauseSl = append(clauseSl, skipped...)
		}
		cmdDef = next
	}
	clause = strings.Join(clauseSl, " ")
	return
}

func (c *CommandDef) skipOptionalTo(token string) (*CommandDef, []string) {
	if len(c.CommandDefs) == 0 {
		return nil, []string{}
	}
	return c._skipOptionalToWork(token, []string{})
}

func (c *CommandDef) _skipOptionalToWork(token string, skipped []string) (*CommandDef, []string) {
	for _, next := range c.CommandDefs {
		if next.Name == strings.ToLower(token) {
			return next, skipped
		}
		if next.IsOptional() {
			skipped = append(skipped, next.Name)
			return next._skipOptionalToWork(token, skipped)
		}
	}
	return nil, skipped
}

// ShortDesc returns the `ShortDesc` field of a command. ShortDesc panics
// if the command is zero length or contains an unrecognized command.
func ShortDesc(command string) string {
	cmd, _, err := TryParse(command)
	if cmd == nil {
		panic(err.Error())
	}

	return cmd.CommandDef.ShortDesc
}

func makeCommand(cmdDef *CommandDef, args []string) *Command {
	lastArg := ""
	if len(args) > 0 {
		lastArg = args[len(args)-1]
	}

	var ref *string
	if len(lastArg) > 0 && strings.HasPrefix(lastArg, "@") {
		r := lastArg[1:]
		ref = &r
		args = args[:len(args)-1]
	}
	return &Command{
		CommandDef: cmdDef,
		RootDef:    cmdDef.getRoot(),
		Args:       args,
		Ref:        ref,
	}
}
