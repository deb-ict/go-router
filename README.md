# go-router

go-router is a high performance and secure http router for go

## Installation
`go get -u github.com/deb-ict/go-router`

## Routing
``` 
func apiHandler(w http.ResponseWriter, r* http.Request) {
    
}

router := NewRouter()
router.HandleFunc("/api/example", apiHandler)
```

## Middlewares
```
router := NewRouter()
router.Use(authenticationMiddleware, authorizationMiddleware)
```

## Authentication

### Basic authentication
```
m := authentication.NewAuthenticationMiddleware(authentication.NewBasicAuthenticationHandler(validator))
```

### Api-Key authentication
```
m := authentication.NewAuthenticationMiddleware(authentication.NewApiKeyAuthenticationHandler(validator))
```

#### Customization
```
m := authentication.NewAuthenticationMiddleware(authentication.NewApiKeyAuthenticationHandler(validator,
    authentication.WithApiKeyAuthenticationHeaderName("MY-API-HEADER"),
    authentication.WithApiKeyAuthenticationQueryParamName("security_key"),
))
```

### Bearer token authentication
```
m := authentication.NewAuthenticationMiddleware(authentication.NewBearerAuthenticationHandler(validator))
```

## Authorization

### Policy definitions
```
m := authentication.NewAuthorizationMiddleware(
    authentication.WithPolicy(
        "management",
        authentication.NewUserRequirement(),
        authentication.NewRoleRequirement("admin", "manager"),
    ),
    authentication.WithPolicy(
        "read",
        authentication.NewScopeRequirement("data.read"),
    ),
    authentication.WithPolicy(
        "management_and_read",
        authentication.NewRoleRequirement("admin", "manager"),
        authentication.NewScopeRequirement("data.read"),
    ),
    authentication.WithPolicy(
        "management_or_read",
        authentication.NewCombinedRequirement(false,
            authentication.NewRoleRequirement("admin", "manager"),
            authentication.NewScopeRequirement("data.read"),
        ),
    ),
)
```

In the following policy definition, a tree of combined requirements is defined.  
The "tree" policy requires the following combination of scopes:  
`((a or b) and c) or ((d or e) and f)`
```
m := authentication.NewAuthorizationMiddleware(
    authentication.WithPolicy("tree", authentication.NewCombinedRequirement(false,
        authentication.NewCombinedRequirement(true,
            authentication.NewScopeRequirement("a", "b"),
            authentication.NewScopeRequirement("c"),
        ),
        authentication.NewCombinedRequirement(true,
            authentication.NewScopeRequirement("d", "e"),
            authentication.NewScopeRequirement("f"),
        ),
    )),
)
```

## Full example

```
package main

import (
	"fmt"
	"net/http"

	"github.com/deb-ict/go-router"
	"github.com/deb-ict/go-router/authentication"
	"github.com/deb-ict/go-router/authorization"
)

type ApiKeyValidator struct {
}

func (v *ApiKeyValidator) GetApiKeyAuthenticationData(apiKey string) (authentication.ClaimMap, error) {
	if apiKey == "123" {
		claims := authentication.ClaimMap{}
		claims.AddClaim("id", "12345")
		claims.AddClaim("name", "John Doe")
		return claims, nil
	}
	return nil, fmt.Errorf("invalid API key")
}

func callHandler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func callHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Space!")
}

func main() {
	// Setup the routes
	r := router.NewRouter()
	r.HandleFunc("/world", callHandler1).AllowedMethod(http.MethodGet).Authorize("example")
	r.HandleFunc("/space", callHandler2).AllowedMethod(http.MethodGet)

	// Setup the authentication and authorization middleware
	authentication.UseMiddleware(r, authentication.NewApiKeyAuthenticationHandler(
		&ApiKeyValidator{},
		authentication.WithApiKeyAuthenticationHeaderName("X-API-Key"),
	))
	authorization.UseMiddleware(r, authorization.WithPolicy(
		"example",
		authorization.NewUserRequirement(),
	))

	http.ListenAndServe(":8080", r)
}
```