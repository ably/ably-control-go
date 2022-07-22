package control

import (
	"encoding/json"
	"fmt"
)

type Rule struct {
	ID          string `json:"id,omitempty"`
	AppID       string `json:"appId,omitempty"`
	Version     string `json:"version,omitempty"`
	Status      string `json:"status,omitempty"`
	Created     int    `json:"created"`
	Modified    int    `json:"modified"`
	RequestMode string `json:"requestMode,omitempty"`
	Source      Source `json:"source"`
	Target      Target `json:"target"`
}

func (r *Rule) RuleType() string {
	return r.Target.TargetType()
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	var raw rawRule
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	switch raw.RuleType {
	case "pulsar":
		var t PulsarTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "kafka":
		var t KafkaTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "amqp/external":
		var t AmqpExternalTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "amqp":
		var t AmqpTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "aws/sqs":
		var t AwsSqsTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "aws/kinesis":
		var t AwsKenesisTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "aws/lambda":
		var t AwsLambdaTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http/google-cloud-function":
		var t HttpGoogleCloudFunctionTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http/azure-function":
		var t HttpAzureFunctionTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http/cloudflare-worker":
		var t HttpCloudfareWorkerTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http/zapier":
		var t HttpZapierTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http/ifttt":
		var t HttpIftttTarget
		err = json.Unmarshal(raw.Target, &t)
		r.Target = &t
	case "http":
		var t HttpTarget
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
	r.RequestMode = raw.RequestMode
	r.Source = raw.Source

	return nil
}

type rawRule struct {
	ID          string          `json:"id,omitempty"`
	AppID       string          `json:"appId,omitempty"`
	Version     string          `json:"version,omitempty"`
	Status      string          `json:"status,omitempty"`
	Created     int             `json:"created"`
	Modified    int             `json:"modified"`
	RuleType    string          `json:"ruleType,omitempty"`
	RequestMode string          `json:"requestMode,omitempty"`
	Source      Source          `json:"source"`
	Target      json.RawMessage `json:"target"`
}

func (r *NewRule) RuleType() string {
	return r.Target.TargetType()
}

type Source struct {
	ChannelFilter string `json:"channelFilter,omitempty"`
	Type          string `json:"type,omitempty"`
}

type Target interface {
	TargetType() string
}

type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type rawNewRule struct {
	RuleType string `json:"ruleType,omitempty"`
	*NewRule
}

type NewRule struct {
	Status      string `json:"status,omitempty"`
	RequestMode string `json:"requestMode,omitempty"`
	Source      Source `json:"source"`
	Target      Target `json:"target"`
}

func (r *NewRule) MarshalJSON() ([]byte, error) {
	raw := rawNewRule{
		RuleType: r.Target.TargetType(), NewRule: r}

	return json.Marshal(&raw)
}

type PulsarAuthentication struct {
	AuthenticationMode string `json:"authenticationMode,omitempty"`
	Token              string `json:"token,omitempty"`
}

type Sasl struct {
	Mechanism string `json:"mechanism,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

type KafkaAuthentication struct {
	Sasl Sasl `json:"sasl"`
}

type AwsAuthentication struct {
	Authentication AwsAuthenticationType
}

func (a *AwsAuthentication) UnmarshalJSON(data []byte) error {
	var m map[string]string = make(map[string]string)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil
	}
	mode := m["authenticationMode"]
	switch mode {
	case "assumeRole":
		role := m["assumeRoleArn"]
		a.Authentication = &AuthenticationModeAssumeRole{
			AssumeRoleArn: role,
		}
	case "credentials":
		accessKeyID := m["accessKeyId"]
		secretAccessKey := m["secretAccessKey"]
		a.Authentication = &AuthenticationModeCredentials{
			AccessKeyId:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		}
	default:
		return fmt.Errorf("unknown authentication mode \"%s\"", mode)
	}

	return nil
}

func (a *AwsAuthentication) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"authenticationMode": a.Authentication.AuthenticationMode(),
	}

	switch auth := a.Authentication.(type) {
	case *AuthenticationModeAssumeRole:
		m["assumeRoleArn"] = auth.AssumeRoleArn

	case *AuthenticationModeCredentials:
		m["accessKeyId"] = auth.AccessKeyId
		m["secretAccessKey"] = auth.SecretAccessKey
	default:
		return nil, fmt.Errorf("authentication is an invalid type")
	}

	return json.Marshal(m)

}

type AwsAuthenticationType interface {
	AuthenticationMode() string
}

type AuthenticationModeAssumeRole struct {
	AssumeRoleArn string `json:"assumeRoleArn,omitempty"`
}

func (a *AuthenticationModeAssumeRole) AuthenticationMode() string {
	return "assumeRole"
}

type AuthenticationModeCredentials struct {
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

func (a *AuthenticationModeCredentials) AuthenticationMode() string {
	return "credentials"
}

type PulsarTarget struct {
	RoutingKey     string               `json:"routingKey,omitempty"`
	Topic          string               `json:"topic,omitempty"`
	ServiceURL     string               `json:"serviceUrl,omitempty"`
	TlsTrustCerts  []string             `json:"tlsTrustCerts,omitempty"`
	Authentication PulsarAuthentication `json:"authentication"`
	Enveloped      bool                 `json:"enveloped"`
	Format         string               `json:"format,omitempty"`
}

func (s *PulsarTarget) TargetType() string {
	return "pular"
}

type KafkaTarget struct {
	RoutingKey     string              `json:"routingKey,omitempty"`
	Brokers        []string            `json:"brokers,omitempty"`
	Authentication KafkaAuthentication `json:"auth"`
	Enveloped      bool                `json:"enveloped"`
	Format         string              `json:"format,omitempty"`
}

func (s *KafkaTarget) TargetType() string {
	return "kafka"
}

type AmqpExternalTarget struct {
	Url                string   `json:"url,omitempty"`
	RoutingKey         string   `json:"routingKey,omitempty"`
	MandatoryRoute     bool     `json:"mandatoryRoute"`
	PersistentMessages bool     `json:"persistentMessages"`
	MessageTTL         int      `json:"messageTtl"`
	Headers            []Header `json:"headers,omitempty"`
	Enveloped          bool     `json:"enveloped"`
	Format             string   `json:"format,omitempty"`
}

func (s *AmqpExternalTarget) TargetType() string {
	return "amqp/external"
}

type AmqpTarget struct {
	QueueID   string   `json:"queueId,omitempty"`
	Headers   []Header `json:"headers,omitempty"`
	Enveloped bool     `json:"enveloped"`
	Format    string   `json:"format,omitempty"`
}

func (s *AmqpTarget) TargetType() string {
	return "amqp"
}

type AwsSqsTarget struct {
	Region         string            `json:"region,omitempty"`
	AwsAccountID   string            `json:"awsAccountId,omitempty"`
	QueueName      string            `json:"queueName,omitempty"`
	Authentication AwsAuthentication `json:"authentication"`
	Enveloped      bool              `json:"enveloped"`
	Format         string            `json:"format,omitempty"`
}

func (s *AwsSqsTarget) TargetType() string {
	return "aws/sqs"
}

type AwsKenesisTarget struct {
	Region         string            `json:"region,omitempty"`
	StreamName     string            `json:"streamName,omitempty"`
	PartitionKey   string            `json:"partitionKey,omitempty"`
	Authentication AwsAuthentication `json:"authentication"`
	Enveloped      bool              `json:"enveloped"`
	Format         string            `json:"format,omitempty"`
}

func (s *AwsKenesisTarget) TargetType() string {
	return "aws/kinesis"
}

type AwsLambdaTarget struct {
	Region         string            `json:"region,omitempty"`
	FunctionName   string            `json:"functionName,omitempty"`
	Authentication AwsAuthentication `json:"authentication"`
	Enveloped      bool              `json:"enveloped"`
}

func (s *AwsLambdaTarget) TargetType() string {
	return "aws/lambda"
}

type HttpGoogleCloudFunctionTarget struct {
	Region       string   `json:"region,omitempty"`
	ProjectID    string   `json:"projectId,omitempty"`
	FunctionName string   `json:"functionName,omitempty"`
	Headers      []Header `json:"headers,omitempty"`
	SigningKeyID string   `json:"signingKeyId,omitempty"`
	Enveloped    bool     `json:"enveloped"`
	Format       string   `json:"format,omitempty"`
}

func (s *HttpGoogleCloudFunctionTarget) TargetType() string {
	return "http/google-cloud-function"
}

type HttpAzureFunctionTarget struct {
	AzureAppID        string   `json:"azureAppId,omitempty"`
	AzureFunctionName string   `json:"azureFunctionName,omitempty"`
	Headers           []Header `json:"headers,omitempty"`
	SigningKeyID      string   `json:"signingKeyId,omitempty"`
	Enveloped         bool     `json:"enveloped"`
	Format            string   `json:"format,omitempty"`
}

func (s *HttpAzureFunctionTarget) TargetType() string {
	return "http/azure-function"
}

type HttpCloudfareWorkerTarget struct {
	Url          string   `json:"url,omitempty"`
	Headers      []Header `json:"headers,omitempty"`
	SigningKeyID string   `json:"signingKeyId,omitempty"`
	Enveloped    bool     `json:"enveloped"`
}

func (s *HttpCloudfareWorkerTarget) TargetType() string {
	return "http/cloudflare-worker"
}

type HttpZapierTarget struct {
	Url          string   `json:"url,omitempty"`
	Headers      []Header `json:"headers,omitempty"`
	SigningKeyID string   `json:"signingKeyId,omitempty"`
}

func (s *HttpZapierTarget) TargetType() string {
	return "http/zapier"
}

type HttpIftttTarget struct {
	WebhookKey string `json:"webhookKey,omitempty"`
	EventName  string `json:"eventName,omitempty"`
}

func (s *HttpIftttTarget) TargetType() string {
	return "http/ifttt"
}

type HttpTarget struct {
	Url          string   `json:"url,omitempty"`
	Headers      []Header `json:"headers,omitempty"`
	SigningKeyID string   `json:"signingKeyId,omitempty"`
	Enveloped    bool     `json:"enveloped"`
	Format       string   `json:"format,omitempty"`
}

func (s *HttpTarget) TargetType() string {
	return "http"
}

func (c *Client) Rules(appID string) ([]Rule, error) {
	var rules []Rule
	err := c.request("GET", "/apps/"+appID+"/rules", nil, &rules)
	return rules, err
}

func (c *Client) Rule(appID, ruleID string) (Rule, error) {
	var rule Rule
	err := c.request("GET", "/apps/"+appID+"/rules/"+ruleID, nil, &rule)
	return rule, err
}

func (c *Client) CreateRule(appID string, rule NewRule) (Rule, error) {
	var out Rule
	err := c.request("POST", "/apps/"+appID+"/rules", &rule, &out)
	return out, err
}

func (c *Client) UpdateRule(appID, ruleID string, rule NewRule) (Rule, error) {
	var out Rule
	err := c.request("PATCH", "/apps/"+appID+"/rules/"+ruleID, &rule, &out)
	return out, err
}

func (c *Client) DeleteRule(appID, ruleID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/rules/"+ruleID, nil, nil)
	return err
}
