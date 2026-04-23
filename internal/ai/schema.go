package ai

func investigationJSONSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"properties": map[string]any{
			"summary": map[string]any{
				"type": "string",
			},
			"timeline": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
			"suspicious_signals": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
			"recommendations": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required": []string{
			"summary",
			"timeline",
			"suspicious_signals",
			"recommendations",
		},
	}
}
