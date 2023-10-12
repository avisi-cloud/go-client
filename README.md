# Avisi Cloud API go-client

## Description

This repository contains the Golang client to communicate with the Avisi Cloud Platform API.

[Website](https://avisi.cloud) | [Documentation](https://docs.avisi.cloud)

## Installation

To make use of this go-client include it in your go.mod file, e.g.:

```bash
go get github.com/avisi-cloud/go-client
```

## Getting Started

- [Get Started with Avisi Cloud](https://docs.avisi.cloud/docs/get-started/introduction/)
- [Platform Overview](https://docs.avisi.cloud/product/introduction/)
- [References](https://docs.avisi.cloud/references/references-overview/)
- [Example](example)

## Usage

### Setting up a client

To create a client that uses a personal access token (PAT) for authentication we first need to create a NewPersonalAccessTokenAuthenticator to which you pass in the PAT, e.g.:

```go
	authenticator := acloudapi.NewPersonalAccessTokenAuthenticator(token)
```

Then we create a `ClientOpts` object in which you configure the URL of the Avisi Cloud Platform API you want to connect to, e.g.:

```go
	clientOpts := acloudapi.ClientOpts{
		APIUrl: "https://example.com",
	}
```

Next, create a `NewCLient`, using the `Authenticator` and `ClientOpts` created in the previous steps:

```go
	c := acloudapi.NewClient(authenticator, clientOpts)
```

### Example

Full example:

```go
	personalAccessToken := os.Getenv("ACLOUD_PERSONAL_ACCESS_TOKEN")
	authenticator := acloudapi.NewPersonalAccessTokenAuthenticator(personalAccessToken)
	clientOpts := acloudapi.ClientOpts{
		APIUrl: "https://example.com",
	}
	c := acloudapi.NewClient(authenticator, clientOpts)

	createEnvironment := acloudapi.CreateEnvironment{
		Name:        "name",
		Type:        "staging",
		Description: "a description of the environment",
	}

	org := "organisation-slug"
	environment, err := client.CreateEnvironment(ctx, createEnvironment, org)

	if err != nil {
		return err
	}
	// environment has been created
```

## License

[Apache 2.0 License 2.0](lICENSE)
