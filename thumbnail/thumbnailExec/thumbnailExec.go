package thumbnailExec

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ExecThumbnail(bucketName string, objectKey string) {
	log.Println(fmt.Sprintf("画像リサイズ開始。対象オブジェクト: %s", objectKey))
	sess := session.Must(session.NewSession())

	// S3から元画像をダウンロード
	s3svc := s3.New(sess)
	s3Object, err := s3svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Fatal(err)
	}
	s3ObjectBody := s3Object.Body
	defer s3ObjectBody.Close()

	// TODO: 画像リサイズ

	// 元画像を別バケットにアップロード
	uploader := s3manager.NewUploader(sess)
	uploadKey := strings.Replace(objectKey, "tmp/", "upload/", 1)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("tatsukoni-lambda-demo-upload"),
		Key:    aws.String(uploadKey),
		Body:   s3ObjectBody,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("画像リサイズ完了。実施オブジェクト: %s", uploadKey))
}
