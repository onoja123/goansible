package modules

import (
	"fmt"
	"goansible/pkg/playbook"
	"io/ioutil"
)

type CopyModule struct {
	args *playbook.CopyArgs
}

func NewCopyModule(args *playbook.CopyArgs) *CopyModule {
	return &CopyModule{args: args}
}

func (m *CopyModule) Name() string {
	return "copy"
}

func (m *CopyModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: true, Output: "check mode"}, nil
	}

	content, err := ioutil.ReadFile(m.args.Src)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", m.args.Dest, string(content))

	if m.args.Mode != "" {
		cmd += fmt.Sprintf(" && chmod %s %s", m.args.Mode, m.args.Dest)
	}

	var output string
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
