package providers

import "context"

// Provider defines a minimal interface for language model providers.
// It generates text based on a prompt and optional parameters.
type Provider interface {
	Name() string
	Generate(ctx context.Context, prompt string, opts *Options) (string, error)
}

// Options holds generation options common across providers.
type Options struct {
	Model       string
	Temperature float64
	MaxTokens   int
}
