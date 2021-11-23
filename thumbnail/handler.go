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
    nextEndpoint string
)

func init() {
	runtimeApiEndpointPrefix = "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
    nextEndpoint = runtimeApiEndpointPrefix + "next"
}

func main() {
    log.Println("handler started")

	for {
        func() {
            // コンテキスト情報を取得
            resp, _ := http.Get(nextEndpoint)
            defer func() {
                resp.Body.Close()
            }()

            // ヘッダにはリクエストIDが含まれています。
            rId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
            log.Printf("実行中のリクエストID" + rId)

			// 処理本体
			data, _ := ioutil.ReadAll(resp.Body)
			handle(data)

            // 最終的に真のランタイムに返すコンテンツはInvocation Response APIのリクエストボディに含めます。
            http.Post(respEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
        }()
    }
}

func handle(payload []byte) (string, error) {
	var snsEvent events.SNSEvent
	if err := json.Unmarshal(payload, &snsEvent); err != nil {
		log.Fatal(err)
	}
	S3TrigerInfo := event.GetS3TrigerInfo(snsEvent)
	thumbnailExec.ExecThumbnail(S3TrigerInfo.Bucket, S3TrigerInfo.Key)
	return "", nil
}

func respEndpoint(requestId string) string {
    return runtimeApiEndpointPrefix + requestId + "/response"
}
