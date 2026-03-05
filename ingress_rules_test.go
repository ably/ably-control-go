package control

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleIngressMongo(t *testing.T) {
	target := &IngressMongoTarget{
		Url:                      "mongodb://coco:nut@coco.io:27017",
		Database:                 "coconut",
		Collection:               "coconut",
		Pipeline:                 `[{"$set": {"_ablyChannel": "myChannel"}}]`,
		FullDocument:             "off",
		FullDocumentBeforeChange: "off",
		PrimarySite:              "us-east-1-A",
	}

	testIngressRule(t, target)
}

func TestRuleIngressPostgresOutbox(t *testing.T) {
	target := &IngressPostgresOutboxTarget{
		Url:               "postgres://user:password@example.com:5432/your-database-name",
		OutboxTableSchema: "public",
		OutboxTableName:   "outbox",
		NodesTableSchema:  "public",
		NodesTableName:    "nodes",
		SslMode:           "prefer",
		SslRootCert:       "-----BEGIN CERTIFICATE----- MIIFiTCCA3GgAwIBAgIUYO1Lomxzj7VRawWwEFiQht9OLpUwDQYJKoZIhvcNAQEL BQAwTDELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE1pY2hpZ2FuMQ8wDQYDVQQHDAZX ...snip... TOfReTlUQzgpXRW5h3n2LVXbXQhPGcVitb88Cm2R8cxQwgB1VncM8yvmKhREo2tz 7Y+sUx6eIl4dlNl9kVrH1TD3EwwtGsjUNlFSZhg= -----END CERTIFICATE-----",
		PrimarySite:       "us-east-1-A",
	}

	testIngressRule(t, target)
}

func testIngressRule(t *testing.T, target IngressTarget) {
	t.Helper()
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	rule := NewIngressRule{
		Status: "enabled",
		Target: target,
	}

	r, err := client.CreateIngressRule(app.ID, &rule)
	require.NoError(t, err)
	assert.Equal(t, rule.Target.TargetType(), r.Target.TargetType())
	assert.NotEmpty(t, r.ID)
	assert.NotEmpty(t, r.AppID)
	assert.NotEmpty(t, r.Version)
	assert.NotEmpty(t, r.Status)
	assert.NotEmpty(t, r.Created)
	assert.NotEmpty(t, r.Modified)

	r2, err := client.IngressRule(app.ID, r.ID)
	require.NoError(t, err)
	assert.Equal(t, r, r2)

	err = client.DeleteIngressRule(app.ID, r.ID)
	assert.NoError(t, err)
}

func TestIngressMongoTargetIncludesNewFields(t *testing.T) {
	target := IngressMongoTarget{
		Url:                 "mongodb://localhost",
		Database:            "testdb",
		Collection:          "testcol",
		Watch:               "collection",
		ProvisionedCapacity: 1.5,
	}

	data, err := json.Marshal(target)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Equal(t, "collection", raw["watch"])
	assert.Equal(t, 1.5, raw["provisionedCapacity"])
}

func TestIngressRuleUnmarshalIncludesLinks(t *testing.T) {
	data := `{
		"id": "rule789",
		"appId": "app456",
		"status": "enabled",
		"created": 1700000000,
		"modified": 1700001000,
		"_links": {
			"self": "https://control.ably.net/v1/apps/app456/rules/rule789"
		},
		"ruleType": "ingress/mongodb",
		"target": {
			"url": "mongodb://localhost",
			"database": "testdb",
			"collection": "testcol"
		}
	}`

	var rule IngressRule
	err := json.Unmarshal([]byte(data), &rule)
	assert.NoError(t, err)
	assert.NotNil(t, rule.Links)
	assert.Equal(t, "https://control.ably.net/v1/apps/app456/rules/rule789", rule.Links["self"])
}
