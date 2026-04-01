package control

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleMarshalRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		target Target
	}{
		{
			name: "aws/lambda/before-publish",
			target: &AwsLambdaBeforePublishTarget{
				Region:       "us-east-1",
				FunctionName: "my-func",
				Authentication: AwsAuthentication{
					Authentication: &AuthenticationModeAssumeRole{
						AssumeRoleArn: "arn:aws:iam::role/test",
					},
				},
				BeforePublishConfig: BeforePublishConfig{
					RetryTimeout:          5000,
					MaxRetries:            3,
					FailedAction:          "REJECT",
					TooManyRequestsAction: "RETRY",
				},
			},
		},
		{
			name: "http/before-publish",
			target: &HttpBeforePublishTarget{
				Url:     "https://example.com/hook",
				Headers: []Header{{Name: "X-Key", Value: "abc"}},
				Format:  Json,
				BeforePublishConfig: BeforePublishConfig{
					RetryTimeout:          1000,
					MaxRetries:            2,
					FailedAction:          "PUBLISH",
					TooManyRequestsAction: "FAIL",
				},
			},
		},
		{
			name: "hive/text-model-only",
			target: &HiveTextModelOnlyTarget{
				ApiKey:     "hive-key",
				ModelUrl:   "https://api.hive.com/model",
				Thresholds: map[string]int{"sexual": 2, "violence": 1},
				BeforePublishConfig: BeforePublishConfig{
					FailedAction: "REJECT",
				},
			},
		},
		{
			name: "hive/dashboard",
			target: &HiveDashboardTarget{
				ApiKey:          "hive-key",
				CheckWatchLists: boolPtr(true),
			},
		},
		{
			name: "bodyguard/text-moderation",
			target: &BodyguardTextModerationTarget{
				ApiKey:          "bg-key",
				ChannelId:       "chan-1",
				ApiUrl:          "https://api.bodyguard.ai",
				DefaultLanguage: "en",
				BeforePublishConfig: BeforePublishConfig{
					FailedAction: "PUBLISH",
				},
			},
		},
		{
			name: "tisane/text-moderation",
			target: &TisaneTextModerationTarget{
				ApiKey:          "tisane-key",
				ModelUrl:        "https://api.tisane.ai",
				Thresholds:      map[string]int{"hate": 2},
				DefaultLanguage: "en",
				BeforePublishConfig: BeforePublishConfig{
					FailedAction: "REJECT",
				},
			},
		},
		{
			name: "azure/text-moderation",
			target: &AzureTextModerationTarget{
				ApiKey:     "azure-key",
				Endpoint:   "https://contentsafety.azure.com",
				Thresholds: map[string]int{"Hate": 4, "Violence": 2},
				BeforePublishConfig: BeforePublishConfig{
					FailedAction: "REJECT",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newRule := NewRule{
				Status:      "enabled",
				RequestMode: Single,
				Source: Source{
					ChannelFilter: "test.*",
					Type:          ChannelMessage,
				},
				Target: tt.target,
			}

			data, err := json.Marshal(&newRule)
			assert.NoError(t, err)

			var raw map[string]interface{}
			err = json.Unmarshal(data, &raw)
			assert.NoError(t, err)
			assert.Equal(t, tt.name, raw["ruleType"])

			apiResponse := buildMockRuleResponse(t, data, tt.name)

			var rule Rule
			err = json.Unmarshal(apiResponse, &rule)
			assert.NoError(t, err)
			assert.Equal(t, tt.name, rule.Target.TargetType())
			assert.Equal(t, "rule-123", rule.ID)
			assert.Equal(t, "app-456", rule.AppID)
		})
	}
}

func TestRuleUnmarshalUnknownType(t *testing.T) {
	data := []byte(`{
		"id": "rule-1",
		"appId": "app-1",
		"ruleType": "totally/unknown",
		"source": {"channelFilter": "", "type": "channel.message"},
		"target": {}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown rule type")
}

func TestRuleWithInvocationMode(t *testing.T) {
	data := []byte(`{
		"id": "rule-1",
		"appId": "app-1",
		"ruleType": "http/before-publish",
		"invocationMode": "BEFORE_PUBLISH",
		"requestMode": "single",
		"source": {"channelFilter": "", "type": "channel.message"},
		"target": {"url": "https://example.com", "beforePublishConfig": {"retryTimeout": 1000, "maxRetries": 2, "failedAction": "REJECT", "tooManyRequestsAction": "RETRY"}}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)
	assert.Equal(t, BeforePublish, rule.InvocationMode)
}

func TestSourceChatRoomFilter(t *testing.T) {
	data := []byte(`{
		"id": "rule-1",
		"appId": "app-1",
		"ruleType": "http",
		"source": {"channelFilter": "test.*", "type": "channel.message", "chatRoomFilter": "room-.*"},
		"target": {"url": "https://example.com"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)
	assert.Equal(t, "room-.*", rule.Source.ChatRoomFilter)
}

func buildMockRuleResponse(t *testing.T, newRuleData []byte, ruleType string) []byte {
	t.Helper()

	var partial map[string]interface{}
	err := json.Unmarshal(newRuleData, &partial)
	assert.NoError(t, err)

	partial["id"] = "rule-123"
	partial["appId"] = "app-456"
	partial["version"] = "1.0"
	partial["created"] = 1234567890
	partial["modified"] = 1234567890
	partial["ruleType"] = ruleType

	data, err := json.Marshal(partial)
	assert.NoError(t, err)
	return data
}

func boolPtr(b bool) *bool {
	return &b
}
