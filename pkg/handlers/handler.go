package handlers

import (
	"goansible/pkg/playbook"
)

type HandlerRegistry struct {
	handlers map[string]*playbook.Task
	notified map[string]bool
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]*playbook.Task),
		notified: make(map[string]bool),
	}
}

func (r *HandlerRegistry) Register(name string, handler *playbook.Task) {
	r.handlers[name] = handler
}

func (r *HandlerRegistry) Notify(name string) {
	r.notified[name] = true
}

func (r *HandlerRegistry) GetNotified() []*playbook.Task {
	var tasks []*playbook.Task
	for name := range r.notified {
		if handler, ok := r.handlers[name]; ok {
			tasks = append(tasks, handler)
		}
	}
	r.notified = make(map[string]bool)
	return tasks
}
