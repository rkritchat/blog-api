package helper

import (
	"blog-api/internal/config"
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsS3 interface {
	Upload(filename string, b []byte) error
	Download(filename string) ([]byte, error)
}

type awsS3 struct {
	session *session.Session
	env     config.Env
}

func NewS3(session *session.Session, env config.Env) AwsS3 {
	return &awsS3{
		env:     env,
		session: session,
	}
}

func (h awsS3) Upload(filename string, b []byte) error {
	conn := s3manager.NewUploader(h.session)
	_, err := conn.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.env.S3BucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return err
	}

	return nil
}

func (h awsS3) Download(filename string) ([]byte, error) {
	conn := s3manager.NewDownloader(h.session)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := conn.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(h.env.S3BucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
