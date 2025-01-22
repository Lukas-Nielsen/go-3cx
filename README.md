# go-3cx

package to wrap the api of [3cx®](https://www.3cx.com/)

- first steps: https://www.3cx.com/docs/configuration-rest-api/
- openapi documentation: https://downloads-global.3cx.com/downloads/misc/restapi/3cxconfigurationapi.yaml

## used packages

- [Resty](https://github.com/go-resty/resty)

## getting started

```sh
go get github.com/Lukas-Nielsen/go-3cx
```

```go
import threecx "github.com/Lukas-Nielsen/go-3cx"
```

## usage

### conf

you need the FQDN of the instanz, username and password (with permissions for the action)

```go
client, err := new threecx.NewClient(threecx.Host{FQDN string, Port int, Debug bool});

client.SetUser(threecx.User{Username string, Password string, MFA string});
// or
client.SetRest(Rest{ClientID string, ClientSecret string});
// or
client.SetToken(Token{TokenType string, Expires int, AccessToken string, RefreshToken string});
```

### functions

#### get

```go
err := c.Get(<uri>, <result struct>, <query params as map[string]string>)
```

#### delete

```go
err := c.Delete(<uri>, <query params as map[string]string>)
```

#### post, put, patch

```go
err := c.<Post|Put|Patch>(<uri>, <payload>, <query params as map[string]string>)
```
