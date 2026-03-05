package control

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleTargets(t *testing.T) {
	tests := []struct {
		name   string
		target Target
	}{
		{"Pulsar", &PulsarTarget{
			RoutingKey:    "aaaaa",
			Topic:         "my-tenant/my-namespace/my-topic",
			ServiceURL:    "pulsar://test.com:1234",
			TlsTrustCerts: []string{"-----BEGIN CERTIFICATE-----\naaaaa\n-----END CERTIFICATE-----"},
			Authentication: PulsarAuthentication{
				AuthenticationMode: "token",
				Token:              "1234",
			},
			Enveloped: true,
			Format:    Json,
		}},
		{"Kafka", &KafkaTarget{
			RoutingKey: "1234",
			Brokers:    []string{"a", "b", "c"},
			Authentication: KafkaAuthentication{
				Sasl: Sasl{
					Mechanism: Plain,
					Username:  "b",
					Password:  "c",
				},
			},
			Enveloped: false,
			Format:    Json,
		}},
		{"AmqpExternal", &AmqpExternalTarget{
			Url:                "amqps://test.com",
			RoutingKey:         "key",
			MandatoryRoute:     true,
			PersistentMessages: true,
			MessageTTL:         50,
			Headers:            []Header{{Name: "a", Value: "b"}},
			Enveloped:          true,
			Format:             Json,
		}},
		{"Amqp", &AmqpTarget{
			Headers:   []Header{{Name: "a", Value: "b"}},
			Enveloped: true,
			Format:    Json,
		}},
		{"AwsSqs", &AwsSqsTarget{
			Region:       "us-east-2",
			AwsAccountID: "b",
			QueueName:    "c",
			Authentication: AwsAuthentication{
				Authentication: &AuthenticationModeAssumeRole{
					AssumeRoleArn: "aaaaaaa",
				},
			},
			Enveloped: true,
			Format:    Json,
		}},
		{"AwsKinesis", &AwsKinesisTarget{
			Region:       "us-east-2",
			StreamName:   "aaaaaaa",
			PartitionKey: "bbbbbbb",
			Authentication: AwsAuthentication{
				Authentication: &AuthenticationModeAssumeRole{
					AssumeRoleArn: "aaaaaaa",
				},
			},
			Enveloped: true,
			Format:    Json,
		}},
		{"AwsLambda", &AwsLambdaTarget{
			Region:       "us-east-2",
			FunctionName: "heck",
			Authentication: AwsAuthentication{
				Authentication: &AuthenticationModeAssumeRole{
					AssumeRoleArn: "aaaaaaa",
				},
			},
			Enveloped: true,
		}},
		{"HttpGoogleCloudFunction", &HttpGoogleCloudFunctionTarget{
			Region:       "us",
			ProjectID:    "1234",
			FunctionName: "heck",
			Headers:      []Header{{Name: "a", Value: "b"}},
			SigningKeyID: "1234",
			Enveloped:    true,
			Format:       Json,
		}},
		{"HttpAzureFunction", &HttpAzureFunctionTarget{
			AzureAppID:        "420",
			AzureFunctionName: "heck",
			Headers:           []Header{{Name: "a", Value: "b"}},
			Enveloped:         true,
			Format:            Json,
		}},
		{"HttpCloudfareWorker", &HttpCloudfareWorkerTarget{
			Url:     "https://example.com",
			Headers: []Header{{Name: "a", Value: "b"}},
		}},
		{"HttpZapier", &HttpZapierTarget{
			Url:     "https://example.com",
			Headers: []Header{{Name: "a", Value: "b"}},
		}},
		{"HttpIfttt", &HttpIftttTarget{
			WebhookKey: "aaa",
			EventName:  "bbb",
		}},
		{"Http", &HttpTarget{
			Url:       "https://example.com",
			Headers:   []Header{{Name: "a", Value: "b"}},
			Enveloped: true,
			Format:    MsgPack,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRule(t, tt.target)
		})
	}
}

