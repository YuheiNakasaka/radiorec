package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/YuheiNakasaka/radiorec/config"
	"github.com/YuheiNakasaka/radiorec/internal/uploader"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AwsS3 is struct to upload to aws s3, which fill uploader interface.
type AwsS3 struct {
	uploader.Uploader
}

// Upload : upload to s3
func (s *AwsS3) Upload(path string, filename string) error {
	fmt.Println("Uploading...")

	// Read config file
	myconf := config.Config{}
	err := myconf.Init()
	if err != nil {
		return fmt.Errorf("Failed to load config %v", err)
	}

	accessKeyID := fmt.Sprintf("%v", myconf.List.Get("aws.s3.access_key_id"))
	secretAccessKey := fmt.Sprintf("%v", myconf.List.Get("aws.s3.secret_access_key"))
	region := fmt.Sprintf("%v", myconf.List.Get("aws.s3.region"))
	bucketName := fmt.Sprintf("%v", myconf.List.Get("aws.s3.bucket"))

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	cli := s3.New(session.New(), &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:      aws.String(region),
	})

	resp, err := cli.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("video/mp4"),
		Body:        file,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload: %v", err)
	}

	log.Println(awsutil.StringValue(resp))

	return err
}
