package main

import (
	"net/http"

	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
)

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

func init() {
	ln := net.Listen()

	// Amazon API Gateway binary media types are supported out of the box.
	// If you don't send or receive binary data, you can safely set it to nil.
	Handle = apigatewayproxy.New(ln, nil).Handle

	// Any Go framework complying with the Go http.Handler interface can be used.
	// This includes, but is not limited to, Vanilla Go, Gin, Echo, Gorrila, Goa, etc.
	go http.Serve(ln, http.HandlerFunc(RSSMergeHandler))
}