func testRule(t *testing.T, target Target) {
	t.Helper()
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	queue := NewQueue{
		Name:      "test-queue",
		Ttl:       60,
		MaxLength: 100,
		Region:    "us-east-1-a",
	}

	key := NewKey{
		Name:       "test-key",
		Capability: map[string][]string{"foo": {"publish"}},
	}

	q, err := client.CreateQueue(app.ID, &queue)
	require.NoError(t, err)

	k, err := client.CreateKey(app.ID, &key)
	require.NoError(t, err)

	switch t := target.(type) {
	case *AmqpTarget:
		t.QueueID = q.ID
	case *HttpGoogleCloudFunctionTarget:
		t.SigningKeyID = k.ID
	case *HttpAzureFunctionTarget:
		t.SigningKeyID = k.ID
	case *HttpCloudfareWorkerTarget:
		t.SigningKeyID = k.ID
	case *HttpZapierTarget:
		t.SigningKeyID = k.ID
	case *HttpTarget:
		t.SigningKeyID = k.ID
	}

	rule := NewRule{
		Status:      "enabled",
		RequestMode: Single,
		Source: Source{
			ChannelFilter: "aaa",
			Type:          ChannelMessage,
		},
		Target: target,
	}

	r, err := client.CreateRule(app.ID, &rule)
	require.NoError(t, err)
	assert.Equal(t, rule.RequestMode, r.RequestMode)
	assert.Equal(t, rule.Source, r.Source)
	assert.Equal(t, rule.Target.TargetType(), r.Target.TargetType())

	// TlsTrustCerts is write-only in the API — accepted on create but never
	// returned. Nil it out on the sent target before comparing so the full
	// struct equality check still works.
	if p, ok := rule.Target.(*PulsarTarget); ok {
		assert.Nil(t, r.Target.(*PulsarTarget).TlsTrustCerts, "API should not echo back TlsTrustCerts")
		p.TlsTrustCerts = nil
	}
	// The API populates default Format when not explicitly set.
	if l, ok := rule.Target.(*AwsLambdaTarget); ok && l.Format == "" {
		l.Format = r.Target.(*AwsLambdaTarget).Format
	}
	assert.Equal(t, rule.Target, r.Target)
	assert.NotEmpty(t, r.ID)
	assert.NotEmpty(t, r.AppID)
	assert.NotEmpty(t, r.Version)
	assert.NotEmpty(t, r.Status)
	assert.NotEmpty(t, r.Created)
	assert.NotEmpty(t, r.Modified)

	r2, err := client.Rule(app.ID, r.ID)
	require.NoError(t, err)
	assert.Equal(t, r, r2)

	err = client.DeleteRule(app.ID, r.ID)
	assert.NoError(t, err)
}

func TestSaslMechanismConstants(t *testing.T) {
	assert.Equal(t, SaslMechanism("scram-sha-256"), ScramSha256)
	assert.Equal(t, SaslMechanism("scram-sha-512"), ScramSha512)

	// Deprecated constants preserve old (buggy) values for backwards compat.
	assert.Equal(t, SaslMechanism("scra-sha-256"), Scram_sha_256)
	assert.Equal(t, SaslMechanism("scra-sha-512"), Scram_sha_512)
}

func TestPulsarRuleUnmarshalNoStdout(t *testing.T) {
	payload := `{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 1000,
		"modified": 2000,
		"ruleType": "pulsar",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {
			"routingKey": "key",
			"topic": "persistent://tenant/ns/topic",
			"serviceUrl": "pulsar://localhost:6650",
			"tlsTrustCerts": [],
			"authentication": {"authenticationMode": "token", "token": "abc"},
			"enveloped": false
		}
	}`

	// Capture stdout to verify no debug printing occurs.
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var rule Rule
	err := json.Unmarshal([]byte(payload), &rule)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Empty(t, buf.String(), "UnmarshalJSON should not print to stdout")
	assert.Equal(t, "pulsar", rule.Target.TargetType())

	pt, ok := rule.Target.(*PulsarTarget)
	assert.True(t, ok)
	assert.Equal(t, "persistent://tenant/ns/topic", pt.Topic)
	assert.Equal(t, "pulsar://localhost:6650", pt.ServiceURL)
	assert.Equal(t, PulsarAuthenticationMode("token"), pt.Authentication.AuthenticationMode)
}

func TestAmqpExternalTargetIncludesExchange(t *testing.T) {
	target := AmqpExternalTarget{
		Url:      "amqp://example.com",
		Exchange: "my-exchange",
	}

	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Equal(t, "my-exchange", raw["exchange"])
}

func TestAwsLambdaTargetIncludesFormat(t *testing.T) {
	target := AwsLambdaTarget{
		Region:       "us-east-1",
		FunctionName: "my-func",
		Format:       Json,
	}

	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Equal(t, "json", raw["format"])
}

