package repl

import (
	"github.com/nrfta/ddsl/exec"
	"github.com/nrfta/ddsl/log"
	"github.com/nrfta/ddsl/parser"
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
