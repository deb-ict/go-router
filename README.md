# go-router

go-router is a high performance and secure http router for go

## Installation
`go get -u github.com/deb-ict/go-router`

## Authentication

### Basic authentication
```
m := NewAuthenticationMiddleware(NewBasicAuthenticationHandler(validator))
```

### Api-Key authentication
```
m := NewAuthenticationMiddleware(NewApiKeyAuthenticationHandler(validator))
```

#### Customization
```
m := NewAuthenticationMiddleware(NewApiKeyAuthenticationHandler(validator,
    WithApiKeyAuthenticationHeaderName("MY-API-HEADER"),
    WithApiKeyAuthenticationQueryParamName("security_key"),
))
```

### Bearer token authentication
```
m := NewAuthenticationMiddleware(NewBearerAuthenticationHandler(validator))
```

## Authorization

### Policy definitions
```
m := NewAuthorizationMiddleware(
    WithPolicy(
        "management",
        NewUserRequirement(),
        NewRoleRequirement("admin", "manager"),
    ),
    WithPolicy(
        "read",
        NewScopeRequirement("data.read"),
    ),
    WithPolicy(
        "management_and_read",
        NewRoleRequirement("admin", "manager"),
        NewScopeRequirement("data.read"),
    ),
    WithPolicy(
        "management_or_read",
        NewCombinedRequirement(false,
            NewRoleRequirement("admin", "manager"),
            NewScopeRequirement("data.read"),
        ),
    ),
)
```

In the following policy definition, a tree of combined requirements is defined.  
The "tree" policy requires the following combination of scopes:  
`((a or b) and c) or ((d or e) and f)`
```
m := NewAuthorizationMiddleware(
    WithPolicy("tree", NewCombinedRequirement(false,
        NewCombinedRequirement(true,
            NewScopeRequirement("a", "b"),
            NewScopeRequirement("c"),
        ),
        NewCombinedRequirement(true,
            NewScopeRequirement("d", "e"),
            NewScopeRequirement("f"),
        ),
    )),
)
```