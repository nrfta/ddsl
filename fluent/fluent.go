package fluent

import (
	"github.com/neighborly/ddsl/exec"
	"github.com/neighborly/ddsl/parser"
)

type FluentDDSL struct {
	ctx *exec.Context
	cmd string
}

func New(repo, dbURL string) *FluentDDSL {
	f := &FluentDDSL{}
	f.ctx = exec.NewContext(repo, dbURL, false, false)
	return f
}

func (f *FluentDDSL) Create() *Create {
	f.cmd = exec.CREATE
	return &Create{f}
}

func (f *FluentDDSL) Drop() *Drop {
	f.cmd = exec.DROP
	return &Drop{f}
}

func (f *FluentDDSL) Grant() *Grant {
	f.cmd = exec.GRANT
	return &Grant{f}
}

func (f *FluentDDSL) Revoke() *Revoke {
	f.cmd = exec.REVOKE
	return &Revoke{f}
}

func (f *FluentDDSL) Migrate() *Migrate {
	f.cmd = exec.MIGRATE
	return &Migrate{f}
}

func (f *FluentDDSL) execute() error {
	cmds, _, _, err := parser.Parse(f.cmd)
	if err != nil {
		return err
	}
	err = exec.ExecuteBatch(f.ctx, cmds)
	if err != nil {
		return err
	}
	return nil
}

func (f *FluentDDSL) addTokens(tokens ...string) {
	for _, t := range tokens {
		f.cmd += " " + t
	}
}