func TestRuleUnmarshalIncludesLinks(t *testing.T) {
	data := `{
		"id": "rule123",
		"appId": "app456",
		"status": "enabled",
		"created": 1700000000,
		"modified": 1700001000,
		"_links": {
			"self": "https://control.ably.net/v1/apps/app456/rules/rule123"
		},
		"ruleType": "http",
		"requestMode": "single",
		"source": {
			"channelFilter": "",
			"type": "channel.message"
		},
		"target": {
			"url": "https://example.com/webhook",
			"signingKeyId": "",
			"enveloped": false
		}
	}`

	var rule Rule
	err := json.Unmarshal([]byte(data), &rule)
	assert.NoError(t, err)
	assert.NotNil(t, rule.Links)
	assert.Equal(t, "https://control.ably.net/v1/apps/app456/rules/rule123", rule.Links["self"])
}

func TestHttpBeforePublishTargetRoundTrip(t *testing.T) {
	target := &HttpBeforePublishTarget{
		Url:     "https://example.com/webhook",
		Headers: []Header{{Name: "X-Custom", Value: "test"}},
		Format:  Json,
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded HttpBeforePublishTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.Url, decoded.Url)
	assert.Equal(t, target.Headers, decoded.Headers)
	assert.Equal(t, target.Format, decoded.Format)
	assert.Equal(t, "http/before-publish", target.TargetType())
}

func TestAwsLambdaBeforePublishTargetRoundTrip(t *testing.T) {
	target := &AwsLambdaBeforePublishTarget{
		Region:       "us-east-1",
		FunctionName: "myFunction",
		Authentication: AwsAuthentication{
			Authentication: &AuthenticationModeAssumeRole{
				AssumeRoleArn: "arn:aws:iam::123456789:role/my-role",
			},
		},
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded AwsLambdaBeforePublishTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.Region, decoded.Region)
	assert.Equal(t, target.FunctionName, decoded.FunctionName)
	assert.Equal(t, "aws/lambda/before-publish", target.TargetType())
}

func TestHiveTextModelOnlyTargetRoundTrip(t *testing.T) {
	target := &HiveTextModelOnlyTarget{
		ApiKey:     "key123",
		ModelUrl:   "https://example.com/model",
		Thresholds: map[string]int{"violence": 2, "hate": 1},
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded HiveTextModelOnlyTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.ApiKey, decoded.ApiKey)
	assert.Equal(t, target.ModelUrl, decoded.ModelUrl)
	assert.Equal(t, target.Thresholds, decoded.Thresholds)
	assert.Equal(t, "hive/text-model-only", target.TargetType())
}

func TestHiveDashboardTargetRoundTrip(t *testing.T) {
	check := true
	target := &HiveDashboardTarget{
		ApiKey:          "key123",
		CheckWatchLists: &check,
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded HiveDashboardTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.ApiKey, decoded.ApiKey)
	assert.NotNil(t, decoded.CheckWatchLists)
	assert.True(t, *decoded.CheckWatchLists)
	assert.Equal(t, "hive/dashboard", target.TargetType())
}

func TestBodyguardTextModerationTargetRoundTrip(t *testing.T) {
	target := &BodyguardTextModerationTarget{
		ApiKey:          "key123",
		ChannelId:       "chan123",
		ApiUrl:          "https://api.bodyguard.ai",
		DefaultLanguage: "en",
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded BodyguardTextModerationTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.ApiKey, decoded.ApiKey)
	assert.Equal(t, target.ChannelId, decoded.ChannelId)
	assert.Equal(t, target.ApiUrl, decoded.ApiUrl)
	assert.Equal(t, target.DefaultLanguage, decoded.DefaultLanguage)
	assert.Equal(t, "bodyguard/text-moderation", target.TargetType())
}

func TestTisaneTextModerationTargetRoundTrip(t *testing.T) {
	target := &TisaneTextModerationTarget{
		ApiKey:          "key123",
		ModelUrl:        "https://example.com/tisane",
		Thresholds:      map[string]int{"profanity": 1, "bigotry": 2},
		DefaultLanguage: "en",
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded TisaneTextModerationTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.ApiKey, decoded.ApiKey)
	assert.Equal(t, target.ModelUrl, decoded.ModelUrl)
	assert.Equal(t, target.Thresholds, decoded.Thresholds)
	assert.Equal(t, target.DefaultLanguage, decoded.DefaultLanguage)
	assert.Equal(t, "tisane/text-moderation", target.TargetType())
}

