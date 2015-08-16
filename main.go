package main

// Put a file to an S3 bucket
// Reference: http://docs.aws.amazon.com/sdk-for-go/api/service/s3/S3.html#PutObject-instance_method

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket      = flag.String("bucket", "", "S3 bucket to upload to")
	key         = flag.String("key", "", "S3 key pointing to bucket object")
	file        = flag.String("file", "", "File to put to the S3 bucket")
	awsKey      = flag.String("access-key", "", "AWS Access-Key")
	awsSecret   = flag.String("access-secret", "", "AWS Access-Secret")
	versionFlag = flag.Bool("version", false, "Print version info and exit.")

	buildInfo string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	}
}

func main() {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Error reading file %s: %v\n", *file, err)
	}

	config := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(*awsKey, *awsSecret, "")).WithRegion("us-west-1")

	params := &s3.PutObjectInput{
		Key:    aws.String(*key),
		Bucket: aws.String(*bucket),
		Body:   bytes.NewReader(data),
	}

	svc := s3.New(config)

	resp, err := svc.PutObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}
	log.Println(awsutil.Prettify(resp))
}
