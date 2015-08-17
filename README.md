Push file to an S3 bucket
-------------------------

```
Usage of ./s3dynamo-darwin-amd64:
  -access-key="": AWS Access-Key (required)
  -access-secret="": AWS Access-Secret (required)
  -bucket="": S3 bucket to upload to.  (required for s3)
  -build-duration=0: build duration in seconds.   (required for dynamodb)
  -build-id="": S3 key pointing to bucket object.  (required)
  -build-status=true: build passed (true) or failed (false).   (required for dynamodb)
  -build-time=0: build time in seconds since the epoch.   (required for dynamodb)
  -file="": File to put to the S3 bucket. (required for s3)
  -op-code="": s3 or dynamodb for whether this invocation pushes to S3 or adds an item to DynamodDb.  Command line options will vary depending. (required)
  -project-key="": project key, e.g., plat/users.   (required for dynamodb)
  -version=false: Print version info and exit.
```
