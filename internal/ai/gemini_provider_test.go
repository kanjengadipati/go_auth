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

func TestGeminiProviderGenerateSuccess(t *testing.T) {
	provider := NewGeminiProvider("https://gemini.test", "gem-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/v1beta/models/gemini-2.5-flash:generateContent", req.URL.Path)
			assert.Equal(t, "gem-key", req.Header.Get("x-goog-api-key"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"candidates":[
						{"content":{"parts":[{"text":"{\"summary\":\"ok\",\"timeline\":[],\"suspicious_signals\":[],\"recommendations\":[]}"}]}}
					]
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:        "gemini-2.5-flash",
		SystemPrompt: "system",
		UserPrompt:   "user",
		MaxTokens:    300,
	})

	assert.NoError(t, err)
	assert.Contains(t, result.Text, `"summary":"ok"`)
}

func TestGeminiProviderGenerateReturnsAPIError(t *testing.T) {
	provider := NewGeminiProvider("https://gemini.test", "gem-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"unsupported model"}}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:      "gemini-2.5-flash",
		UserPrompt: "hello",
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "gemini error: unsupported model")
}
