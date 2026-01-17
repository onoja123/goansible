package modules

import (
	"fmt"
	"goansible/pkg/playbook"
)

type ServiceModule struct {
	args *playbook.ServiceArgs
}

func NewServiceModule(args *playbook.ServiceArgs) *ServiceModule {
	return &ServiceModule{args: args}
}

func (m *ServiceModule) Name() string {
	return "service"
}

func (m *ServiceModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: true, Output: "check mode"}, nil
	}

	var cmd string
	
	switch m.args.State {
	case "started":
		cmd = fmt.Sprintf("systemctl start %s", m.args.Name)
	case "stopped":
		cmd = fmt.Sprintf("systemctl stop %s", m.args.Name)
	case "restarted":
		cmd = fmt.Sprintf("systemctl restart %s", m.args.Name)
	case "reloaded":
		cmd = fmt.Sprintf("systemctl reload %s", m.args.Name)
	}

	if m.args.Enabled != nil && *m.args.Enabled {
		cmd += fmt.Sprintf(" && systemctl enable %s", m.args.Name)
	}

	output, err := ctx.Connection.ExecuteWithSudo(cmd, "root")

	return &Result{
		Changed: true,
		Output:  output,
		Failed:  err != nil,
	}, err
}

// ==================== pkg/modules/apt.go ====================
package modules

import (
	"fmt"
	"goansible/pkg/playbook"
)

type AptModule struct {
	args *playbook.PackageArgs
}

func NewAptModule(args *playbook.PackageArgs) *AptModule {
	return &AptModule{args: args}
}

func (m *AptModule) Name() string {
	return "apt"
}

func (m *AptModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: true, Output: "check mode"}, nil
	}

	var cmd string
	
	if m.args.Update {
		cmd = "apt-get update && "
	}

	switch m.args.State {
	case "present", "installed":
		cmd += fmt.Sprintf("apt-get install -y %s", m.args.Name)
	case "absent", "removed":
		cmd += fmt.Sprintf("apt-get remove -y %s", m.args.Name)
	case "latest":
		cmd += fmt.Sprintf("apt-get install -y --only-upgrade %s", m.args.Name)
	}

	output, err := ctx.Connection.ExecuteWithSudo(cmd, "root")

	return &Result{
		Changed: true,
		Output:  output,
		Failed:  err != nil,
	}, err
}

// ==================== pkg/modules/template.go ====================
package modules

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
	"goansible/pkg/playbook"
)

type TemplateModule struct {
	args *playbook.TemplateArgs
}

func NewTemplateModule(args *playbook.TemplateArgs) *TemplateModule {
	return &TemplateModule{args: args}
}

func (m *TemplateModule) Name() string {
	return "template"
}

func (m *TemplateModule) Execute(ctx *ExecutionContext) (*Result, error) {
	if ctx.CheckMode {
		return &Result{Changed: true, Output: "check mode"}, nil
	}

	tmplContent, err := ioutil.ReadFile(m.args.Src)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("template").Parse(string(tmplContent))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx.Vars); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", m.args.Dest, buf.String())
	
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

// ==================== pkg/logger/logger.go ====================
package logger

type Logger interface {
	Info(format string, args ...interface{})
	Success(format string, args ...interface{})
	Changed(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// ==================== pkg/logger/console.go ====================
package logger

import (
	"fmt"
	"log"
	"os"
)

type ConsoleLogger struct {
	verbose bool
}

func NewConsoleLogger(verbose bool) *ConsoleLogger {
	return &ConsoleLogger{verbose: verbose}
}

func (l *ConsoleLogger) Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (l *ConsoleLogger) Success(format string, args ...interface{}) {
	fmt.Printf("\033[32m"+format+"\033[0m\n", args...)
}

func (l *ConsoleLogger) Changed(format string, args ...interface{}) {
	fmt.Printf("\033[33m"+format+"\033[0m\n", args...)
}

func (l *ConsoleLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("\033[33mWARNING: "+format+"\033[0m\n", args...)
}

func (l *ConsoleLogger) Error(format string, args ...interface{}) {
	fmt.Printf("\033[31mERROR: "+format+"\033[0m\n", args...)
}

func (l *ConsoleLogger) Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
	os.Exit(1)
}
