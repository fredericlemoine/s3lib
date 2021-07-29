package s3lib

// This package provides simple functions to interact with an s3 bucket
// More information on the usage of aws-sdk-go:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Lib struct {
	accessKey      string
	secretKey      string
	endpoint       string
	region         string
	forcePathStyle bool
}

type S3File struct {
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass string
}

func NewS3Cli(accessKey, secretKey, endpoint, region string, forcePathStyle bool) *S3Lib {
	return &S3Lib{
		accessKey:      accessKey,
		secretKey:      secretKey,
		endpoint:       endpoint,
		region:         region,
		forcePathStyle: forcePathStyle,
	}
}

// Provider Interface
func (s *S3Lib) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     s.accessKey,
		SecretAccessKey: s.secretKey,
	}, nil
}

// Provider Interface
func (s *S3Lib) IsExpired() bool {
	return false
}

func (s *S3Lib) awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(s.endpoint),
		Credentials:      credentials.NewCredentials(s),
		Region:           aws.String(s.region),
		S3ForcePathStyle: &s.forcePathStyle,
	}))
}

// UploadReader uploads the content of the Reader to the s3 bucket
func (s *S3Lib) UploadReader(r io.Reader, bucket, path string) (err error) {
	var result *s3manager.UploadOutput

	// The session the S3 Uploader will use
	sess := s.awsSession()

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   r,
	})

	if err == nil {
		log.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	}
	return
}

// Upload uploads a file to the s3 bucket
func (s *S3Lib) Upload(filename string, bucket, path string) (err error) {
	var f *os.File

	if f, err = os.Open(filename); err != nil {
		err = fmt.Errorf("failed to open file %q, %v", filename, err)
		return
	}
	defer f.Close()
	if err = s.UploadReader(f, bucket, path); err != nil {
		return
	}

	return
}

// Download downloads an item from an s3 bucket and write it to the given writer
func (s *S3Lib) DownloadWriter(w io.WriterAt, bucket, path string) (err error) {
	// The session the S3 Downloader will use
	sess := s.awsSession()

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Write the contents of S3 Object to the Writer
	n, err := downloader.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err == nil {
		fmt.Printf("file downloaded, %d bytes\n", n)
	}

	return
}

// Download downloads an item from an s3 bucket and write it to the given file
func (s *S3Lib) Download(filename string, bucket, path string) (err error) {
	var f *os.File

	// Create a file to write the S3 Object contents to.
	if f, err = os.Create(filename); err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}
	defer f.Close()
	err = s.DownloadWriter(f, bucket, path)

	return
}

// List lists items in an s3 bucket
// given by a prefix. An empty prefix will list all the items
func (s *S3Lib) List(bucket, path string) (files []S3File, err error) {
	var resp *s3.ListObjectsV2Output
	files = make([]S3File, 0)

	// The session the S3 Downloader will use
	sess := s.awsSession()

	// Create a s3 client
	sc := s3.New(sess)

	resp, err = sc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return
	}
	for _, c := range resp.Contents {
		files = append(files, S3File{
			Key:          *c.Key,
			LastModified: *c.LastModified,
			Size:         *c.Size,
			StorageClass: *c.StorageClass,
		})
	}
	return
}

// Delete deletes an item from an s3 bucket
// if the object does not exists, returns no error
func (s *S3Lib) Delete(bucket, path string) (err error) {
	sess := s.awsSession()

	sc := s3.New(sess)

	_, err = sc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		err = fmt.Errorf("Unable to delete object %q from bucket %q, %v", path, bucket, err)
		return
	}

	err = sc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	if err == nil {
		log.Printf("Object %q successfully deleted\n", path)
	}

	return
}
