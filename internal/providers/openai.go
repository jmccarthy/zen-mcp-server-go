package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider implements Provider using the OpenAI HTTP API.
type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1/chat/completions",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string { return "openai" }

// generateRequest is the payload sent to OpenAI.
type generateRequest struct {
	Model       string              `json:"model"`
	Messages    []map[string]string `json:"messages"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

type generateResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error any `json:"error"`
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{}
	}
	reqBody := generateRequest{
		Model:       opts.Model,
		Messages:    []map[string]string{{"role": "user", "content": prompt}},
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai: status %d: %s", resp.StatusCode, string(buf))
	}

	var gr generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return "", err
	}
	if len(gr.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices")
	}
	return gr.Choices[0].Message.Content, nil
}
