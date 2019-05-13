package cmd

import (
	"github.com/neighborly/ddsl/exec"
	"io/ioutil"
	"strings"
)

func runCommand(repo string, url string, command string) (exitCode int, err error) {
	cmd := strings.Replace(command, ";", "\n", -1)
	if err := exec.Execute(repo, url, cmd); err != nil {
		return 1, err
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


