# Gomid Middleware Library

**Gomid** is a lightweight and flexible middleware library for Go. 
It provides a simple way to wrap functions with middleware, enabling pre- and post-processing of function calls. 
Whether you're building Lambda handlers or general-purpose applications, Gomid makes it easy to manage middleware logic.

## Features

- **Middleware Support**: Add `Before` and `After` hooks to your functions.
- **Dynamic Function Wrapping**: Automatically wraps functions with middleware using reflection.
- **Custom Middleware**: Create and integrate your own middleware implementations.
- **Error Handling**: Supports error propagation in middleware.

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
