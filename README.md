# Reloop Go SDK

Official Go client for the Reloop API.

## Install

```bash
go get github.com/reloop-labs/reloop-go
```

## Usage

```go
package main

import "github.com/reloop-labs/reloop-go"

func main() {
    client, err := reloop.NewClient(reloop.ClientOptions{
        APIKey: "re_123456789",
    })
    if err != nil {
        panic(err)
    }

    contact, err := client.Contacts.Create(map[string]interface{}{
        "email":      "john.doe@example.com",
        "first_name": "John",
        "last_name":  "Doe",
        "unsubscribed": false,
    })
    if err != nil {
        panic(err)
    }

    _ = contact
}
```

## API Keys

```go
apiKey, err := client.ApiKeys.Create(reloop.CreateApiKeyParams{
    Name:             "Production key",
    Enabled:          reloop.Bool(true),
    RateLimitEnabled: reloop.Bool(true),
})
```

## Contacts

```go
contacts, err := client.Contacts.List(map[string]interface{}{
    "page":  1,
    "limit": 10,
})

group, err := client.Contacts.CreateGroup(map[string]interface{}{
    "name": "Beta Testers",
})

_, err = client.Contacts.Groups.AddContact("grp_123456789", map[string]interface{}{
    "contact_id": "cont_123456789",
})
```
