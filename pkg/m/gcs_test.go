package m

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/storage"
)

func TestNewGCSManager(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{projectID: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewGCSManager(tt.args.projectID)
		})
	}
}

func TestGCSModel_ListObjects(t *testing.T) {
	client, _ := storage.NewClient(context.TODO())
	bucket := client.Bucket("k-jun-bucket")
	type fields struct {
		client           *storage.Client
		bucket           *storage.BucketHandle
		availableBuckets []*storage.BucketAttrs
		prefix           string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "debug",
			args:   args{key: ""},
			fields: fields{bucket: bucket},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcsm := &GCSModel{
				client:           tt.fields.client,
				bucket:           tt.fields.bucket,
				availableBuckets: tt.fields.availableBuckets,
				prefix:           tt.fields.prefix,
			}
			got, err := gcsm.ListObjects(tt.args.key)
			fmt.Println(got)
			fmt.Println(len(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("GCSModel.ListObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GCSModel.ListObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
