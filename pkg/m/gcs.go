package m

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCSBucket struct {
	Name   string
	Region string
}

type GCSModel struct {
	client           *storage.Client
	bucket           *storage.BucketHandle
	availableBuckets []*storage.BucketAttrs
	prefix           string
}

type GCSManager interface {
	Bucket() *storage.BucketHandle
	SetBucket(bucket string) error
	AvailableBuckets() []*storage.BucketAttrs
	Prefix() string
	setPrefix(prefix string) error
	MoveUp() error
	MoveDown(prefix string) error
	List() (prefixes []string, keys []string, err error)
	ListObjects(key string) ([]*storage.ObjectAttrs, error)
	Download(object *storage.ObjectAttrs, destPath string) (n int64, err error)
}

func NewGCSManager(project string) GCSManager {
	gcsm := GCSModel{}
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		panic(err)
	}
	gcsm.client = client
	if err := gcsm.setAvailableBuckets(project); err != nil {
		panic(err)
	}
	return &gcsm
}

func (gcsm *GCSModel) Bucket() *storage.BucketHandle {
	return gcsm.bucket
}

func (gcsm *GCSModel) SetBucket(bucket string) error {
	// if gcsm.bucket != "" {
	// 	return fmt.Errorf("bucket is already set: %s", gcsm.bucket)
	// }

	for _, ab := range gcsm.AvailableBuckets() {
		if ab.Name != bucket {
			continue
		}

		// found
		gcsm.bucket = gcsm.client.Bucket(ab.Name)
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
	query := &storage.Query{
		Prefix:    gcsm.prefix,
		Delimiter: "/",
	}

	it := gcsm.bucket.Objects(context.TODO(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return prefixes, keys, err
		}

		if attrs.Name == "" {
			prefix := strings.Replace(attrs.Prefix, gcsm.prefix, "", 1)
			prefixes = append(prefixes, prefix)
			continue
		}
		name := strings.Replace(attrs.Name, gcsm.prefix, "", 1)
		keys = append(keys, name)

	}
	return prefixes, keys, err
}

func (gcsm *GCSModel) ListObjects(key string) ([]*storage.ObjectAttrs, error) {
	query := &storage.Query{
		Prefix: gcsm.prefix,
	}

	var objs []*storage.ObjectAttrs
	it := gcsm.bucket.Objects(context.TODO(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objs, err
		}
		objs = append(objs, attrs)

	}
	return objs, nil
}

func (gcsm *GCSModel) Download(objectAttrs *storage.ObjectAttrs, destPath string) (n int64, err error) {
	object := gcsm.bucket.Object(objectAttrs.Name)
	if err = os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
		return 0, err
	}

	_, err = os.Stat(destPath)
	if err == nil {
		return 0, fmt.Errorf("exists")
	}

	f, err := os.Create(destPath)
	if err != nil {
		return 0, err
	}
	// #gosec G307
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	rc, err := object.NewReader(context.TODO())
	if err != nil {
		return 0, fmt.Errorf("Object(%v).NewReader: %v", object, err)
	}
	defer rc.Close()
	return io.Copy(f, rc)
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
