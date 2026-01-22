package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DefaultRegion to use for S3 credential creation
const defaultRegion = "us-east-1"

func buildConfig(cfg *Config) (aws.Config, error) {
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(defaultRegion),
	}
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		loadOpts = append(loadOpts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")))
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), loadOpts...)
	if err != nil {
		return aws.Config{}, err
	}

	return awsCfg, nil
}

func buildEndpointResolver(hostname string) aws.EndpointResolverWithOptions {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL: hostname,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
}

// SessionNew - Read the config for default credentials, if not provided use environment based variables
func SessionNew(config *Config) *s3.Client {
	cfg, err := buildConfig(config)
	if err != nil {
		panic(err) // or handle error properly
	}

	if config.HostBase != "" && config.HostBase != "s3.amazon.com" {
		fixedHost := config.HostBase
		if !strings.HasPrefix(config.HostBase, "http") {
			fixedHost = "https://" + config.HostBase
		}
		cfg.EndpointResolverWithOptions = buildEndpointResolver(fixedHost)
	}

	return s3.NewFromConfig(cfg)
}

// SessionForBucket - For a given S3 bucket, create an appropriate session that references the region
// that this bucket is located in
func SessionForBucket(config *Config, bucket string) (*s3.Client, error) {
	cfg, err := buildConfig(config)
	if err != nil {
		return nil, err
	}

	if config.HostBucket == "" || config.HostBucket == "%(bucket)s.s3.amazonaws.com" {
		svc := SessionNew(config)

		loc, err := svc.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{Bucket: &bucket})
		if err != nil {
			return nil, err
		}
		if loc.LocationConstraint == "" {
			// Use default service
			return svc, nil
		} else {
			cfg.Region = string(loc.LocationConstraint)
		}
	} else {
		host := strings.ReplaceAll(config.HostBucket, "%(bucket)s", bucket)
		fixedHost := host
		if !strings.HasPrefix(host, "http") {
			fixedHost = "https://" + host
		}
		cfg.EndpointResolverWithOptions = buildEndpointResolver(fixedHost)
	}

	return s3.NewFromConfig(cfg), nil
}
