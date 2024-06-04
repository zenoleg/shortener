package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestinationURL_String(t *testing.T) {
	t.Parallel()

	t.Run("When SSL is true, then return link with https protocol", func(t *testing.T) {
		url := newDestinationURL(true, "example.com", shortID{encoded: "short"})

		assert.Equal(t, "https://example.com/link/short", url.String())
	})

	t.Run("When SSL is false, then return link with http protocol", func(t *testing.T) {
		url := newDestinationURL(false, "example.com", shortID{encoded: "short"})

		assert.Equal(t, "http://example.com/link/short", url.String())
	})
}
