package tools

import (
	"context"
	"fmt"

	"github.com/BeehiveInnovations/zen-mcp-server-go/internal/conversation"
	"github.com/BeehiveInnovations/zen-mcp-server-go/internal/providers"
)

// ChatTool is a minimal chat implementation that forwards prompts to a provider.
type ChatTool struct {
	ProviderName string
	Store        *conversation.Store
}

func (t *ChatTool) Name() string { return "chat" }

func (t *ChatTool) Description() string {
	return "GENERAL CHAT - Send a prompt to the language model"
}

func (t *ChatTool) Execute(ctx context.Context, params map[string]any) (any, error) {
	prompt, _ := params["prompt"].(string)
	if prompt == "" {
		return nil, fmt.Errorf("prompt required")
	}
	convID, _ := params["continuation_id"].(string)
	var thread *conversation.Thread
	var err error
	if convID != "" {
		thread, err = t.Store.Get(ctx, convID)
		if err != nil {
			thread = nil
		}
	}
	if thread == nil {
		thread, err = t.Store.CreateThread(ctx)
		if err != nil {
			return nil, err
		}
	}
	prov := providers.Get(t.ProviderName)
	if prov == nil {
		return nil, fmt.Errorf("provider not configured")
	}
	text, err := prov.Generate(ctx, prompt, &providers.Options{})
	if err != nil {
		return nil, err
	}
	t.Store.AddTurn(ctx, thread.ID, conversation.Turn{Role: "user", Content: prompt, Tool: t.Name()})
	t.Store.AddTurn(ctx, thread.ID, conversation.Turn{Role: "assistant", Content: text, Tool: t.Name()})
	return map[string]any{"continuation_id": thread.ID, "response": text}, nil
}
