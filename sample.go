package main

import (
	"context"
	"errors"
	"github.com/jonasespelita/gomid/gomid"
	"reflect"
)

// lambda handler sigs https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures

// lambda sample 1
func handle(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
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
	wallace := gomid.New(handle1,
		gomid.AddMid(HiByeGomidWare{}),
		gomid.AddMid(NewSlogGomidWare()),
	)

	test(wallace)
	//lambda.Start(handle)

}

// test function to invoke handler
func test(h any) {
	hVal := reflect.ValueOf(h)

	//hVal.Call([]reflect.Value{
	//	reflect.ValueOf(context.Background()),
	//	reflect.ValueOf(events.APIGatewayProxyRequest{Body: "hello worlds"}),
	//})

	hVal.Call([]reflect.Value{
		// no args
	})

	// check response
	//println("got response:", response[0].Interface().(*events.APIGatewayProxyResponse).Body)
	//println("got response", *(response[0].Interface().(*string)))
}

// sample middlewares (GomidWare)

type HiByeGomidWare struct{}

func (h HiByeGomidWare) Before() func(out ...any) {
	return func(out ...any) {
		println("hi")
	}
}

func (h HiByeGomidWare) After() func(out ...any) {
	return func(out ...any) {
		println("bye")
	}
}

//func (h HiByeGomidWare) After() func(out ...any) {
//	return func(out ...any) {
//		println("adding bye to response")
//		if len(out) > 0 {
//			// dynamically handle out type
//			switch rType := out[0].(type) {
//			case *string:
//				*rType += " bye bye"
//			case *events.APIGatewayProxyResponse:
//				rType.Body = rType.Body + " bye bye"
//			default:
//				println("unsupported response type")
//			}
//		} else {
//			println("bye")
//		}
//	}
//}

var _ gomid.GomidWare = (*HiByeGomidWare)(nil)
