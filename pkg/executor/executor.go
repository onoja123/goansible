package executor

import (
	"fmt"
	"sync"

	"goansible/pkg/config"
	"goansible/pkg/inventory"
	"goansible/pkg/logger"
	"goansible/pkg/modules"
	"goansible/pkg/playbook"
)

type Executor struct {
	config    *config.Config
	logger    logger.Logger
	checkMode bool
	conns     map[string]*SSHConnection
	mu        sync.Mutex
}

func NewExecutor(cfg *config.Config, log logger.Logger) *Executor {
	return &Executor{
		config: cfg,
		logger: log,
		conns:  make(map[string]*SSHConnection),
	}
}

func (e *Executor) SetCheckMode(check bool) {
	e.checkMode = check
}

func (e *Executor) ExecutePlaybook(pb *playbook.Playbook, inv *inventory.Inventory) error {
	for i, play := range pb.Plays {
		e.logger.Info("\nPLAY [%s] ******************", play.Name)

		hosts := inv.GetHosts(play.Hosts)
		if len(hosts) == 0 {
			e.logger.Warn("No hosts matched for play: %s", play.Name)
			continue
		}

		if err := e.executePlay(play, hosts); err != nil {
			return fmt.Errorf("play %d failed: %w", i, err)
		}
	}

	return nil
}

func (e *Executor) executePlay(play *playbook.Play, hosts []*inventory.Host) error {
	// Execute pre_tasks
	if len(play.PreTasks) > 0 {
		e.logger.Info("Running pre_tasks...")
		if err := e.executeTasks(play.PreTasks, hosts, play); err != nil {
			return err
		}
	}

	// Execute main tasks
	if len(play.Tasks) > 0 {
		if err := e.executeTasks(play.Tasks, hosts, play); err != nil {
			return err
		}
	}

	// Execute post_tasks
	if len(play.PostTasks) > 0 {
		e.logger.Info("Running post_tasks...")
		if err := e.executeTasks(play.PostTasks, hosts, play); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) executeTasks(tasks []*playbook.Task, hosts []*inventory.Host, play *playbook.Play) error {
	for _, task := range tasks {
		e.logger.Info("\nTASK [%s] ******************", task.Name)

		var wg sync.WaitGroup
		results := make(chan *TaskResult, len(hosts))

		for _, host := range hosts {
			wg.Add(1)
			go func(h *inventory.Host) {
				defer wg.Done()
				result := e.executeTask(task, h, play)
				results <- result
			}(host)
		}

		wg.Wait()
		close(results)

		for result := range results {
			if result.Error != nil {
				e.logger.Error("failed: [%s] => %v", result.Host.Name, result.Error)
				if !task.Ignore {
					return result.Error
				}
			} else if result.Changed {
				e.logger.Changed("changed: [%s]", result.Host.Name)
			} else {
				e.logger.Success("ok: [%s]", result.Host.Name)
			}

			if result.Output != "" {
				e.logger.Info("  Output: %s", result.Output)
			}
		}
	}

	return nil
}

func (e *Executor) executeTask(task *playbook.Task, host *inventory.Host, play *playbook.Play) *TaskResult {
	conn, err := e.getConnection(host)
	if err != nil {
		return &TaskResult{Host: host, Error: err}
	}

	ctx := &modules.ExecutionContext{
		Host:       host,
		Connection: conn,
		CheckMode:  e.checkMode,
		Become:     play.Become,
		BecomeUser: play.BecomeUser,
		Vars:       play.Vars,
	}

	mod := e.resolveModule(task)
	if mod == nil {
		return &TaskResult{Host: host, Error: fmt.Errorf("unknown module")}
	}

	result, err := mod.Execute(ctx)
	return &TaskResult{
		Host:    host,
		Changed: result.Changed,
		Output:  result.Output,
		Error:   err,
	}
}

func (e *Executor) resolveModule(task *playbook.Task) modules.Module {
	if task.Command != "" {
		return modules.NewCommandModule(task.Command)
	}
	if task.Shell != "" {
		return modules.NewShellModule(task.Shell)
	}
	if task.Copy != nil {
		return modules.NewCopyModule(task.Copy)
	}
	if task.File != nil {
		return modules.NewFileModule(task.File)
	}
	if task.Template != nil {
		return modules.NewTemplateModule(task.Template)
	}
	if task.Service != nil {
		return modules.NewServiceModule(task.Service)
	}
	if task.Apt != nil {
		return modules.NewAptModule(task.Apt)
	}
	if task.Debug != nil {
		return modules.NewDebugModule(task.Debug)
	}
	return nil
}

func (e *Executor) getConnection(host *inventory.Host) (*SSHConnection, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if conn, ok := e.conns[host.Name]; ok {
		return conn, nil
	}

	conn, err := NewSSHConnection(host)
	if err != nil {
		return nil, err
	}

	e.conns[host.Name] = conn
	return conn, nil
}

type TaskResult struct {
	Host    *inventory.Host
	Changed bool
	Output  string
	Error   error
}
