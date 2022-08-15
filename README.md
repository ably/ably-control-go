# [Ably Control Go](https://ably.com/)

[![Go Reference](https://pkg.go.dev/badge/github.com/ably/ably-control-go.svg)](https://pkg.go.dev/github.com/ably/ably-control-go)

_[Ably](https://ably.com) is the platform that powers synchronized digital experiences in realtime. Whether attending an event in a virtual venue, receiving realtime financial information, or monitoring live car performance data – consumers simply expect realtime digital experiences as standard. Ably provides a suite of APIs to build, extend, and deliver powerful digital experiences in realtime for more than 250 million devices across 80 countries each month. Organizations like Bloomberg, HubSpot, Verizon, and Hopin depend on Ably’s platform to offload the growing complexity of business-critical realtime data synchronization at global scale. For more information, see the [Ably documentation](https://ably.com/documentation)._

## Overview

This is a Go client library for the [Ably control API](https://ably.com/docs/control-api).

Use the Control API to manage your applications, namespaces, keys, queues, rules, and more.
Detailed information on using this API can be found in the Ably
[Control API documentation](https://ably.com/docs/api/control-api). Control API is currently in Preview.

## OpenAPI

An OpenAPI document for the control API can be found at https://github.com/ably/open-specs. This repo
is not generated from the OpenAPI but instead written manually. This is because the OpenAPI generator
did not produce an erganomic or easy to use library.

## Installation

```bash
~ $ go get -u github.com/ably/ably-go/ably
```

## Examples

### Create a client

```go
token := os.Getenv("ABLY_ACCOUNT_TOKEN")
client, _, err := control.NewClient(token)
if err != nil {
	panic(err)
}

fmt.Println(client)
```

### Get account and user info

```go
me, err := client.Me()
if err != nil {
	panic(err)
}

println(me.User.Email)
```

### List apps

```go
apps, err := client.Apps()
if err != nil {
	panic(err)
}

for _, app := range apps {
	fmt.Println(app.Name)
}
```

### Create app

```go
newapp := control.App{
	Name:    "Foo",
	TLSOnly: true,
}
app, err := client.CreateApp(&newapp)
if err != nil {
	panic(err)
}

fmt.Println(app.Name)
```

### Update app

```go
app.Name = "Bar"
app, err := client.UpdateApp(app.ID, &app)
if err != nil {
	panic(err)
}
```

### Create key

```go
newkey := control.NewKey{
	Name:       "KeyName",
	Capability: map[string][]string{"a": {"subscribe"}},
}
key, err := client.CreateKey(app.ID, &newkey)
if err != nil {
	panic(err)
}

println(key.Key)
```

### Create rule

```go
target := &control.PulsarTarget{
	RoutingKey:    "aaaaa",
	Topic:         "my-tenant/my-namespace/my-topic",
	ServiceURL:    "pulsar://test.com:1234",
	TlsTrustCerts: []string{"-----BEGIN CERTIFICATE-----\naaaaa\n-----END CERTIFICATE-----"},
	Authentication: control.PulsarAuthentication{
		AuthenticationMode: "token",
		Token:              "1234",
	},
	Enveloped: true,
	Format:    control.Json,
}

newrule := control.NewRule{
	Status:      "enabled",
	RequestMode: control.Single,
	Source: control.Source{
		ChannelFilter: "aaa",
		Type:          control.ChannelMessage,
	},
	Target: target,
}

rule, err := client.CreateRule(app.ID, &newrule)
if err != nil {
	panic(err)
}

fmt.Println(rule.ID)
```

## Supported Versions of Go

Whenever a new version of Go is released, Ably adds support for that version. The [Go Release Policy](https://golang.org/doc/devel/release#policy)
supports the last two major versions. This SDK follows the same policy of supporting the last two major versions of Go.

## Contributing

For guidance on how to contribute to this project, see [CONTRIBUTING.md](CONTRIBUTING.md).
