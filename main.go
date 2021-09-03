package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/twitchtv/twirp"
	pb "github.com/upadhyayap/twirp-cat/twirp/service"
)

func main() {
	protobufClient := pb.NewHelloWorldProtobufClient("http://localhost:8080", &http.Client{})
	// with custom path prefix
	//protobufClient := pb.NewHelloWorldProtobufClient("http://localhost:8080", &http.Client{}, twirp.WithClientPathPrefix("/my/custom/prefix"))
	jsonClient := pb.NewHelloWorldJSONClient("http://localhost:8080", &http.Client{})

	protoResp, err := protobufClient.Hello(context.Background(), &pb.HelloReq{Subject: "World from protobuf client"})
	if err == nil {
		fmt.Println(protoResp.Text) // prints "Hello World" via protobuf client
	} else {
		if twerr, ok := err.(twirp.Error); ok {
			fmt.Println(twerr.Code())
		}
	}

	// unwrapping internal erros, few twirp erros are also wrapped in as internal erros like tranport level erros. to get the
	// actual cause, you need to unwrap it.
	/* 	if err != nil {
		if twerr, ok := err.(twirp.Error); ok {
			if twerr.Code() == twirp.Internal {
				if transportErr := errors.Unwrap(twerr); transportErr != nil {
					// transportErr could be something like an HTTP connection error
				}
			}
		}
	} */

	jsonResp, err := jsonClient.Hello(context.Background(), &pb.HelloReq{Subject: "World from json client"})

	if err == nil {
		fmt.Println(jsonResp.Text) // prints "Hello world" via json client
	}

	// Setting up custom HTTP headers
	header := make(http.Header)
	header.Set("bearerToken", "some token")
	ctx, err := twirp.WithHTTPRequestHeaders(context.Background(), header)
	// Now you will have access to http header through context like ctx.value("customHeader") in the server, https://twitchtv.github.io/twirp/docs/headers.html
	if err != nil {
		log.Println("Error setting custom header")
	}

	res, err := protobufClient.Hello(ctx, &pb.HelloReq{Subject: "Twirp"})
	fmt.Println(res.Text)
}

func NewLoggingClientHooks() *twirp.ClientHooks {
	return &twirp.ClientHooks{
		RequestPrepared: func(ctx context.Context, r *http.Request) (context.Context, error) {
			fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) {
			log.Println("Error: " + string(twerr.Code()))
		},
		ResponseReceived: func(ctx context.Context) {
			log.Println("Success")
		},
	}
}
