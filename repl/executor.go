package repl

import (
	"fmt"
	"github.com/neighborly/ddsl/exec"
)

func executor(command string) {
	ctx := exec.NewContext(cache.repo, cache.url, false)
	err := exec.Execute(ctx, command)
	if err != nil {
		fmt.Println(err)
	}
}
