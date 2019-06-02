package cmd

import (
	"github.com/neighborly/ddsl/exec"
	"github.com/neighborly/ddsl/parser"
	"io/ioutil"
)

func runCLICommand(command string) (exitCode int, err error) {
	cmds, _, hasDB, err := parser.Parse(command)
	if err != nil {
		return 1, err
	}
	ctx := makeExecContext(!hasDB)
	err = exec.ExecuteBatch(ctx, cmds)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func runFile(file string) (exitCode int, err error) {
	commandBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return 1, err
	}

	command := string(commandBytes)
	cmds, hasTx, hasDB, err := parser.Parse(command)

	ctx := makeExecContext(!hasTx && !hasDB)
    err = exec.ExecuteBatch(ctx, cmds)
    if err != nil {
    	return 1, err
	}
	return 0, nil
}
