package main

import (
	"context"
	"errors"
	//"github.com/aws/aws-lambda-go/events"
	"github.com/jonasespelita/gomid/gomid"
	"reflect"
)

// lambda handler sigs https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures

// lambda sample 1
func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	println("mxg from context " + ctx.Value("msg1").(string))
	return events.APIGatewayProxyResponse{
		Body:       "i'm an app with request body: " + request.Body,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// lambda sample 2
func handle1() (*string, error) {
	s := "hello there"
	return &s, errors.New("sample error")
}

// lambda sample 3
func handle2() {
	println("hello there")
}

func main() {
	wallace := gomid.New(handle,
		gomid.AddMid(NewSlogGomidWare()),
		gomid.AddMid(HiByeGomidWare{}),
	)

	test(wallace)
	//lambda.Start(handle)

}

// test function to invoke handler
func test(h any) {
	hVal := reflect.ValueOf(h)

	bckCtx := context.Background()
	hVal.Call([]reflect.Value{
		reflect.ValueOf(context.WithValue(bckCtx, "msg", "hello from test")),
		reflect.ValueOf(events.APIGatewayProxyRequest{Body: "hello worlds"}),
	})
	//
	//hVal.Call([]reflect.Value{
	//	// no args
	//})

	// check response
	//println("got response:", response[0].Interface().(*events.APIGatewayProxyResponse).Body)
	//println("got response", *(response[0].Interface().(*string)))
}

// sample middlewares (GomidWare)

type HiByeGomidWare struct{}

func (h HiByeGomidWare) Before() func(...any) []any {
	return func(in ...any) []any {
		if len(in) < 1 {
			println("hi no edit")
			return in
		}
		ctx, ok := in[0].(context.Context)
		if ok {
			in[0] = context.WithValue(ctx, "msg1", "hello from gomid")
		}
		return in
	}
}

//func (h HiByeGomidWare) After() func(out ...any) []any {
//	return func(out ...any) []any {
//		println("bye")
//		return out
//	}
//}

func (h HiByeGomidWare) After() func(...any) []any {
	return func(out ...any) []any {
		if len(out) > 0 {
			// dynamically handle out type
			switch rType := out[0].(type) {
			case *string:
				*rType += " bye bye"
			case events.APIGatewayProxyResponse:
				rType.Body = rType.Body + " bye bye"
				out[0] = rType
			default:
				println("unsupported response type")
			}
		} else {
			println("bye")
		}
		return out
	}
}

var _ gomid.GomidWare = (*HiByeGomidWare)(nil)
