package thumbnailExec

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
)

const (
	uploadBacketName = "tatsukoni-lambda-demo-upload"
	readPrefix = "tmp/"
	uploadPrefix = "upload/"
)

var (
	resizedWidth uint = 256
	resizedHeight uint = 0
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

	// 画像リサイズ
	img, data, err := image.Decode(s3ObjectBody)
	if err != nil {
		log.Fatal(err)
	}
	resizedImg := resize.Resize(resizedWidth, resizedHeight, img, resize.NearestNeighbor)
	// リサイズした画像をエンコード
	buf := new(bytes.Buffer)
	switch data {
		case "png":
			if err := png.Encode(buf, resizedImg); err != nil {
				log.Fatal(err)
			}
		case "jpeg", "jpg":
			opts := &jpeg.Options{Quality: 100}
			if err := jpeg.Encode(buf, resizedImg, opts); err != nil {
				log.Fatal(err)
			}
		default:
			if err := png.Encode(buf, resizedImg); err != nil {
				log.Fatal(err)
			}
	}

	// 元画像を別バケットにアップロード
	uploader := s3manager.NewUploader(sess)
	uploadKey := strings.Replace(objectKey, readPrefix, uploadPrefix, 1)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(uploadBacketName),
		Key:    aws.String(uploadKey),
		Body:   buf,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("画像リサイズ完了。実施オブジェクト: %s", uploadKey))
}
