package m

import (
	"testing"
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
