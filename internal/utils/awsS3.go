package utils

import (
	"context"
	"fileTransfer/internal/repository"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsS3 struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	BucketName string
}

func NewAwsS3() *AwsS3 {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	return &AwsS3{
		Client:     s3Client,
		Uploader:   uploader,
		BucketName: os.Getenv("S3_BUCKET_NAME"),
	}
}

// UploadFile : Upload file to S3-AWS
func (a *AwsS3) UploadFile(file *multipart.FileHeader) (*manager.UploadOutput, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	key := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)

	res, err := a.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(a.BucketName),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		return nil, err
	}

	//log.Println("1: ", *res.Expiration, "1: ", *res.Key, "1: ", res.Location)
	return res, nil
}

// ListFiles Lists all files on S3-AWS
func (a *AwsS3) ListFiles() (*s3.ListObjectsV2Output, error) {
	resp, err := a.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(a.BucketName),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DownloadFile Downloads file from S3-AWS
func (a *AwsS3) DownloadFile(key string) (*s3.GetObjectOutput, error) {
	resp, err := a.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(a.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GenerateSignedURL Generates signed URL for uploaded files
func (a *AwsS3) GenerateSignedURL(key string, expiry time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.BucketName),
		Key:    aws.String(key),
	}

	presignClient := s3.NewPresignClient(a.Client)

	presignResult, err := presignClient.PresignGetObject(context.TODO(), input, func(po *s3.PresignOptions) {
		po.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

// DeleteExpiredFiles Check and Delete Expired file from AWS
func (a *AwsS3) DeleteExpiredFiles(repo repository.FileDbRepo) {
	log.Println("Checking for expired files...")

	expiredFiles, err := repo.GetExpiredFiles(time.Now())
	if err != nil {
		log.Printf("Failed to get expired files from DB: %v", err)
		return
	}

	for _, file := range expiredFiles {
		// Delete from S3
		_, err := a.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(a.BucketName),
			Key:    aws.String(file.S3Key),
		})
		if err != nil {
			log.Printf("Failed to delete from S3: %v", err)
			continue
		}

		// Delete DB record
		err = repo.DeleteFileByID(file.ID)
		if err != nil {
			log.Printf("Failed to delete DB record: %v", err)
		} else {
			log.Printf("Deleted file: %s", file.S3Key)
		}
	}
}
