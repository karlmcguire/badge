# badge
Minimal, highly performant authentication tokens. Like a JWT without JSON and arbitrary payloads.

## purpose

Badges are strictly for authentication, not carrying around data. For this reason, badges only consist of a username, id, and signature. They are meant to be sent with every request (as a cookie, header, etc.) and checked on the server to determine if the user is authenticated.

## rules

- Usernames must be <= 12 bytes (for now).
- Ids must be uint32.

## examples

### creating a badge

```go
var ExampleBadge []byte
ExampleBadge, _ = badge.New([]byte("karl"), uint32(1), []byte("secret"))
// string(ExampleBadge) == "dXNlcm5hbWU9PT09.AQAAAA==.HjL9WjyH6hIKHWaR_pwujS7eHU0P2tQRuSIGFnmUEzE="
```

### checking a badge

```go
ok := badge.Check(ExampleBadge, []byte("secret"))
// ok == true

ok = badge.Check(ExampleBadge, []byte("wrong secret"))
// ok == false
```

### getting badge values

```go
username, id, ok := badge.Get(ExampleBadge, []byte("secret"))
// username == "karl"
// id == 1
// ok == true

username, id, ok = badge.Get(ExampleBadge, []byte("wrong secret"))
// username == ""
// id == 0
// ok == false
```
