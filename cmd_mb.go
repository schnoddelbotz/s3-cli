package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"
)

func MakeBucket(config *Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	u, err := FileURINew(args[0])
	if err != nil || u.Scheme != "s3" {
		return fmt.Errorf("ls requires buckets to be prefixed with s3://")
	}
	if u.Path != "/" {
		return fmt.Errorf("Parameter problem: Expecting S3 URI with just the bucket name set instead of '%s'", args[0])
	}

	params := &s3.CreateBucketInput{
		Bucket: &u.Bucket,
	}
	if _, err := svc.CreateBucket(context.TODO(), params); err != nil {
		return err
	}

	fmt.Printf("Bucket 's3://%s/' created\n", u.Bucket)
	return nil
}
