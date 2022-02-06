package m

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iterator"
)

type GCSBucket struct {
	Name   string
	Region string
}

type GCSModel struct {
	client           *storage.Client
	bucket           string
	availableBuckets []*storage.BucketAttrs
	prefix           string
}

type GCSManager interface {
	Bucket() string
	SetBucket(bucket string) error
	AvailableBuckets() []*storage.BucketAttrs
	Prefix() string
	setPrefix(prefix string) error
	MoveUp() error
	MoveDown(prefix string) error
	List() (prefixes []string, keys []string, err error)
	ListObjects(key string) []interface{}
	Download(object interface{}, destPath string) (n int64, err error)
}

func NewGCSManager(projectID string) GCSManager {
	gcsm := GCSModel{}
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		panic(err)
	}
	gcsm.client = client

	if projectID == "" {
		credentials := defaultCredentials(context.TODO(), compute.ComputeScope)
		projectID = credentials.ProjectID
	}
	gcsm.setAvailableBuckets(projectID)
	return &gcsm
}

func defaultCredentials(ctx context.Context, scope string) *google.Credentials {
	credentials, err := google.FindDefaultCredentials(ctx, scope)
	if err != nil {
		panic(err)
	}
	return credentials
}

func (gcsm *GCSModel) Bucket() string {
	return gcsm.bucket
}

func (gcsm *GCSModel) SetBucket(bucket string) error {
	if gcsm.bucket != "" {
		return fmt.Errorf("bucket is already set: %s", gcsm.bucket)
	}

	for _, ab := range gcsm.AvailableBuckets() {
		if ab.Name != bucket {
			continue
		}

		// found
		gcsm.bucket = bucket

		// TODO
		// opts := []optsFunc{
		// 	config.WithRegion(ab.Region),
		// }
		// if gcsm.endpointURL != "" {
		// 	endpoint := aws.EndpointResolverFunc(func(service, r string) (aws.Endpoint, error) {
		// 		return aws.Endpoint{
		// 			URL:               gcsm.endpointURL,
		// 			SigningRegion:     r,
		// 			HostnameImmutable: gcsm.pathStyle,
		// 		}, nil
		// 	})
		// 	opts = append(opts, config.WithEndpointResolver(endpoint))
		// }

		// // re-create client with region
		// cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
		// if err != nil {
		// 	panic(err)
		// }

		// gcsm.client = s3.NewFromConfig(cfg)
		// gcsm.downloader = gcsmanager.NewDownloader(gcsm.client)

		return nil
	}

	return fmt.Errorf("not found in available buckets: %s", bucket)

}

func (gcsm *GCSModel) AvailableBuckets() []*storage.BucketAttrs {
	return gcsm.availableBuckets
}

func (gcsm *GCSModel) setAvailableBuckets(projectID string) error {
	buckets := []*storage.BucketAttrs{}
	it := gcsm.client.Buckets(context.TODO(), projectID)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		buckets = append(buckets, attrs)
	}
	gcsm.availableBuckets = buckets
	return nil
}

func (gcsm *GCSModel) Prefix() string {
	return gcsm.prefix
}

func (gcsm *GCSModel) setPrefix(prefix string) error {
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		return fmt.Errorf("prefix must be end with slash: %s", prefix)
	}

	gcsm.prefix = prefix
	return nil
}

func (gcsm *GCSModel) MoveUp() error {
	return gcsm.setPrefix(upperPrefix((gcsm.prefix)))
}

func (gcsm *GCSModel) MoveDown(prefix string) error {
	return gcsm.setPrefix(gcsm.prefix + prefix)

}

func (gcsm *GCSModel) List() (prefixes []string, keys []string, err error) {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) ListObjects(key string) []interface{} {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) Download(object interface{}, destPath string) (n int64, err error) {
	panic("not implemented") // TODO: Implement
}

func upperPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}

	prefixNoslash := prefix[:len(prefix)-1]
	i := strings.LastIndex(prefixNoslash, "/")

	if i == -1 {
		// "foo/" => ""
		return ""
	}

	// "foo/bar/baz/" => "foo/bar/"
	return prefixNoslash[:i+1]
}

func lastPartPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}

	prefixNoslash := prefix[:len(prefix)-1]
	i := strings.LastIndex(prefixNoslash, "/")

	if i == -1 {
		// "foo/" => "foo/"
		return prefix
	}

	// "foo/bar/baz/" => "baz/"
	return prefix[i+1:]
}
