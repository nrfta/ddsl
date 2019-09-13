package exec

import (
	"fmt"
	dbdr "github.com/neighborly/ddsl/drivers/database"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/util"
	"os"
	"strings"
)

type processor struct {
	ctx *Context
}

func (p *processor) process() error {
	dbDriver, err := dbdr.Open(p.ctx.DatbaseUrl)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	p.ctx.dbDriver = dbDriver

	err = ensureAuditTable(p.ctx)
	if err != nil {
		return err
	}

	if p.ctx.AutoTransaction {
		if err = p.beginTransaction(); err != nil {
			return err
		}
	}

	if err = p.processInstructions(); err != nil {
		if p.ctx.AutoTransaction && p.ctx.inTransaction {
			p.rollbackTransaction()
		}
		return err
	}

	if p.ctx.AutoTransaction && p.ctx.inTransaction {
		return p.commitTransaction()
	}

	return nil
}

func (p *processor) processInstructions() error {
	for _, instr := range p.ctx.instructions {
		var err error
		switch instr.instrType {
		case INSTR_BEGIN:
			err = p.beginTransaction()
		case INSTR_COMMIT:
			err = p.commitTransaction()
		case INSTR_ROLLBACK:
			err = p.rollbackTransaction()
		case INSTR_SQL_FILE:
			err = p.executeSQLFile(instr)
		case INSTR_SH_FILE:
			err = p.executeShellScriptFile(instr)
		case INSTR_CSV_FILE:
			err = p.importCSV(instr)
		case INSTR_SQL_SCRIPT:
			err = p.executeSQLScript(instr)
		case INSTR_SH_SCRIPT:
			err = p.executeShellScript(instr)
		case INSTR_DDSL:
			err = p.executeDDSL(instr)
		case INSTR_DDSL_FILE:
			log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing DDSL file %s", instr.params[FILE_PATH].(string))
		case INSTR_DDSL_FILE_END:
			log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "completed executing DDSL file")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *processor) beginTransaction() error {
	if p.ctx.inTransaction {
		return fmt.Errorf("already in transaction")
	}

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "beginning transaction")
	if !p.ctx.DryRun {
		if err := p.ctx.dbDriver.Begin(); err != nil {
			return err
		}
	}
	p.ctx.inTransaction = true

	return nil
}

func (p *processor) commitTransaction() error {
	if !p.ctx.inTransaction {
		return fmt.Errorf("not in transaction")
	}

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "committing transaction")
	if !p.ctx.DryRun {
		if err := p.ctx.dbDriver.Commit(); err != nil {
			return err
		}
	}
	p.ctx.inTransaction = false

	return nil
}

func (p *processor) rollbackTransaction() error {
	if !p.ctx.inTransaction {
		return fmt.Errorf("not in transaction")
	}

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_WARN), "rolling back transaction")
	if !p.ctx.DryRun {
		p.ctx.dbDriver.Rollback()
	}
	p.ctx.inTransaction = false

	return nil
}

func (p *processor) executeSQLFile(instr *instruction) error {
	filePath := instr.params[FILE_PATH].(string)
	fr, err := os.Open(filePath)
	if err != nil {
		return err
	}

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing SQL file %s", filePath)
	if !p.ctx.DryRun {
		return p.ctx.dbDriver.Exec(fr)
	}
	return nil
}

func (p *processor) executeSQLScript(instr *instruction) error {
	sql := instr.params[SQL].(string)

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing SQL script")
	log.Log(levelOrDryRun(p.ctx, log.LEVEL_DEBUG), sql)
	if !p.ctx.DryRun {
		return p.ctx.dbDriver.Exec(strings.NewReader(sql))
	}
	return nil
}

func (p *processor) executeShellScriptFile(instr *instruction) error {
	filePath := instr.params[FILE_PATH].(string)

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing shell script file %s", filePath)
	if !p.ctx.DryRun {
		out, err := util.OSExec("sh", filePath)
		if err != nil {
			return err
		}

		if len(out) > 0 {
			log.Info(out)
		}
	}
	return nil
}

func (p *processor) executeShellScript(instr *instruction) error {
	command := instr.params[COMMAND].(string)
	args := instr.params[ARGS].([]string)

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "executing shell script")
	log.Log(levelOrDryRun(p.ctx, log.LEVEL_DEBUG), command)
	log.Log(levelOrDryRun(p.ctx, log.LEVEL_DEBUG), "[%s]", strings.Join(args, ", "))
	if !p.ctx.DryRun {
		out, err := util.OSExec(command, args...)
		if err != nil {
			return err
		}

		if len(out) > 0 {
			log.Info(out)
		}
	}
	return nil
}

func (p *processor) importCSV(instr *instruction) error {
	filePath := instr.params[FILE_PATH].(string)
	schemaName := instr.params[SCHEMA_NAME].(string)
	tableName := instr.params[TABLE_NAME].(string)

	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "importing CSV %s", filePath)
	if !p.ctx.DryRun {
		// TODO: provide options for delimiter and header
		output, err := p.ctx.dbDriver.ImportCSV(filePath, schemaName, tableName, ",", true)
		if err != nil {
			return err
		}

		if len(output) > 0 {
			log.Info(output)
		}
	}

	return nil
}

func (p *processor) executeDDSL(instr *instruction) error {
	ddslCommand := instr.params[COMMAND].(string)
	log.Log(levelOrDryRun(p.ctx, log.LEVEL_INFO), "DDSL command: %s", ddslCommand)
	if !p.ctx.DryRun {
		err := p.audit(ddslCommand)
		if err != nil {
			return err
		}
	}
	return nil
}

func levelOrDryRun(ctx *Context, level log.LogLevel) log.LogLevel {
	if ctx.DryRun {
		return log.LEVEL_DRY_RUN
	}
	return level
}
