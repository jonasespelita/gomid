# Gomid Middleware Library

**Gomid** is a lightweight and flexible middleware library for Go. 
It provides a simple way to wrap functions with middleware, enabling pre- and post-processing of function calls. 
Whether you're building Lambda handlers or general-purpose applications, Gomid makes it easy to manage middleware logic.

Gomid implements the classic onion-like middleware pattern (FILO).

## Features

- **Middleware Support**: Add `Before` and `After` hooks to your functions.
- **Dynamic Function Wrapping**: Automatically wraps functions with middleware using reflection.
- **Custom Middleware**: Create and integrate your own middleware implementations.
- **Error Handling**: Supports error propagation in middleware.

## Execution Order

Middlewares have two phases: before and after.

The before phase, happens before the handler is executed. In this code the response is not created yet, so you will have access only to the request.

The after phase, happens after the handler is executed. In this code you will have access to the response.

If you have three middlewares attached (as in the image above), this is the expected order of execution:

- `middleware1` (before)
- `middleware2` (before)
- `middleware3` (before)
- `handler`
- `middleware3` (after)
- `middleware2` (after)
- `middleware1` (after)

Notice that in the after phase, middlewares are executed in inverted order, this way the first handler attached is the one with the highest priority as it will be the first able to change the request and last able to modify the response before it gets sent to the user.

## Installation

To install Gomid, use `go get`:

```bash
go get github.com/jonasespelita/gomid
```

## Usage

### Basic Example

Wrap a function with middleware:

```go
func handle(ctx context.Context, request string) (string, error) {
	return "Hello, " + request, nil
}

type ExampleMiddleware struct{}

func (e ExampleMiddleware) Before() func(...any) {
	return func(args ...any) {
		println("Before middleware executed")
	}
}

func (e ExampleMiddleware) After() func(...any) {
	return func(results ...any) {
		println("After middleware executed")
	}
}

func main() {
	wrapped := gomid.New(handle, gomid.AddMid(ExampleMiddleware{}))
	lambda.Start(wrapped)
}
```

### AWS Lambda Example

Gomid is compatible with AWS Lambda handlers:

```go
func handle(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Body:       "Hello, " + request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	wrapped := gomid.New(handle)
	lambda.Start(wrapped)
}
```

### Custom Middlewares

Implement the `GomidWare` interface to create custom middleware:

```go
type CustomMiddleware struct{}

func (c CustomMiddleware) Before() func(...any) {
	return func(args ...any) {
		println("Custom Before logic")
	}
}

func (c CustomMiddleware) After() func(...any) {
	return func(results ...any) {
		println("Custom After logic")
	}
}
```

See `sloggonmidware.go` for sample middleware that logs function calls and their results.

## API Reference

### `gomid.New(handler any, options ...GomidOption) any`

Initializes function as a gomid function.

### `gomid.AddMid(middleware GomidWare) GomidOption`

Adds a middleware to the function.

### `GomidWare` Interface

- `Before() func(...any)`: Logic to execute before the function call.
- `After() func(...any)`: Logic to execute after the function call.

## Inspiration

The name **Gomid** is a playful nod to Gromit from the beloved Wallace and Gromit series, symbolizing the library's role as a helpful companion in your Go projects.

![gromit.png](gromit.png)
