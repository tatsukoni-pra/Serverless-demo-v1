package event

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type SnsRecordCollection struct {
	Record []SnsMessage `json:"Records"`
}

type SnsMessage struct {
	EventName string `json:"eventName"`
	S3Message SnsS3Message `json:"s3"`
}

type SnsS3Message struct {
	Bucket SnsS3BucketMessage `json:"bucket"`
	Object SnsS3ObjectMessage `json:"object"`
}

type SnsS3BucketMessage struct {
	Name string `json:"name"`
}

type SnsS3ObjectMessage struct {
	Key string `json:"key"`
}

type S3Info struct {
	Bucket string
	Key  string
}

func GetS3TrigerInfo(SNSEvent events.SNSEvent) *S3Info {
	for _, record := range SNSEvent.Records {
		var snsRecordCollection SnsRecordCollection
		if err := json.Unmarshal([]byte(record.SNS.Message), &snsRecordCollection); err != nil {
			return getS3TrigerErrorInfo()
		}
		for _, snsRecord := range snsRecordCollection.Record {
			return &S3Info {
				Bucket: snsRecord.S3Message.Bucket.Name,
				Key: snsRecord.S3Message.Object.Key}
		}
		return getS3TrigerErrorInfo()
	}
	return getS3TrigerErrorInfo()
}

func getS3TrigerErrorInfo() *S3Info {
	return &S3Info {
		Bucket: "no-info",
		Key: "no-info"}
}
