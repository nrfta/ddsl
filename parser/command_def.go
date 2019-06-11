package parser

type CommandDef struct {
	Name        string
	ShortDesc   string
	Props       map[string]string
	Level       int
	Parent      *CommandDef
	CommandDefs map[string]*CommandDef
	ArgDefs     map[string]*ArgDef
}

type ArgDef struct {
	Name      string
	ShortDesc string
}

func (c *CommandDef) IsOptional() bool {
	return c.hasProp("optional")
}
func (c *CommandDef) IsRoot() bool {
	return c.hasProp("root")
}
func (c *CommandDef) IsNonExec() bool {
	return c.hasProp("non-exec")
}
func (c *CommandDef) IsPrimary() bool {
	return c.hasProp("primary")
}
func (c *CommandDef) HasExtArgs() bool {
	return c.hasProp("ext-args")
}

func (c *CommandDef) hasProp(name string) bool {
	_, ok := c.Props[name]
	return ok
}

func (c *CommandDef) getRoot() *CommandDef {
	if c.IsRoot() {
		return c
	}
	return c.Parent.getRoot()
}
