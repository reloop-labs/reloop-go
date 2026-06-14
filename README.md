# Reloop Go SDK

## Before you send

You need two things:

1. **API key** — create one in your Reloop account
2. **Verified domain** — add and verify a sending domain; use it in the `from` address

For setup details and the full API reference, see [reloop.sh/docs](https://reloop.sh/docs).

## Send email

```bash
go get github.com/reloop-labs/reloop-go
```

```go
import (
    "fmt"

    reloop "github.com/reloop-labs/reloop-go"
)

client, err := reloop.NewClient(reloop.ClientOptions{
    APIKey: "rl_your_api_key_here",
})
if err != nil {
    panic(err)
}

result, err := client.Mail.Send(reloop.SendMailParams{
    From:    "Reloop <hello@your-verified-domain.com>",
    To:      "user@example.com",
    Subject: "Welcome to Reloop",
    HTML:    reloop.String("<p>Thanks for signing up.</p>"),
    Text:    reloop.String("Thanks for signing up."),
})
if err != nil {
    panic(err)
}

fmt.Println(result.MessageID, result.ID)
```

More examples and optional fields: [reloop.sh/docs](https://reloop.sh/docs)

## License

Licensed under the [Apache License 2.0](./LICENSE) with additional use restrictions from Reloop Labs (same as the [Reloop project](https://github.com/reloop-labs/reloop/blob/main/LICENSE)).
