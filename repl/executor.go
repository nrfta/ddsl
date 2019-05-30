package repl

import (
	"fmt"
	"github.com/neighborly/ddsl/exec"
)

func executor(command string) {
	err := exec.Execute(cache.repo, cache.url, command)
	if err != nil {
		fmt.Println(err)
	}
}
