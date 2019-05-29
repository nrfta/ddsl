package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/exec"
	"io/ioutil"
	"strings"
)

func runCommand(repo string, url string, command string) (exitCode int, err error) {
	cmds := strings.Split(command, ";")
	for _, cmd := range cmds {
		fmt.Printf("[INFO] *** command: %s ***\n", cmd)
		if err := exec.Execute(repo, url, cmd); err != nil {
			return 1, err
		}
	}
	return 0, nil
}

func runFile(repo string, url string, file string) (exitCode int, err error) {
	commandBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return 1, err
	}

	command = string(commandBytes)

	return runCommand(repo, url, command)
}
