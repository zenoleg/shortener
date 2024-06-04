package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOriginalURL(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is empty, then return an error", func(t *testing.T) {
		url, err := newOriginalURL(" ")

		assert.Error(t, err)
		assert.Equal(t, originalURL{}, url)
	})

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		url, err := newOriginalURL("invalid-url")

		assert.Error(t, err)
		assert.Equal(t, originalURL{}, url)
	})

	t.Run("When passed URL contains spaces, then trim them", func(t *testing.T) {
		url, err := newOriginalURL(" https://google.com ")

		assert.NoError(t, err)
		assert.Equal(t, "https://google.com", url.String())
	})

	t.Run("When passed URL value is valid, then return an error", func(t *testing.T) {
		original := "https://google.com/?key=value"
		url, err := newOriginalURL(original)

		assert.NoError(t, err)
		assert.Equal(t, original, url.String())
	})
}

func TestShortURL(t *testing.T) {
	t.Run("When passed original URL is valid, then return a shortID URL", func(t *testing.T) {
		original, err := newOriginalURL("https://google.com")
		short := newShortID(original)

		assert.NoError(t, err)
		assert.Equal(t, "t92YuUGbn92bn9yL6MHc0RHa", short.String())
	})
}

func TestLink(t *testing.T) {
	t.Parallel()

	t.Run("When passed original URL is invalid, then return an error", func(t *testing.T) {
		lnk, err := newLink("invalid-url")

		assert.Error(t, err)
		assert.Equal(t, link{}, lnk)
	})

	t.Run("When passed original URL is valid, then return a link", func(t *testing.T) {
		lnk, err := newLink("https://google.com")

		assert.NoError(t, err)
		assert.Equal(t, link{
			original: originalURL{original: "https://google.com"},
			shortID:  shortID{encoded: "t92YuUGbn92bn9yL6MHc0RHa"},
		}, lnk)
	})
}
