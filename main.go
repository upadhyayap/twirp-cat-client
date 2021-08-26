package main

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/upadhyayap/twirp-cat/twirp/service"
)

func main() {
	client := pb.NewHelloWorldProtobufClient("http://localhost:8080", &http.Client{})

	resp, err := client.Hello(context.Background(), &pb.HelloReq{Subject: "World"})
	if err == nil {
		fmt.Println(resp.Text) // prints "Hello World"
	}
}