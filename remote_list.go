package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func remotePager(_ *Config, svc *s3.Client, uri string, delim bool, pager func(page *s3.ListObjectsV2Output)) error {
	u, err := FileURINew(uri)
	if err != nil || u.Scheme != "s3" {
		return fmt.Errorf("requires buckets to be prefixed with s3://")
	}

	maxKeys := int32(1000)
	params := &s3.ListObjectsV2Input{
		Bucket:  &u.Bucket,
		MaxKeys: &maxKeys,
	}
	if u.Path != "" && u.Path != "/" {
		params.Prefix = u.Key()
	}
	if delim {
		delimiter := "/"
		params.Delimiter = &delimiter
	}

	paginator := s3.NewListObjectsV2Paginator(svc, params)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return err
		}
		pager(page)
	}

	return nil
}

func remoteList(config *Config, svc *s3.Client, args []string) ([]FileObject, error) {
	result := make([]FileObject, 0)

	for _, arg := range args {
		pager := func(page *s3.ListObjectsV2Output) {
			for _, obj := range page.Contents {
				result = append(result, FileObject{
					Name:     *obj.Key,
					Size:     *obj.Size,
					Checksum: *obj.ETag,
				})
			}
		}

		remotePager(config, svc, arg, false, pager)
	}

	return result, nil
}
