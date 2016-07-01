# badge
Minimal, highly performant authentication tokens. Like a JWT without JSON and arbitrary payloads.

## purpose

Badges are strictly for authentication, not carrying around data. For this reason, badges only consist of a username, id, and signature. They are meant to be sent with every request (as a cookie, header, etc.) and checked on the server to determine if the user is authenticated.

## restrictions

- Usernames must be >= 1 && <= 255 bytes.
- Ids must be uint32.

## examples

### creating a badge

```go
var ExampleBadge []byte
ExampleBadge, _ = badge.New([]byte("karl"), uint32(1), []byte("secret"))
// string(ExampleBadge) == "04karl01000000Usco0AX0HdHJnBrdjdqfi2uyH-mO0KrSpkLiQNJ3BCw"
```

### getting badge values

```go
username, id, auth := badge.Get(ExampleBadge, []byte("secret"))
// username == "karl"
// id == 1
// auth == true

username, id, auth = badge.Get(ExampleBadge, []byte("wrong secret"))
// username == ""
// id == 0
// auth == false
```
