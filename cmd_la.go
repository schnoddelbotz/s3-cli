package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"
)

func ListAll(config *Config, c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) != 0 {
		return fmt.Errorf("la shouldn't have arguments")
	}

	svc := SessionNew(config)

	resp, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	for _, bucket := range resp.Buckets {
		uri := fmt.Sprintf("s3://%s", *bucket.Name)

		// Shared with "ls"
		listBucket(config, svc, []string{uri})
	}

	return nil
}
