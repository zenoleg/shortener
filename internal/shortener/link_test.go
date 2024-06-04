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

func TestShortLink(t *testing.T) {
	t.Parallel()

	t.Run("When passed short link does not ends with /link/{shortID}, then return an error", func(t *testing.T) {
		cases := []string{
			"https://example.com",
			"https://example.com/link",
			"https://example.com/link/",
			"https://example.com/link/shortID/word",
		}

		for _, testCase := range cases {
			lnk, err := newShortLink(testCase)

			assert.Error(t, err)
			assert.Equal(t, shortLink{}, lnk)
		}
	})

	t.Run("When passed short link ends with /link/{shortID}, then return short id", func(t *testing.T) {
		cases := []string{
			"https://example.com/link/123",
			"https://example.com/link/123/",
		}

		for _, testCase := range cases {
			lnk, err := newShortLink(testCase)

			assert.NoError(t, err)
			assert.Equal(t, shortID{encoded: "123"}, lnk.shortID())
		}
	})
}
