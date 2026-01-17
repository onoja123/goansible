package modules

import (
	"goansible/pkg/executor"
	"goansible/pkg/inventory"
)

type Module interface {
	Execute(ctx *ExecutionContext) (*Result, error)
	Name() string
}

type ExecutionContext struct {
	Host       *inventory.Host
	Connection *executor.SSHConnection
	CheckMode  bool
	Become     bool
	BecomeUser string
	Vars       map[string]interface{}
}

type Result struct {
	Changed bool
	Output  string
	Failed  bool
	Data    map[string]interface{}
}
