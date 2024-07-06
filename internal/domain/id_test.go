package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62IDGenerator_Generate(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL URL
	}
	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{
			name:    "when original url is valid then return id",
			args:    args{originalURL: URL("https://example.com")},
			want:    ID("t92YuUGbw1WY4V2LvozcwRHdoB"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Base62IDGenerator{}
			got, err := g.Generate(tt.args.originalURL)

			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
