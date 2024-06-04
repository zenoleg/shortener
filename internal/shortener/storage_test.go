package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("When called, then return a new in-memory storage", func(t *testing.T) {
		storage := NewInMemoryStorage()

		assert.NotNil(t, storage.links)
		assert.Len(t, storage.links, 0)
	})
}

func TestInMemoryStorage_Store(t *testing.T) {
	t.Parallel()

	t.Run("When called with a link, then store it in the storage", func(t *testing.T) {
		storage := NewInMemoryStorage()
		lnk, _ := newLink("https://google.com")

		err := storage.Store(lnk)

		assert.NoError(t, err)
		assert.Len(t, storage.links, 1)
	})
}

func TestInMemoryStorage_GetOriginalURL(t *testing.T) {
	t.Parallel()

	t.Run("When called with a shortID URL, then return the original URL", func(t *testing.T) {
		storage := NewInMemoryStorage()
		lnk, _ := newLink("https://google.com")
		_ = storage.Store(lnk)

		original, err := storage.GetOriginalURL(lnk.shortID)

		assert.NoError(t, err)
		assert.Equal(t, lnk.Original(), original)
	})

	t.Run("When called with an unknown shortID URL, then return an error", func(t *testing.T) {
		storage := NewInMemoryStorage()
		lnk, _ := newLink("https://google.com")

		_, err := storage.GetOriginalURL(lnk.shortID)

		assert.Error(t, err)
	})
}
