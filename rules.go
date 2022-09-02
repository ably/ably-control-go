package control

import (
	"encoding/json"
	"fmt"
)

type NewRuleNoJson NewRule

// PularAuthenticationMode is an enum of authentication modes used by Pulsar rules.
type PularAuthenticationMode string

// AuthToken AuthenticationMode.
const AuthToken PularAuthenticationMode = "token"

// SaslMechanism is the hash type used for Sasl authentication.
type SaslMechanism string

// Plain do not use a hash.
const Plain SaslMechanism = "plain"

// Scram_sha_256 use sha256 hashes.
const Scram_sha_256 SaslMechanism = "scra-sha-256"

// Scram_sha_512 use sha512 hashes.
const Scram_sha_512 SaslMechanism = "scra-sha-512"

// Format is the format used for encoding.
type Format string

// Json encodes using json.
const Json Format = "json"

// MsgPack encodes using message pack.
const MsgPack Format = "msgpack"

// SourceType is the type of messages a source applies to.
type SourceType string

// ChannelMessage represents message published to a channel.
const ChannelMessage SourceType = "channel.message"

// ChannelPresence represents presence events on a channel.
const ChannelPresence SourceType = "channel.presence"

// ChannelLifeCycle represents channel lifecycle events.
const ChannelLifeCycle SourceType = "channel.lifecycle"

// ChannelOccupancy representing channel occupancy events.
const ChannelOccupancy SourceType = "channel.occupancy"

// RequestMode is a source's request mode.
type RequestMode string

// Single is the Single Request Mode
const Single RequestMode = "single"

// Batch is the Batch Request Mode
const Batch RequestMode = "batch"

// Rule is a struct representing an Ably rule.
type Rule struct {
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
	// RequestMode. You can read more about the difference between single and batched
	// events in the Ably documentation. https://ably.com/documentation/general/events#batching
	RequestMode RequestMode `json:"requestMode,omitempty"`
	// The rule source.
	Source Source `json:"source"`
	// The rule target.
	Target Target `json:"target"`
}

// RuleType gets the type of target this rule has.
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
		if err != nil {
			fmt.Println(err, string(raw.Target))
		}
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
		var t AwsKinesisTarget
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
	RequestMode RequestMode     `json:"requestMode,omitempty"`
	Source      Source          `json:"source"`
	Target      json.RawMessage `json:"target"`
}

// RuleType gets the type of target this rule has.
func (r *NewRule) RuleType() string {
	return r.Target.TargetType()
}

// Source controls how a rule gets data from channels.
type Source struct {
	// ChannelFilter allows you to filter your rule based on a regular expression that is matched against the complete channel name.
	// Leave this empty if you want the rule to apply to all channels.
	ChannelFilter string `json:"channelFilter,omitempty"`
	// Type controls the type of messages that are sent to the rule.
	Type SourceType `json:"type,omitempty"`
}

// The Target interface is implemented by targets and
// allows querying what kind of target they are.
type Target interface {
	// TargetType returns the kind of target.
	TargetType() string
}

// Headers that can be used for some rule kinds.
type Header struct {
	// The name of the header.
	Name string `json:"name,omitempty"`
	// The value of the header.
	Value string `json:"value,omitempty"`
}

type rawNewRule struct {
	RuleType string `json:"ruleType,omitempty"`
	*NewRuleNoJson
}

// NewRule is used to create a new rule.
type NewRule struct {
	// The status of the rule. Rules can be enabled or disabled.
	Status string `json:"status,omitempty"`
	// RequestMode. You can read more about the difference between single and batched
	// events in the Ably documentation. https://ably.com/documentation/general/events#batching
	RequestMode RequestMode `json:"requestMode,omitempty"`
	// The rule source.
	Source Source `json:"source"`
	// The rule target.
	Target Target `json:"target"`
}

func (r *NewRule) MarshalJSON() ([]byte, error) {
	raw := rawNewRule{
		RuleType: r.Target.TargetType(), NewRuleNoJson: (*NewRuleNoJson)(r)}

	return json.Marshal(&raw)
}

// PulsarAuthentication is used to authenticate for Pulsar rules
type PulsarAuthentication struct {
	// Authentication mode.
	AuthenticationMode PularAuthenticationMode `json:"authenticationMode,omitempty"`
	// The JWT string.
	Token string `json:"token,omitempty"`
}

// SASL (Simple Authentication Security Layer) / SCRAM (Salted Challenge Response Authentication Mechanism)
// uses usernames and passwords stored in ZooKeeper. Credentials are created during installation.
// See documentation on configuring SCRAM.
type Sasl struct {
	// The hash type to use.
	Mechanism SaslMechanism `json:"mechanism,omitempty"`
	// Kafka login credential.
	Username string `json:"username,omitempty"`
	// Kafka login credential.
	Password string `json:"password,omitempty"`
}

