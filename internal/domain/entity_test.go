package domain

import "testing"

func TestNewID(t *testing.T) {
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
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewURL(t *testing.T) {
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
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
