package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockProviderGenerateBuildsContextualOutput(t *testing.T) {
	provider := NewMockProvider()

	result, err := provider.Generate(context.Background(), GenerateInput{
		UserPrompt: "Audit log count: 3\nOverview: failed_logs=2 unique_ip_addresses=1 unique_actor_user_ids=0 first_seen=2026-04-22T02:31:14Z last_seen=2026-04-22T02:35:14Z\n- log 1\n  time: 2026-04-22T02:31:14Z\n  action: login\n  resource: auth\n  status: failed\n  actor_user_id: n/a\n  ip_address: 203.0.113.10\n  user_agent: Postman\n  description: \"invalid credentials\"\n",
	})

	assert.NoError(t, err)
	assert.Contains(t, result.Text, "203.0.113.10")
	assert.Contains(t, result.Text, "invalid credential")
	assert.Contains(t, result.Text, "\"timeline\"")
}