// KafkaAuthentication are params used for authenticating with Kafka.
type KafkaAuthentication struct {
	// SASL (Simple Authentication Security Layer) / SCRAM (Salted Challenge Response Authentication Mechanism)
	// uses usernames and passwords stored in ZooKeeper. Credentials are created during installation.
	// See documentation on configuring SCRAM.
	Sasl Sasl `json:"sasl"`
}

// AwsAuthentication are the params used for authenticating with AWS.
type AwsAuthentication struct {
	// Authentication can be any supported AWS authentication type.
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

// AwsAuthenticationType is an interface implemented by the different structs that
// can be used in AWS authentication.
type AwsAuthenticationType interface {
	AuthenticationMode() string
}

// AuthenticationModeAssumeRole is the assume role authentication mode for AWS.
type AuthenticationModeAssumeRole struct {
	// If you are using the "ARN of an assumable role" authentication method,
	// this is your Assume Role ARN. See this Ably knowledge base article for details.
	// https://ably.com/documentation/general/aws-authentication/
	AssumeRoleArn string `json:"assumeRoleArn,omitempty"`
}

// AuthenticationModeAssumeRole implements AwsAuthenticationType.
func (a *AuthenticationModeAssumeRole) AuthenticationMode() string {
	return "assumeRole"
}

// AuthenticationModeCredentials is the credentials authentication mode for AWS.
type AuthenticationModeCredentials struct {
	// The AWS key ID for the AWS IAM user. See this Ably knowledge base article for details.
	// https://ably.com/documentation/general/aws-authentication/
	AccessKeyId string `json:"accessKeyId,omitempty"`
	// The AWS secret key for the AWS IAM user. See this Ably knowledge base article for details.
	// https://ably.com/documentation/general/aws-authentication/
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// AuthenticationModeCredentials implements the AwsAuthenticationType interface.
func (a *AuthenticationModeCredentials) AuthenticationMode() string {
	return "credentials"
}

// PulsarTarget is the type used for Pular rules.
type PulsarTarget struct {
	// The optional routing key (partition key) used
	// to publish messages. Supports interpolation as described in the Ably FAQs.
	// https://faqs.ably.com/what-is-the-format-of-the-routingkey-for-an-amqp-or-kinesis-reactor-rule
	RoutingKey string `json:"routingKey,omitempty"`
	// A Pulsar topic. This is a named channel for transmission of messages between producers and consumers.
	// The topic has the form: {persistent|non-persistent}://tenant/namespace/topic
	Topic string `json:"topic,omitempty"`
	// The URL of the Pulsar cluster in the form pulsar://host:port or pulsar+ssl://host:port.
	ServiceURL string `json:"serviceUrl,omitempty"`
	// All connections to a Pulsar endpoint require TLS. The tlsTrustCerts option
	// allows you to configure different or additional trust anchors for those TLS
	// connections. This enables server verification. You can specify an optional
	// list of trusted CA certificates to use to verify the TLS certificate presented
	// by the Pulsar cluster. Each certificate should be encoded in PEM format.
	TlsTrustCerts []string `json:"tlsTrustCerts"`
	// Pulsar supports authenticating clients using security tokens that are based on JSON Web Tokens.
	Authentication PulsarAuthentication `json:"authentication"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// PulsarTarget implements the Target interface.
func (s *PulsarTarget) TargetType() string {
	return "pulsar"
}

// KafkaTarget is the type for Kafka targets.
type KafkaTarget struct {
	// The Kafka partition key. This is used to determine which partition a message should be routed to,
	// where a topic has been partitioned. routingKey should be in the format topic:key where topic is
	// the topic to publish to, and key is the value to use as the message key.
	RoutingKey string `json:"routingKey,omitempty"`
	// This is an array of brokers that host your Kafka partitions. Each broker is specified
	// using the format host, host:port or ip:port.
	Brokers []string `json:"brokers,omitempty"`
	// The Kafka authentication mechanism.
	Authentication KafkaAuthentication `json:"auth"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// KafkaTarget implements the Target interface.
func (s *KafkaTarget) TargetType() string {
	return "kafka"
}

// AmqpExternalTarget is the type used for AMQP/external rules.
type AmqpExternalTarget struct {
	// The webhook URL that Ably will POST events to.
	Url string `json:"url,omitempty"`
	// The AMQP routing key. The routing key is used by the AMQP exchange to route messages to a
	// physical queue. See this Ably knowledge base article for details.
	// https://knowledge.ably.com/what-is-the-format-of-the-routingkey-for-an-amqp-or-kinesis-reactor-rule
	RoutingKey string `json:"routingKey,omitempty"`
	// Reject delivery of the message if the route does not exist, otherwise fail silently.
	MandatoryRoute bool `json:"mandatoryRoute"`
	// Marks the message as persistent, instructing the broker to write it to disk if it is in a durable queue.
	PersistentMessages bool `json:"persistentMessages"`
	//You can optionally override the default TTL on a queue and specify a TTL in minutes for
	//messages to be persisted. It is unusual to change the default TTL, so if this field is
	//left empty, the default TTL for the queue will be used.
	MessageTTL int `json:"messageTtl,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// AmqpExternalTarget implements the Target interface.
func (s *AmqpExternalTarget) TargetType() string {
	return "amqp/external"
}

// AmqpTarget is the type used for AMQP rules.
type AmqpTarget struct {
	// The ID of the Ably queue for input.
	QueueID string `json:"queueId,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// AmqpTarget implements the Target interface.
func (s *AmqpTarget) TargetType() string {
	return "amqp"
}

// AwsSqsTarget is the type used for aws/sqs rules.
type AwsSqsTarget struct {
	// The region is which AWS SQS is hosted. See the AWS documentation for more detail.
	// https://docs.aws.amazon.com/general/latest/gr/rande.html#lambda_region
	Region string `json:"region,omitempty"`
	// Your AWS account ID.
	AwsAccountID string `json:"awsAccountId,omitempty"`
	// The AWS SQS queue name.
	QueueName string `json:"queueName,omitempty"`
	// Authentication details.
	Authentication AwsAuthentication `json:"authentication"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// AwsSqsTarget implements the Target interface.
func (s *AwsSqsTarget) TargetType() string {
	return "aws/sqs"
}

// AwsKinesisTarget is the type used for aws/kinesis rules.
type AwsKinesisTarget struct {
	// The region is which AWS Kinesis is hosted. See the AWS documentation for more detail.
	// https://docs.aws.amazon.com/general/latest/gr/rande.html#lambda_region
	Region string `json:"region,omitempty"`
	// The name of your AWS Kinesis Stream.
	StreamName string `json:"streamName,omitempty"`
	// The AWS Kinesis partition key. The partition key is used by Kinesis to route messages
	// to one of the stream shards. See this Ably knowledge base article for details.
	// https://knowledge.ably.com/what-is-the-format-of-the-routingkey-for-an-amqp-or-kinesis-reactor-rule
	PartitionKey string `json:"partitionKey,omitempty"`
	// Authentication details.
	Authentication AwsAuthentication `json:"authentication"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// AwsKinesisTarget implements the Target interface.
func (s *AwsKinesisTarget) TargetType() string {
	return "aws/kinesis"
}

// AwsLambdaTarget is the type used for aws/lambda rules.
type AwsLambdaTarget struct {
	// The region is which your AWS Lambda Function is hosted. See the AWS documentation for more detail.
	// https://docs.aws.amazon.com/general/latest/gr/rande.html#lambda_region
	Region string `json:"region,omitempty"`
	// The name of your AWS Lambda Function.
	FunctionName string `json:"functionName,omitempty"`
	// Authentication details.
	Authentication AwsAuthentication `json:"authentication"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
}

// AwsLambdaTarget implements the Target interface.
func (s *AwsLambdaTarget) TargetType() string {
	return "aws/lambda"
}

// HttpGoogleCloudFunctionTarget is type type used for http/goole-cloud-function targets.
type HttpGoogleCloudFunctionTarget struct {
	// The region in which your Google Cloud Function is hosted. See the Google documentation for more details.
	// https://cloud.google.com/compute/docs/regions-zones/
	Region string `json:"region,omitempty"`
	// The project ID for your Google Cloud Project that was generated when you created your project.
	ProjectID string `json:"projectId,omitempty"`
	// The name of your Google Cloud Function.
	FunctionName string `json:"functionName,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// The signing key ID for use in batch mode. Ably will optionally sign the payload using
	// an API key ensuring your servers can validate the payload using the private API key.
	// See the webhook security docs for more information.
	// https://ably.com/documentation/general/events#security
	SigningKeyID string `json:"signingKeyId"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// HttpGoogleCloudFunctionTarget implements the Target interface.
func (s *HttpGoogleCloudFunctionTarget) TargetType() string {
	return "http/google-cloud-function"
}

// HttpAzureFunctionTarget is the type used for http/azure-function targets.
type HttpAzureFunctionTarget struct {
	// The Microsoft Azure Application ID. You can find your Microsoft Azure Application
	// ID as shown in this article. https://dev.applicationinsights.io/documentation/Authorization/API-key-and-App-ID
	AzureAppID string `json:"azureAppId,omitempty"`
	// The name of your Microsoft Azure Function.
	AzureFunctionName string `json:"azureFunctionName,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// The signing key ID for use in batch mode. Ably will optionally sign the payload using an API key
	// ensuring your servers can validate the payload using the private API key. See the webhook security
	// docs for more information. https://ably.com/documentation/general/events#security
	SigningKeyID string `json:"signingKeyId"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// HttpAzureFunctionTarget implements the Target interface.
func (s *HttpAzureFunctionTarget) TargetType() string {
	return "http/azure-function"
}

// HttpCloudfareWorkerTarget is the type used for http/cloudflare-worker rules.
type HttpCloudfareWorkerTarget struct {
	// The webhook URL that Ably will POST events to.
	Url string `json:"url,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// The signing key ID for use in batch mode. Ably will optionally sign the payload using an API key
	// ensuring your servers can validate the payload using the private API key. See the webhook security
	// docs for more information. https://ably.com/documentation/general/events#security
	SigningKeyID string `json:"signingKeyId"`
}

// HttpCloudfareWorkerTarget implements the Target interface.
func (s *HttpCloudfareWorkerTarget) TargetType() string {
	return "http/cloudflare-worker"
}

// HttpZapierTarget is the type used for http/zapier rules.
type HttpZapierTarget struct {
	// The webhook URL that Ably will POST events to.
	Url string `json:"url,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// The signing key ID for use in batch mode. Ably will optionally sign the payload using an API key
	// ensuring your servers can validate the payload using the private API key. See the webhook security
	// docs for more information. https://ably.com/documentation/general/events#security
	SigningKeyID string `json:"signingKeyId"`
}

// HttpZapierTarget implements the Target interface.
func (s *HttpZapierTarget) TargetType() string {
	return "http/zapier"
}

// HttpIftttTarget is the type used for http/ifttt rules.
type HttpIftttTarget struct {
	// The key in the Webhook Service Documentation page of your IFTTT account.
	WebhookKey string `json:"webhookKey,omitempty"`
	// The Event name is used to identify the IFTTT applet that will receive the Event,
	// make sure the name matches the name of the IFTTT applet.
	EventName string `json:"eventName,omitempty"`
}

// HttpIftttTarget implements the Target interface.
func (s *HttpIftttTarget) TargetType() string {
	return "http/ifttt"
}

// HttpTarget is the type used for http rules.
type HttpTarget struct {
	// The webhook URL that Ably will POST events to.
	Url string `json:"url,omitempty"`
	// If you have additional information to send, you'll need to include the relevant headers.
	Headers []Header `json:"headers,omitempty"`
	// The signing key ID for use in batch mode. Ably will optionally sign the payload using an
	// API key ensuring your servers can validate the payload using the private API key.
	// See the webhook security docs for more information.
	// https://ably.com/documentation/general/events#security
	SigningKeyID string `json:"signingKeyId"`
	// Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message
	// and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a
	// Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by
	// unchecking "Enveloped" when setting up the rule.
	Enveloped bool `json:"enveloped"`
	// JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding.
	Format Format `json:"format,omitempty"`
}

// HttpTarget implements the Target interface.
func (s *HttpTarget) TargetType() string {
	return "http"
}

// Lists the rules for the application specified by the application ID.
func (c *Client) Rules(appID string) ([]Rule, error) {
	var rules []Rule
	err := c.request("GET", "/apps/"+appID+"/rules", nil, &rules)
	return rules, err
}

// Returns the rule specified by the rule ID, for the application specified by application ID.
func (c *Client) Rule(appID, ruleID string) (Rule, error) {
	var rule Rule
	err := c.request("GET", "/apps/"+appID+"/rules/"+ruleID, nil, &rule)
	return rule, err
}

// Creates a rule for the application with the specified application ID.
func (c *Client) CreateRule(appID string, rule *NewRule) (Rule, error) {
	var out Rule
	err := c.request("POST", "/apps/"+appID+"/rules", &rule, &out)
	return out, err
}

// Updates the rule specified by the rule ID, for the application specified by application ID.
func (c *Client) UpdateRule(appID, ruleID string, rule *NewRule) (Rule, error) {
	var out Rule
	err := c.request("PATCH", "/apps/"+appID+"/rules/"+ruleID, &rule, &out)
	return out, err
}

// Deletes the rule specified by the rule ID, for the application specified by application ID.
func (c *Client) DeleteRule(appID, ruleID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/rules/"+ruleID, nil, nil)
	return err
}
