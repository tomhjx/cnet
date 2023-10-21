package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInputOption_NewInputer(t *testing.T) {
	type fields struct {
		Path         string
		PollInterval time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    Inputer
		wantErr error
	}{
		{
			name:    "LocalFileInput",
			fields:  fields{Path: "/path/to/my/config.json"},
			want:    &LocalFileInput{},
			wantErr: nil,
		},
		{
			name:    "RemoteFileInput",
			fields:  fields{Path: "http://path.to/my/config.json"},
			want:    &RemoteFileInput{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := InputOption{
				Path:         tt.fields.Path,
				PollInterval: tt.fields.PollInterval,
			}
			got, err := o.NewInputer()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, reflect.TypeOf(tt.want), reflect.TypeOf(got))
		})
	}
}
