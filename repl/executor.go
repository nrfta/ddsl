package repl

import (
	"github.com/neighborly/ddsl/exec"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
)

func executor(command string) {
	cmds, _,_,err := parser.Parse(command)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = exec.ExecuteBatch(cache.context, cmds)
	if err != nil {
		log.Error(err.Error())
	}
}
