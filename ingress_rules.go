package control

import (
	"encoding/json"
	"fmt"
)

type NewIngressRuleNoJson NewIngressRule

// IngressRule is a struct representing an Ably Ingress rule.
type IngressRule struct {
	// The rule ID.
	ID string `json:"id,omitempty"`
	// The Ably application ID.
	AppID string `json:"appId,omitempty"`
	// API version. Events and the format of their payloads are versioned.
	// Please see the Events documentation. https://ably.com/documentation/general/events
	Version string `json:"version,omitempty"`
	// The status of the rule. Rules can be enabled or disabled.
	Status string `json:"status,omitempty"`
	// Unix timestamp representing the date and time of creation of the rule.
	Created int `json:"created"`
	// Unix timestamp representing the date and time of last modification of the rule.
	Modified int `json:"modified"`
	// The rule target.
	Target IngressTarget `json:"target"`
}

// IngressRuleType gets the type of target this rule has.
func (r *IngressRule) IngressRuleType() string {
	return r.Target.TargetType()
}

func (r *IngressRule) UnmarshalJSON(data []byte) error {
	var raw rawIngressRule
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	switch raw.RuleType {
	case "ingress/mongodb":
		var t IngressMongoTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "ingress-postgres-outbox":
		var t IngressPostgresOutboxTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	default:
		return fmt.Errorf("unknown rule type \"%s\"", raw.RuleType)
	}

	if err != nil {
		return err
	}

	r.ID = raw.ID
	r.AppID = raw.AppID
	r.Version = raw.Version
	r.Status = raw.Status
	r.Created = raw.Created
	r.Modified = raw.Modified

	return nil
}

type rawIngressRule struct {
	ID       string          `json:"id,omitempty"`
	AppID    string          `json:"appId,omitempty"`
	Version  string          `json:"version,omitempty"`
	Status   string          `json:"status,omitempty"`
	Created  int             `json:"created"`
	Modified int             `json:"modified"`
	RuleType string          `json:"ruleType,omitempty"`
	Target   json.RawMessage `json:"target"`
}

// IngressRuleType gets the type of target this rule has.
func (r *NewIngressRule) IngressRuleType() string {
	return r.Target.TargetType()
}

type IngressTarget interface {
	// TargetType returns the kind of target.
	TargetType() string
}

type rawNewIngressRule struct {
	RuleType string `json:"ruleType,omitempty"`
	*NewIngressRuleNoJson
}

// NewRule is used to create a new rule.
type NewIngressRule struct {
	// The status of the rule. Rules can be enabled or disabled.
	Status string `json:"status,omitempty"`
	// The rule target.
	Target IngressTarget `json:"target"`
}

func (r *NewIngressRule) MarshalJSON() ([]byte, error) {
	raw := rawNewIngressRule{
		RuleType: r.Target.TargetType(), NewIngressRuleNoJson: (*NewIngressRuleNoJson)(r)}

	return json.Marshal(&raw)
}

// IngressMongoTarget is the type used for MongoDB Ingress rules.
type IngressMongoTarget struct {
	// The URL of the MongoDB server.
	Url string `json:"url,omitempty"`
	// Watch is the collection to watch.
	Database string `json:"database,omitempty"`
	// The collection to watch.
	Collection string `json:"collection,omitempty"`
	// The pipeline to use.
	Pipeline string `json:"pipeline,omitempty"`
	// FullDocument controls how the full document is sent.
	FullDocument string `json:"fullDocument,omitempty"`
	// FullDocumentBeforeChange controls how the full document is sent.
	FullDocumentBeforeChange string `json:"fullDocumentBeforeChange,omitempty"`
	// The primary site.
	PrimarySite string `json:"primarySite,omitempty"`
}

// IngressMongoTarget implements the Target interface.
func (s *IngressMongoTarget) TargetType() string {
	return "ingress/mongodb"
}

// RuleType gets the type of target this rule has.
func (r *IngressRule) RuleType() string {
	return r.Target.TargetType()
}

type IngressPostgresOutboxTarget struct {
	// The URL for your Postgres database.
	Url string `json:"url,omitempty"`
	// Schema for the outbox table in your database which allows for the
	// reliable publication of an ordered sequence of change event messages
	// over Ably.
	OutboxTableSchema string `json:"outboxTableSchema,omitempty"`
	// Name for the outbox table.
	OutboxTableName string `json:"outboxTableName,omitempty"`
	// Schema for the nodes table in your database to allow for operation as a
	// cluster to provide fault tolerance.
	NodesTableSchema string `json:"nodesTableSchema,omitempty"`
	// Name for the nodes table.
	NodesTableName string `json:"nodesTableName,omitempty"`
	// Determines the level of protection provided by the SSL connection.
	// Options are: prefer, require, verify-ca, verify-full;
	// default value is prefer.
	SslMode string `json:"sslMode,omitempty"`
	// Optional. Specifies the SSL certificate authority (CA) certificates.
	// Required if SSL mode is verify-ca or verify-full.
	SslRootCert string `json:"sslRootCert,omitempty"`
	//The primary data center in which to run the integration rule.
	PrimarySite string `json:"primarySite,omitempty"`
}

// IngressPostgresTarget implements the Target interface.
func (s *IngressPostgresOutboxTarget) TargetType() string {
	return "ingress-postgres-outbox"
}

// Creates an Ingress rule for the application with the specified application ID.
func (c *Client) CreateIngressRule(appID string, rule *NewIngressRule) (IngressRule, error) {
	var out IngressRule
	err := c.request("POST", "/apps/"+appID+"/rules", &rule, &out)
	return out, err
}

// Lists the rules for the application specified by the application ID.
func (c *Client) IngressRules(appID string) ([]IngressRule, error) {
	var rules []IngressRule
	err := c.request("GET", "/apps/"+appID+"/rules", nil, &rules)
	return rules, err
}

// Returns the ingess rule specified by the rule ID, for the application specified by application ID.
func (c *Client) IngressRule(appID, ruleID string) (IngressRule, error) {
	var rule IngressRule
	err := c.request("GET", "/apps/"+appID+"/rules/"+ruleID, nil, &rule)
	return rule, err
}

// Updates the rule specified by the rule ID, for the application specified by application ID.
func (c *Client) UpdateIngressRule(appID, ruleID string, rule *NewIngressRule) (IngressRule, error) {
	var out IngressRule
	err := c.request("PATCH", "/apps/"+appID+"/rules/"+ruleID, &rule, &out)
	return out, err
}

// Deletes the rule specified by the rule ID, for the application specified by application ID.
func (c *Client) DeleteIngressRule(appID, ruleID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/rules/"+ruleID, nil, nil)
	return err
}
