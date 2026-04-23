package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const defaultOpenAIBaseURL = "https://api.openai.com"

type OpenAIProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type openAIResponseRequest struct {
	Model       string               `json:"model"`
	Input       []openAIInputMessage `json:"input"`
	Text        openAITextConfig     `json:"text"`
	Temperature *float64             `json:"temperature,omitempty"`
	MaxTokens   int                  `json:"max_output_tokens,omitempty"`
}

type openAIInputMessage struct {
	Role    string                    `json:"role"`
	Content []openAIInputContentBlock `json:"content"`
}

type openAIInputContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAITextConfig struct {
	Format openAITextFormat `json:"format"`
}

type openAITextFormat struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Schema      map[string]any `json:"schema"`
	Strict      bool           `json:"strict"`
	Description string         `json:"description,omitempty"`
}

type openAIResponseEnvelope struct {
	Error      *openAIErrorItem   `json:"error"`
	Output     []openAIOutputItem `json:"output"`
	OutputText string             `json:"output_text"`
}

type openAIErrorItem struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type openAIOutputItem struct {
	Type    string                     `json:"type"`
	Content []openAIOutputContentBlock `json:"content"`
}

type openAIOutputContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewOpenAIProvider(baseURL, apiKey string, timeout time.Duration) *OpenAIProvider {
	return &OpenAIProvider{
		baseURL: normalizeBaseURL(baseURL, defaultOpenAIBaseURL),
		apiKey:  apiKey,
		client:  &http.Client{Timeout: timeout},
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	reqBody := openAIResponseRequest{
		Model: input.Model,
		Input: []openAIInputMessage{
			{
				Role: "system",
				Content: []openAIInputContentBlock{
					{Type: "input_text", Text: input.SystemPrompt},
				},
			},
			{
				Role: "user",
				Content: []openAIInputContentBlock{
					{Type: "input_text", Text: input.UserPrompt},
				},
			},
		},
		Text: openAITextConfig{
			Format: openAITextFormat{
				Type:        "json_schema",
				Name:        "audit_investigation",
				Strict:      true,
				Schema:      investigationJSONSchema(),
				Description: "Structured investigation result for audit log analysis.",
			},
		},
		MaxTokens: input.MaxTokens,
	}
	if input.Temperature > 0 {
		reqBody.Temperature = &input.Temperature
	}

	body, err := marshalBody(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/responses", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		return nil, fmt.Errorf("openai is unavailable: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed openAIResponseEnvelope
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to decode openai response: %s", strings.TrimSpace(string(bodyBytes)))
	}

	if resp.StatusCode >= 400 {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return nil, fmt.Errorf("openai error: %s", parsed.Error.Message)
		}
		return nil, fmt.Errorf("openai returned status %d", resp.StatusCode)
	}

	text := firstNonEmpty(parsed.OutputText, extractOpenAIOutputText(parsed.Output))
	if text == "" {
		return nil, ErrInvalidStructuredOutput
	}

	return &GenerateResult{Text: text}, nil
}

func extractOpenAIOutputText(items []openAIOutputItem) string {
	for _, item := range items {
		for _, block := range item.Content {
			if block.Type == "output_text" && block.Text != "" {
				return block.Text
			}
		}
	}
	return ""
}