func TestAzureTextModerationTargetRoundTrip(t *testing.T) {
	target := &AzureTextModerationTarget{
		ApiKey:     "key123",
		Endpoint:   "https://eastus.api.cognitive.microsoft.com",
		Thresholds: map[string]int{"Hate": 4, "Violence": 2},
	}
	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var decoded AzureTextModerationTarget
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, target.ApiKey, decoded.ApiKey)
	assert.Equal(t, target.Endpoint, decoded.Endpoint)
	assert.Equal(t, target.Thresholds, decoded.Thresholds)
	assert.Equal(t, "azure/text-moderation", target.TargetType())
}

func TestBeforePublishConfigInNewRule(t *testing.T) {
	rule := &NewRule{
		Status:         "enabled",
		InvocationMode: BeforePublish,
		BeforePublishConfig: &BeforePublishConfig{
			RetryTimeout:          5000,
			MaxRetries:            3,
			FailedAction:          FailedActionReject,
			TooManyRequestsAction: TooManyRequestsRetry,
		},
		ChatRoomFilter: "^chat:.*",
		RequestMode:    Single,
		Source:         Source{ChannelFilter: "test", Type: ChannelMessage},
		Target: &HttpBeforePublishTarget{
			Url: "https://example.com/webhook",
		},
	}

	data, err := json.Marshal(rule)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Equal(t, "BEFORE_PUBLISH", raw["invocationMode"])
	assert.Equal(t, "^chat:.*", raw["chatRoomFilter"])
	assert.Equal(t, "http/before-publish", raw["ruleType"])

	bpc, ok := raw["beforePublishConfig"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, float64(5000), bpc["retryTimeout"])
	assert.Equal(t, float64(3), bpc["maxRetries"])
	assert.Equal(t, "REJECT", bpc["failedAction"])
	assert.Equal(t, "RETRY", bpc["tooManyRequestsAction"])
}

