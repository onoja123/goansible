package modules

import (
	"fmt"
	"goansible/pkg/playbook"
)

type DebugModule struct {
	args *playbook.DebugArgs
}

func NewDebugModule(args *playbook.DebugArgs) *DebugModule {
	return &DebugModule{args: args}
}

func (m *DebugModule) Name() string {
	return "debug"
}

func (m *DebugModule) Execute(ctx *ExecutionContext) (*Result, error) {
	output := ""

	if m.args.Msg != "" {
		output = m.args.Msg
	} else if m.args.Var != "" {
		if val, ok := ctx.Vars[m.args.Var]; ok {
			output = fmt.Sprintf("%s: %v", m.args.Var, val)
		}
	}

	return &Result{
		Changed: false,
		Output:  output,
		Failed:  false,
	}, nil
}
