# Fall: A Go Web Framework

Fall is a lightweight and flexible web framework for Go, designed to simplify the development of web applications and APIs. It provides a set of tools and features to help you build robust and scalable applications quickly.

## Features

*   **Routing:** A powerful and flexible router that supports middleware and route grouping.
*   **Dependency Injection:** A simple and effective dependency injection container to manage your application's components.
*   **Middleware:** Support for middleware to add functionality to your request pipeline.
*   **WebSocket Support:** Built-in support for WebSocket connections.
*   **GORM Integration:** Seamless integration with the GORM library for database operations.

## Installation

To install Fall, use `go get`:

```bash
go get github.com/vfcunha/fall
```

## Getting Started

Here's a simple "Hello, World!" example to get you started with Fall:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/vfcunha/fall"
)

func main() {
	// Create a new Fall application
	app, err := fall.NewApp(fall.Development, &fall.DefaultEnvConfig{})
	if err != nil {
		panic(err)
	}

	// Get the router
	router := app.GetRouter()

	// Define a simple route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Start the server
	app.ListenAndServe("8080")
}
```

## Dependency Injection

Fall's dependency injection container allows you to manage your application's components with ease. You can register your services and dependencies, and Fall will automatically inject them where needed.

### Registering a Dependency

To register a dependency, use the `fall.Register` function:

```go
fall.Register("myService", func() (any, error) {
	return &MyService{}, nil
})
```

### Injecting a Dependency

You can inject a dependency into your structs using the `fall` tag:

```go
type MyController struct {
	MyService *MyService `fall:"myService"`
}
```

## Controllers

Controllers are responsible for handling requests and returning responses. To create a controller, you need to implement the `fall.Controller` interface:

```go
type MyController struct {
	// ... dependencies
}

func (c *MyController) Configure(router *fall.Router) {
	router.Get("/my-route", c.MyHandler)
}

func (c *MyController) MyHandler(w http.ResponseWriter, r *http.Request) {
	// ... handler logic
}
```

Fall will automatically discover and register your controllers.

## Routing

The router allows you to define routes and associate them with handlers. You can also use middleware to add functionality to your routes.

### Defining a Route

To define a route, use the methods on the `fall.Router` struct:

```go
router.Get("/users", getUsersHandler)
router.Post("/users", createUserHandler)
```

### Route Groups

You can group routes that share a common path prefix or middleware. This helps in organizing your routes and avoiding repetition.

To create a group, use the `Group` method on the `fall.Router`:

```go
// Create a group for API v1 routes
router.Group("/api/v1", func(apiV1 *fall.Router) {
    // Add routes to the group
    apiV1.Get("/users", getUsersHandler) // Path: /api/v1/users
    apiV1.Post("/products", createProductHandler) // Path: /api/v1/products
})
```

### Middleware

You can add middleware to the entire application, to specific routes, or to route groups.

#### Application-level Middleware

To apply middleware to all routes, pass it to the `fall.NewApp` function:

```go
// This middleware will be applied to every request
app, err := fall.NewApp(fall.Development, &fall.DefaultEnvConfig{}, myMiddleware)
```

#### Route-level Middleware

To apply middleware to a specific route, pass it as an additional argument to the route definition:

```go
// This middleware will only be applied to the /protected route
router.Get("/protected", myProtectedHandler, authMiddleware)
```

#### Group-level Middleware

To apply middleware to all routes within a group, use the `Use` method on the group:

```go
// Create a group for protected routes
protected := router.Group("/protected")

// Apply an authentication middleware to the entire group
protected.Use(authMiddleware)

// Add routes to the group
protected.Get("/profile", getProfileHandler) // Path: /protected/profile
protected.Post("/settings", updateSettingsHandler) // Path: /protected/settings
```
