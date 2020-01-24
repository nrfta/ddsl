package exec

import (
	"fmt"
	"os/user"
	"strings"
)

func ensureAuditTable(ctx *Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS ddsl_audit (
		ddsl_command CHARACTER VARYING,
		performed_at TIMESTAMP WITHOUT TIME ZONE,
		by_db_user CHARACTER VARYING,
		by_os_user CHARACTER VARYING
	)`
	return ctx.dbDriver.Exec(strings.NewReader(sql))
}

func (p *processor) audit(ddslCommand string) error {
	sql := `
	INSERT INTO ddsl_audit (ddsl_command, performed_at, by_db_user, by_os_user)
	VALUES ('%s', NOW(), '%s', '%s');
	`
	osUser, err := user.Current()
	if err != nil {
		return err
	}

	sql = fmt.Sprintf(sql, ddslCommand, p.ctx.dbDriver.User(), osUser.Username)
	return p.ctx.dbDriver.Exec(strings.NewReader(sql))

}
