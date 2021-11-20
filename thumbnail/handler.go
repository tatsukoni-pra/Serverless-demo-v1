package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest() (string, error) {
	fmt.Println("hello world!")
	fmt.Println("output test v2")
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
