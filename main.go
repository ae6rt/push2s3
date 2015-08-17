package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

var (

	// common to s3 and dynamodb
	awsKey    = flag.String("access-key", "", "AWS Access-Key (required)")
	awsSecret = flag.String("access-secret", "", "AWS Access-Secret (required)")
	opCode    = flag.String("op-code", "", "s3 or dynamodb for whether this invocation pushes to S3 or adds an item to DynamodDb.  Command line options will vary depending. (required)")
	buildID   = flag.String("build-id", "", "S3 key pointing to bucket object.  (required)")

	// push to s3 params
	bucket = flag.String("bucket", "", "S3 bucket to upload to.  (required for s3)")
	file   = flag.String("file", "", "File to put to the S3 bucket. (required for s3)")

	// add-item to DynamoDb params
	buildStatus   = flag.Bool("build-status", true, "build passed (true) or failed (false).   (required for dynamodb)")
	projectKey    = flag.String("project-key", "", "project key, e.g., plat/users.   (required for dynamodb)")
	buildTime     = flag.Int("build-time", 0, "build time in seconds since the epoch.   (required for dynamodb)")
	buildDuration = flag.Int("build-duration", 0, "build duration in seconds.   (required for dynamodb)")

	versionFlag = flag.Bool("version", false, "Print version info and exit.")
	buildInfo   string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	}
}

func main() {
	if *opCode != "s3" && *opCode != "dynamodb" {
		log.Fatalf("op-code of s3 or dynamodb required\n")
	}
	if *awsKey == "" {
		log.Fatalf("Nonempty access-key value required\n")
	}
	if *awsSecret == "" {
		log.Fatalf("Nonempty access-secret value required\n")
	}
	if *buildID == "" {
		log.Fatalf("Nonempty build-id value required\n")
	}

	config := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(*awsKey, *awsSecret, "")).WithRegion("us-west-1").WithMaxRetries(3)

	switch *opCode {
	case "s3":
		// Reference: http://docs.aws.amazon.com/sdk-for-go/api/service/s3/S3.html#PutObject-instance_method

		if *bucket == "" {
			log.Fatalf("Nonempty bucket value required\n")
		}
		if *file == "" {
			log.Fatalf("Nonempty file value required\n")
		}

		data, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatalf("Error reading file %s: %v\n", *file, err)
		}

		params := &s3.PutObjectInput{
			Key:    aws.String(*buildID),
			Bucket: aws.String(*bucket),
			Body:   bytes.NewReader(data),
		}

		svc := s3.New(config)

		resp, err := svc.PutObject(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Generic AWS error with Code, Message, and original error (if any)
				log.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					// A service error occurred
					log.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
				}
			} else {
				// This case should never be hit, the SDK should always return an
				// error which satisfies the awserr.Error interface.
				log.Println(err.Error())
			}
		}
		log.Println(awsutil.Prettify(resp))

	case "dynamodb":
		// Reference: http://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/DynamoDB.html#PutItem-instance_method

		if *projectKey == "" {
			log.Fatalf("project-key must be non empty\n")
		}

		svc := dynamodb.New(config)
		params := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"buildID": {
					S: aws.String(*buildID),
				},
				"buildTime": {
					N: aws.String(fmt.Sprintf("%d", *buildTime)),
				},
				"projectKey": {
					S: aws.String(strings.ToLower(*projectKey)),
				},
				"buildElapsedTime": {
					N: aws.String(fmt.Sprintf("%d", *buildDuration)),
				},
				"buildStatus": {
					BOOL: aws.Bool(*buildStatus),
				},
			},
			TableName: aws.String("inf-eng-build-server-2"),
		}
		resp, err := svc.PutItem(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Generic AWS error with Code, Message, and original error (if any)
				log.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					// A service error occurred
					log.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
				}
			} else {
				// This case should never be hit, the SDK should always return an
				// error which satisfies the awserr.Error interface.
				log.Println(err.Error())
			}
		}
		log.Println(awsutil.Prettify(resp))
	}
}
