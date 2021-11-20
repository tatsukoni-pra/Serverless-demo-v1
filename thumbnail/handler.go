package main

import (
	"context"
	"fmt"

	event "thumbnail/event"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, events events.SNSEvent) (string, error) {
	fmt.Println(event.GetS3TrigerInfo(events).Key)
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
