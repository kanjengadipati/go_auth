package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenAIProviderGenerateSuccess(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.test", "test-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/v1/responses", req.URL.Path)
			assert.Equal(t, "Bearer test-key", req.Header.Get("Authorization"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"output":[
						{
							"type":"message",
							"content":[{"type":"output_text","text":"{\"summary\":\"ok\",\"timeline\":[],\"suspicious_signals\":[],\"recommendations\":[]}"}]
						}
					]
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:        "gpt-4.1-mini",
		SystemPrompt: "system",
		UserPrompt:   "user",
		MaxTokens:    300,
	})

	assert.NoError(t, err)
	assert.Contains(t, result.Text, `"summary":"ok"`)
}

func TestOpenAIProviderGenerateReturnsAPIError(t *testing.T) {
	provider := NewOpenAIProvider("https://api.openai.test", "test-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"bad api key"}}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:      "gpt-4.1-mini",
		UserPrompt: "hello",
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "openai error: bad api key")
}
