package ai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func normalizeBaseURL(raw string, fallback string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		value = fallback
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return strings.TrimRight(value, "/")
	}
	return strings.TrimRight(parsed.String(), "/")
}

func marshalBody(value any) ([]byte, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ai request: %w", err)
	}
	return body, nil
}
