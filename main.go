package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket      = flag.String("bucket", "", "S3 bucket to upload to")
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
	config := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(*awsKey, *awsSecret, nil))
	service := s3.New(config)
}
