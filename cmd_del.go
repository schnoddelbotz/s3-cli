package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/urfave/cli/v2"
)

// TODO: Handle --recusrive
func DeleteObjects(config *Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNew(config)

	buckets := make(map[string][]types.ObjectIdentifier, 0)

	for _, path := range args {
		u, err := FileURINew(path)
		if err != nil || u.Scheme != "s3" {
			return fmt.Errorf("rm requires buckets to be prefixed with s3://")
		}

		if (u.Path == "" || strings.HasSuffix(u.Path, "/")) && !config.Recursive {
			return fmt.Errorf("Parameter problem: Expecting S3 URI with a filename or --recursive: %s", path)
		}

		objects := buckets[u.Bucket]
		if objects == nil {
			objects = make([]types.ObjectIdentifier, 0)
		}
		key := u.Key()
		buckets[u.Bucket] = append(objects, types.ObjectIdentifier{Key: key})
	}

	// FIXME: Limited to 1000 objects, that's that shouldn't be an issue, but ...
	for bucket, objects := range buckets {
		bsvc, err := SessionForBucket(config, bucket)
		if err != nil {
			return err
		}

		if config.Recursive {
			for _, obj := range objects {
				uri := fmt.Sprintf("s3://%s/%s", bucket, *obj.Key)

				remotePager(config, svc, uri, false, func(page *s3.ListObjectsV2Output) {
					olist := make([]types.ObjectIdentifier, 0)
					for _, item := range page.Contents {
						olist = append(olist, types.ObjectIdentifier{Key: item.Key})

						fmt.Printf("delete: s3://%s/%s\n", bucket, *item.Key)
					}

					if !config.DryRun {
						params := &s3.DeleteObjectsInput{
							Bucket: &bucket,
							Delete: &types.Delete{
								Objects: olist,
							},
						}

						_, err := bsvc.DeleteObjects(context.TODO(), params)
						if err != nil {
							fmt.Println("Error removing")
						}
					}
				})
			}
		} else if !config.DryRun {
			params := &s3.DeleteObjectsInput{
				Bucket: &bucket,
				Delete: &types.Delete{
					Objects: objects,
				},
			}

			_, err := bsvc.DeleteObjects(context.TODO(), params)
			if err != nil {
				return err
			}
		}
		for _, objs := range objects {
			fmt.Printf("delete: s3://%s/%s\n", bucket, *objs.Key)
		}
	}

	return nil
}
