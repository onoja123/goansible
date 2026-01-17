package modules

import (
	"fmt"
	"goansible/pkg/playbook"
)

type FileModule struct {
	args *playbook.FileArgs
}

func NewFileModule(args *playbook.FileArgs) *FileModule {
	return &FileModule{args: args}
}

func (m *FileModule) Name() string {
	return "file"
}

func (m *FileModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: true, Output: "check mode"}, nil
	}

	var cmd string

	switch m.args.State {
	case "directory":
		cmd = fmt.Sprintf("mkdir -p %s", m.args.Path)
	case "absent":
		cmd = fmt.Sprintf("rm -rf %s", m.args.Path)
	case "touch":
		cmd = fmt.Sprintf("touch %s", m.args.Path)
	case "file":
		cmd = fmt.Sprintf("test -f %s || touch %s", m.args.Path, m.args.Path)
	default:
		return nil, fmt.Errorf("unsupported state: %s", m.args.State)
	}

	if m.args.Mode != "" {
		cmd += fmt.Sprintf(" && chmod %s %s", m.args.Mode, m.args.Path)
	}

	var output string
	var err error

	if ctx.Become {
		output, err = ctx.Connection.ExecuteWithSudo(cmd, ctx.BecomeUser)
	} else {
		output, err = ctx.Connection.Execute(cmd)
	}

	return &Result{
		Changed: true,
		Output:  output,
		Failed:  err != nil,
	}, err
}
