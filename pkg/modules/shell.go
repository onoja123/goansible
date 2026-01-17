package modules

type ShellModule struct {
	command string
}

func NewShellModule(cmd string) *ShellModule {
	return &ShellModule{command: cmd}
}

func (m *ShellModule) Name() string {
	return "shell"
}

func (m *ShellModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: false, Output: "check mode"}, nil
	}

	var output string
	var err error

	shellCmd := "/bin/sh -c '" + m.command + "'"

	if ctx.Become {
		output, err = ctx.Connection.ExecuteWithSudo(shellCmd, ctx.BecomeUser)
	} else {
		output, err = ctx.Connection.Execute(shellCmd)
	}

	return &Result{
		Changed: true,
		Output:  output,
		Failed:  err != nil,
	}, err
}
