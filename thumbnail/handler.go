package main

import (
	"context"

	event "thumbnail/event"
	thumbnailExec "thumbnail/thumbnailExec"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, events events.SNSEvent) (string, error) {
	S3TrigerInfo := event.GetS3TrigerInfo(events)
	thumbnailExec.ExecThumbnail(S3TrigerInfo.Bucket, S3TrigerInfo.Key)
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
