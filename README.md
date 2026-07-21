# Reloop Go SDK

Official Go client for the [Reloop](https://reloop.sh) API.

## Before you send

1. **API key** — create one in your Reloop account
2. **Verified domain** — add and verify a sending domain; use it in the `from` address

Full API reference: [reloop.sh/docs](https://reloop.sh/docs)

## Install

```bash
go get github.com/reloop-labs/reloop-go/v2
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	reloop "github.com/reloop-labs/reloop-go/v2"
)

func main() {
	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey: "rl_your_api_key_here",
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Mail.Send(reloop.SendMailParams{
		From:    "Reloop <hello@your-verified-domain.com>",
		To:      "user@example.com",
		Subject: "Welcome to Reloop",
		HTML:    reloop.String("<p>Thanks for signing up.</p>"),
		Text:    reloop.String("Thanks for signing up."),
	})
	if err != nil {
		if ve, ok := err.(*reloop.ValidationError); ok {
			log.Fatalf("invalid request (%s): %s", ve.Field, ve.Message)
		}
		if ae, ok := err.(*reloop.APIError); ok {
			log.Fatalf("API error %d: %s", ae.Status, ae.Error())
		}
		log.Fatal(err)
	}

	fmt.Println(result.MessageID, result.ID)
}
```

## Services

```go
client.APIKey    // create, list, get, update, delete, rotate, enable, disable
client.Mail      // Send
client.Domain    // create, list, get, update, delete, verify
client.Contacts  // CRUD + .Properties, .Groups, .Channels
client.Webhook   // CRUD, pause/enable/disable, trigger, deliveries + Verify
client.Inbox     // .Mailboxes, .Messages, .Threads
```

## Errors

| Kind | Type | When |
|------|------|------|
| Bad client args | `*ValidationError` | Invalid params — **no HTTP call** |
| HTTP / network | `*APIError` | Non-2xx response or transport failure |

Idiomatic Go `(T, error)` returns — the Go equivalent of Node/Python Result objects.

## Webhook verification

```go
event, err := reloop.ConstructEvent(payload, signature, secret, 300)
// or client.Webhook.Verify(reloop.VerifyWebhookParams{...})
```

## License

Licensed under the [Apache License 2.0](./LICENSE) with additional use restrictions from Reloop Labs (same as the [Reloop project](https://github.com/reloop-labs/reloop/blob/main/LICENSE)).