func TestRuleUnmarshalHiveTextModelOnly(t *testing.T) {
	data := []byte(`{
		"id": "rule123",
		"appId": "app123",
		"version": "1.0",
		"status": "enabled",
		"created": 1234567890,
		"modified": 1234567890,
		"ruleType": "hive/text-model-only",
		"invocationMode": "BEFORE_PUBLISH",
		"beforePublishConfig": {
			"retryTimeout": 5000,
			"maxRetries": 3,
			"failedAction": "REJECT",
			"tooManyRequestsAction": "RETRY"
		},
		"chatRoomFilter": "^chat:.*",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"apiKey": "key123", "modelUrl": "https://example.com", "thresholds": {"violence": 2}}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	assert.Equal(t, "rule123", rule.ID)
	assert.Equal(t, "app123", rule.AppID)
	assert.Equal(t, "1.0", rule.Version)
	assert.Equal(t, "enabled", rule.Status)
	assert.Equal(t, int64(1234567890), rule.Created)
	assert.Equal(t, int64(1234567890), rule.Modified)
	assert.Equal(t, InvocationMode("BEFORE_PUBLISH"), rule.InvocationMode)
	assert.NotNil(t, rule.BeforePublishConfig)
	assert.Equal(t, 5000, rule.BeforePublishConfig.RetryTimeout)
	assert.Equal(t, 3, rule.BeforePublishConfig.MaxRetries)
	assert.Equal(t, FailedActionReject, rule.BeforePublishConfig.FailedAction)
	assert.Equal(t, TooManyRequestsRetry, rule.BeforePublishConfig.TooManyRequestsAction)
	assert.Equal(t, "^chat:.*", rule.ChatRoomFilter)
	assert.Equal(t, Single, rule.RequestMode)
	assert.Equal(t, "test", rule.Source.ChannelFilter)
	assert.Equal(t, ChannelMessage, rule.Source.Type)

	target, ok := rule.Target.(*HiveTextModelOnlyTarget)
	assert.True(t, ok)
	assert.Equal(t, "key123", target.ApiKey)
	assert.Equal(t, "https://example.com", target.ModelUrl)
	assert.Equal(t, map[string]int{"violence": 2}, target.Thresholds)
	assert.Equal(t, "hive/text-model-only", rule.RuleType())
}

func TestRuleUnmarshalHttpBeforePublish(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "http/before-publish",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"url": "https://example.com", "format": "json"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*HttpBeforePublishTarget)
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", target.Url)
	assert.Equal(t, Json, target.Format)
}

func TestRuleUnmarshalAwsLambdaBeforePublish(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "aws/lambda/before-publish",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"region": "us-east-1", "functionName": "myFunc", "authentication": {"authenticationMode": "assumeRole", "assumeRoleArn": "arn:aws:iam::123:role/r"}}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*AwsLambdaBeforePublishTarget)
	assert.True(t, ok)
	assert.Equal(t, "us-east-1", target.Region)
	assert.Equal(t, "myFunc", target.FunctionName)
}

func TestRuleUnmarshalHiveDashboard(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "hive/dashboard",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"apiKey": "key123", "checkWatchLists": true}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*HiveDashboardTarget)
	assert.True(t, ok)
	assert.Equal(t, "key123", target.ApiKey)
	assert.NotNil(t, target.CheckWatchLists)
	assert.True(t, *target.CheckWatchLists)
}

func TestRuleUnmarshalBodyguardTextModeration(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "bodyguard/text-moderation",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"apiKey": "key123", "channelId": "chan1", "apiUrl": "https://api.bg", "defaultLanguage": "en"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*BodyguardTextModerationTarget)
	assert.True(t, ok)
	assert.Equal(t, "key123", target.ApiKey)
	assert.Equal(t, "chan1", target.ChannelId)
	assert.Equal(t, "https://api.bg", target.ApiUrl)
	assert.Equal(t, "en", target.DefaultLanguage)
}

func TestRuleUnmarshalTisaneTextModeration(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "tisane/text-moderation",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"apiKey": "key123", "modelUrl": "https://tisane.ai", "thresholds": {"profanity": 1}, "defaultLanguage": "fr"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*TisaneTextModerationTarget)
	assert.True(t, ok)
	assert.Equal(t, "key123", target.ApiKey)
	assert.Equal(t, "https://tisane.ai", target.ModelUrl)
	assert.Equal(t, map[string]int{"profanity": 1}, target.Thresholds)
	assert.Equal(t, "fr", target.DefaultLanguage)
}

func TestRuleUnmarshalAzureTextModeration(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "azure/text-moderation",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"apiKey": "key123", "endpoint": "https://azure.example.com", "thresholds": {"Hate": 4}}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	target, ok := rule.Target.(*AzureTextModerationTarget)
	assert.True(t, ok)
	assert.Equal(t, "key123", target.ApiKey)
	assert.Equal(t, "https://azure.example.com", target.Endpoint)
	assert.Equal(t, map[string]int{"Hate": 4}, target.Thresholds)
}

func TestRuleUnmarshalUnsupportedTarget(t *testing.T) {
	data := []byte(`{
		"id": "rule123",
		"appId": "app123",
		"version": "1.0",
		"status": "enabled",
		"created": 1234567890,
		"modified": 1234567890,
		"ruleType": "some/future-type",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"url": "https://example.com", "customField": "value"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	assert.Equal(t, "rule123", rule.ID)
	assert.Equal(t, "some/future-type", rule.RuleType())

	target, ok := rule.Target.(*UnsupportedTarget)
	assert.True(t, ok)
	assert.NotNil(t, target.Raw)

	// Verify the raw JSON is preserved
	var rawTarget map[string]any
	err = json.Unmarshal(target.Raw, &rawTarget)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", rawTarget["url"])
	assert.Equal(t, "value", rawTarget["customField"])
}

func TestRuleUnmarshalInvocationModeAndBeforePublishConfigPreserved(t *testing.T) {
	data := []byte(`{
		"id": "rule1",
		"appId": "app1",
		"version": "1.0",
		"status": "enabled",
		"created": 100,
		"modified": 200,
		"ruleType": "http/before-publish",
		"invocationMode": "BEFORE_PUBLISH",
		"beforePublishConfig": {
			"retryTimeout": 8000,
			"maxRetries": 5,
			"failedAction": "PUBLISH",
			"tooManyRequestsAction": "FAIL"
		},
		"chatRoomFilter": "^rooms:.*",
		"requestMode": "single",
		"source": {"channelFilter": "test", "type": "channel.message"},
		"target": {"url": "https://example.com"}
	}`)

	var rule Rule
	err := json.Unmarshal(data, &rule)
	assert.NoError(t, err)

	assert.Equal(t, BeforePublish, rule.InvocationMode)
	assert.Equal(t, "^rooms:.*", rule.ChatRoomFilter)
	assert.NotNil(t, rule.BeforePublishConfig)
	assert.Equal(t, 8000, rule.BeforePublishConfig.RetryTimeout)
	assert.Equal(t, 5, rule.BeforePublishConfig.MaxRetries)
	assert.Equal(t, FailedActionPublish, rule.BeforePublishConfig.FailedAction)
	assert.Equal(t, TooManyRequestsFail, rule.BeforePublishConfig.TooManyRequestsAction)
}
