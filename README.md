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

keys, err := client.ApiKeys.List(&reloop.ApiKeyListParams{
    Page:    reloop.Int(1),
    Limit:   reloop.Int(10),
    Enabled: reloop.Bool(true),
})

_, err = client.ApiKeys.Rotate("key_123456789")
_, err = client.ApiKeys.Pause("key_123456789")
_, err = client.ApiKeys.Enable("key_123456789")
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

## Domains

```go
domain, err := client.Domain.Create(reloop.CreateDomainParams{
    Domain:           "send.example.com",
    CustomReturnPath: reloop.String("inbound"),
    ClickTracking:    reloop.Bool(true),
    OpenTracking:     reloop.Bool(true),
    TLS:              reloop.DomainTLS(reloop.DomainTLSOpportunistic),
    SendingEmail:     reloop.Bool(true),
    ReceivingEmail:   reloop.Bool(true),
})

domains, err := client.Domain.List(&reloop.ListDomainsParams{
    Page:  reloop.Int(1),
    Limit: reloop.Int(10),
})

one, err := client.Domain.Get("domain_123456789")

_, err = client.Domain.Update("domain_123456789", reloop.UpdateDomainParams{
    ClickTracking: reloop.Bool(false),
})

status, err := client.Domain.Verify("domain_123456789")

_, err = client.Domain.ForwardDNS("domain_123456789", reloop.ForwardDNSParams{
    Email: "admin@example.com",
})

nameservers, err := client.Domain.GetNameservers("domain_123456789")

_, err = client.Domain.Delete("domain_123456789")
```
