package server

import (
	"context"
	"fmt"
)

// Tool is the interface all tools must implement.
type Tool interface {
	Name() string
	Description() string
	Execute(context.Context, map[string]any) (any, error)
}

// Dispatcher routes requests to registered tools.
type Dispatcher struct {
	tools map[string]Tool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{tools: map[string]Tool{}}
}

// Register adds a tool to the dispatcher.
func (d *Dispatcher) Register(t Tool) {
	d.tools[t.Name()] = t
}

// Call executes the tool with the given name.
func (d *Dispatcher) Call(ctx context.Context, name string, params map[string]any) (any, error) {
	tool, ok := d.tools[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return tool.Execute(ctx, params)
}
