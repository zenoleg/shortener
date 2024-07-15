package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{
			name:    "when id is empty then return error",
			args:    args{id: "  "},
			want:    ID(""),
			wantErr: true,
		},
		{
			name:    "when id is not empty then return id",
			args:    args{id: " id "},
			want:    ID("id"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewID(tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewURL(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL string
	}
	tests := []struct {
		name    string
		args    args
		want    URL
		wantErr bool
	}{
		{
			name:    "when original url is invalid then return error",
			args:    args{originalURL: "invalid-url"},
			want:    URL(""),
			wantErr: true,
		},
		{
			name:    "when original url is valid then return url",
			args:    args{originalURL: "http://example.com"},
			want:    URL("http://example.com"),
			wantErr: false,
		},
		{
			name:    "when original url is empty then return error",
			args:    args{originalURL: "  "},
			want:    URL(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewURL(tt.args.originalURL)

			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
