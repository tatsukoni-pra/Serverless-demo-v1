package main

import (
	"encoding/json"
    "bytes"
    "log"
    "net/http"
    "os"
	"io/ioutil"

	event "thumbnail/event"
	thumbnailExec "thumbnail/thumbnailExec"
	"github.com/aws/aws-lambda-go/events"
)

var(
	runtimeApiEndpointPrefix string
)

func init() {
	// https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/runtimes-api.html
	runtimeApiEndpointPrefix = "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
}

func main() {
    log.Println("handler started")
	// イベントループ
	for {
        func() {
            // コンテキスト情報を取得
            resp, _ := http.Get(runtimeApiEndpointPrefix + "next")
            defer func() {
                resp.Body.Close()
            }()

            // リクエストIDはヘッダに含まれる
            rId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
            log.Printf("実行中のリクエストID" + rId)

			// 処理本体
			data, _ := ioutil.ReadAll(resp.Body)
			_, err := handle(data)
			if err != nil {
				http.Post(respErrorEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
				log.Fatal(err)
			}

            // 最終的に真のランタイムに返すコンテンツはInvocation Response APIのリクエストボディに含める。
            http.Post(respEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
        }()
    }
}

func handle(payload []byte) (string, error) {
	// https://github.com/aws/aws-lambda-go/blob/main/lambda/handler.go#L115 とかを参考に。
	var snsEvent events.SNSEvent
	if err := json.Unmarshal(payload, &snsEvent); err != nil {
		log.Fatal(err)
	}
	// ここからlambda処理の本体
	S3TrigerInfo := event.GetS3TrigerInfo(snsEvent)
	thumbnailExec.ExecThumbnail(S3TrigerInfo.Bucket, S3TrigerInfo.Key)
	return "", nil
}

func respEndpoint(requestId string) string {
    return runtimeApiEndpointPrefix + requestId + "/response"
}

func respErrorEndpoint(requestId string) string {
    return runtimeApiEndpointPrefix + requestId + "/error"
}
