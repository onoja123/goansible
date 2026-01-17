package modules

type CommandModule struct {
	command string
}

func NewCommandModule(cmd string) *CommandModule {
	return &CommandModule{command: cmd}
}

func (m *CommandModule) Name() string {
	return "command"
}

func (m *CommandModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: false, Output: "check mode"}, nil
	}

	var output string
	var err error

	if ctx.Become {
		output, err = ctx.Connection.ExecuteWithSudo(m.command, ctx.BecomeUser)
	} else {
		output, err = ctx.Connection.Execute(m.command)
	}

	return &Result{
		Changed: true,
		Output:  output,
		Failed:  err != nil,
	}, err
}
