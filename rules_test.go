package control

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRulePulsar(t *testing.T) {
	target := &PulsarTarget{
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
	}

	testRule(t, target)
}

func TestRuleKafka(t *testing.T) {
	target := &KafkaTarget{
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
	}

	testRule(t, target)
}

func TestRuleAmqpExtrernal(t *testing.T) {
	target := &AmqpExternalTarget{
		Url:                "amqps://test.com",
		RoutingKey:         "key",
		Exchange:           "exchange",
		MandatoryRoute:     true,
		PersistentMessages: true,
		MessageTTL:         50,
		Headers:            []Header{{Name: "a", Value: "b"}},
		Enveloped:          true,
		Format:             Json,
	}

	testRule(t, target)
}

func TestRuleAmqp(t *testing.T) {
	target := &AmqpTarget{
		Headers:   []Header{{Name: "a", Value: "b"}},
		Enveloped: true,
		Format:    Json,
	}

	testRule(t, target)
}

func TestRuleAwsSqs(t *testing.T) {
	target := &AwsSqsTarget{
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
	}

	testRule(t, target)
}

func TestRuleAwsKinesis(t *testing.T) {
	target := &AwsKinesisTarget{
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
	}

	testRule(t, target)
}

func TestRuleAwsLambda(t *testing.T) {
	target := &AwsLambdaTarget{
		Region:       "us-east-2",
		FunctionName: "heck",
		Authentication: AwsAuthentication{
			Authentication: &AuthenticationModeAssumeRole{
				AssumeRoleArn: "aaaaaaa",
			},
		},
		Enveloped: true,
	}

	testRule(t, target)
}

func TestRuleHttpGoogleCloudFunction(t *testing.T) {
	target := &HttpGoogleCloudFunctionTarget{
		Region:       "us",
		ProjectID:    "1234",
		FunctionName: "heck",
		Headers:      []Header{{Name: "a", Value: "b"}},
		SigningKeyID: "1234",
		Enveloped:    true,
		Format:       Json,
	}

	testRule(t, target)
}

func TestRuleHttpAzureFunction(t *testing.T) {
	target := &HttpAzureFunctionTarget{
		AzureAppID:        "11111",
		AzureFunctionName: "heck",
		Headers:           []Header{{Name: "a", Value: "b"}},
		Enveloped:         true,
		Format:            Json,
	}

	testRule(t, target)
}

func TestRuleHttpCloudfareWorker(t *testing.T) {
	target := &HttpCloudfareWorkerTarget{
		Url:     "https://test.com",
		Headers: []Header{{Name: "a", Value: "b"}},
	}

	testRule(t, target)
}

func TestRuleHttpZapier(t *testing.T) {
	target := &HttpZapierTarget{
		Url:     "https://test.com",
		Headers: []Header{{Name: "a", Value: "b"}},
	}

	testRule(t, target)
}

func TestRuleHttpIfttt(t *testing.T) {
	target := &HttpIftttTarget{
		WebhookKey: "aaa",
		EventName:  "bbb",
	}

	testRule(t, target)
}

func TestRuleHttp(t *testing.T) {
	target := &HttpTarget{
		Url:       "https://test.com",
		Headers:   []Header{{Name: "a", Value: "b"}},
		Enveloped: true,
		Format:    MsgPack,
	}

	testRule(t, target)
}

func testRule(t *testing.T, target Target) {
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
	assert.NoError(t, err)

	k, err := client.CreateKey(app.ID, &key)
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.Equal(t, rule.RequestMode, r.RequestMode)
	assert.Equal(t, rule.Source, r.Source)
	assert.Equal(t, rule.Target, r.Target)
	assert.Equal(t, rule.Target.TargetType(), r.Target.TargetType())
	assert.NotEmpty(t, r.ID)
	assert.NotEmpty(t, r.AppID)
	assert.NotEmpty(t, r.Version)
	assert.NotEmpty(t, r.Status)
	assert.NotEmpty(t, r.Created)
	assert.NotEmpty(t, r.Modified)

	r2, err := client.Rule(app.ID, r.ID)
	assert.NoError(t, err)
	assert.Equal(t, r, r2)

	err = client.DeleteRule(app.ID, r.ID)
	assert.NoError(t, err)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
