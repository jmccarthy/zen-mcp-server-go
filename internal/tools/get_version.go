package tools

import (
	"context"
	"runtime"

	"github.com/google/uuid"
)

// GetVersionTool returns version and basic configuration info.
type GetVersionTool struct{}

func (t *GetVersionTool) Name() string { return "get_version" }

func (t *GetVersionTool) Description() string {
	return "VERSION & CONFIGURATION - Get server version and runtime info"
}

func (t *GetVersionTool) Execute(ctx context.Context, params map[string]any) (any, error) {
	info := map[string]any{
		"version":    "0.1.0",
		"go_version": runtime.Version(),
		"server_id":  uuid.NewString(),
	}
	return info, nil
}
