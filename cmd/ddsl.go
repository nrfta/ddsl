package cmd

import (
	db "github.com/golang-migrate/migrate/database"
	src "github.com/golang-migrate/migrate/source"
	"github.com/neighborly/ddsl/exec"
	"github.com/neighborly/ddsl/parser"
	"io/ioutil"
)

func runCommand(repo string, url string, command string) (exitCode int, err error) {
	sourceDriver, dbDriver, trees, err := getComponents(repo, url, command)
	if err != nil {
		return 1, err
	}

	for _, t := range trees {
		err := exec.ExecuteTree(sourceDriver, dbDriver, t)
		if err != nil {
			//dbDriver.Rollback()
			return 1, err
		}
	}

	//dbDriver.Commit()
	return 0, nil
}

func runFile(repo string, url string, file string) (exitCode int, err error) {
	commandBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return 1, err
	}

	command = string(commandBytes)

	sourceDriver, dbDriver, trees, err := getComponents(repo, url, command)
	if err != nil {
		return 1, err
	}

	for _, t := range trees {
		err := exec.ExecuteTree(sourceDriver, dbDriver, t)
		if err != nil {
			//dbDriver.Rollback()
			return 1, err
		}
	}

	//dbDriver.Commit()
	return 0, nil
}

func getComponents(repo string, url string, command string) (src.Driver, db.Driver, []*parser.DDSL, error) {
	dbDriver, err := db.Open(url)
	if err != nil {
		return nil, nil, nil, err
	}

	sourceDriver, err := src.Open(repo)
	if err != nil {
		return nil, nil, nil, err
	}

	trees, err := parser.Parse(command)
	if err != nil {
		return nil, nil, nil, err
	}

	return sourceDriver, dbDriver, trees, nil
}
