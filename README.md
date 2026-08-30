# jsonx

A JSON facade that routes each workload to the stronger backend: `jsoniter`
for typed marshaling and Go 1.27 `encoding/json/v2` for dynamic marshaling and
all unmarshaling.

## Installation

```sh
go get github.com/jellybeanci/jsonx
```

## Usage

### JSON Serialization and Deserialization

```go
import "github.com/jellybeanci/jsonx"

// Struct to be serialized
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{Name: "John Doe", Email: "john@example.com"}

// Marshal to JSON
jsonData, err := jsonx.Marshal(user)
if err != nil {
    log.Fatal(err)
}

// Unmarshal from JSON
var parsedUser User
if err := jsonx.Unmarshal(jsonData, &parsedUser); err != nil {
    log.Fatal(err)
}
```

### Casting Helpers

```go
// Cast JSON response directly into a struct
rawData, err := fetchData()
user, err := jsonx.Cast[User](rawData, err)
if err != nil {
    log.Fatal(err)
}

// Cast JSON response into a slice of structs
rawList, err := fetchList()
users, err := jsonx.CastSlice[User](rawList, err)
if err != nil {
	log.Fatal(err)
}
```

## Features

- **Typed marshaling** powered by jsoniter's standard-library-compatible config
- **Dynamic marshaling** powered by `encoding/json/v2`
- **Unmarshaling** powered by `encoding/json/v2`
- **Convenient casting functions** for structured data
- **Minimal performance overhead** with efficient memory usage

## Backend routing

| Operation | Top-level value | Backend |
|---|---|---|
| Marshal | `map[string]any` or `[]any` | `encoding/json/v2` |
| Marshal | All other concrete types | jsoniter compatible config |
| Unmarshal | All destination types | `encoding/json/v2` |

Routing inspects only the top-level runtime type. An `any` variable containing
a `map[string]any` or `[]any` is detected correctly. A typed struct stored in
`any` remains a typed struct and uses jsoniter.

Nested `any` fields and named map/slice aliases are intentionally not scanned.
This keeps the typed hot path at a constant-cost type switch. JSON v2 semantics
apply to routed dynamic containers; for example, nil maps and slices marshal as
`{}` and `[]` instead of `null`.

