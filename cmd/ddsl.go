package cmd

import (
	"github.com/neighborly/ddsl/exec"
	"io/ioutil"
)

func runCommand(repo string, url string, command string) (exitCode int, err error) {
	if err := exec.Execute(repo, url, command); err != nil {
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


