package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type alwaysErrorFakeStorage struct{}

func (s alwaysErrorFakeStorage) Store(l link) error {
	return assert.AnError
}

func TestShortenUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		uc := NewShortenUseCase(NewInMemoryStorage(map[string]string{}))
		destination, err := uc.Handle(NewShortenQuery(false, "localhost", " "))

		assert.Error(t, err)
		assert.Equal(t, "", destination.String())
	})

	t.Run("When passed URL value is valid, then save URL into a storage and return nil", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{})

		uc := NewShortenUseCase(storage)
		destination, err := uc.Handle(NewShortenQuery(false, "localhost", "http://example.com"))

		assert.NoError(t, err)
		assert.Len(t, storage.links, 1)
		assert.Equal(t, "http://localhost/link/t92YuUGbw1WY4V2LvoDc0RHa", destination.String())
	})

	t.Run("When storage return an error, then use case must return it", func(t *testing.T) {
		storage := alwaysErrorFakeStorage{}

		uc := NewShortenUseCase(storage)
		destination, err := uc.Handle(NewShortenQuery(false, "localhost", "http://example.com"))

		assert.Error(t, err)
		assert.Equal(t, "", destination.String())
	})
}

func TestGetShortUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{})
		uc := NewGenerateShortenUseCase(storage)
		short, err := uc.Handle(NewShortenQuery(false, "localhost", " "))

		assert.Error(t, err)
		assert.Equal(t, "", short.String())
	})

	t.Run("When passed URL value is valid but not found, then return error", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{})

		uc := NewGenerateShortenUseCase(storage)
		short, err := uc.Handle(NewShortenQuery(false, "localhost", "https://google.com"))

		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, "", short.String())
	})

	t.Run("When passed URL value is valid and found in storage, then return shorten url", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{
			"t92YuUGbn92bn9yL6MHc0RHa": "https://google.com",
		})

		uc := NewGenerateShortenUseCase(storage)
		short, err := uc.Handle(NewShortenQuery(false, "localhost", "https://google.com"))

		assert.NoError(t, err)
		assert.Equal(t, "http://localhost/link/t92YuUGbn92bn9yL6MHc0RHa", short.String())
	})
}

func TestGetOriginalUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When no originals found by short ID, then return an error", func(t *testing.T) {
		uc := NewGetOriginalUseCase(NewInMemoryStorage(map[string]string{}))
		original, err := uc.Handle("https://example.com/link/123")

		assert.Error(t, err)
		assert.Equal(t, "", original)
	})

	t.Run("When short link is invalid, then return error", func(t *testing.T) {
		uc := NewGetOriginalUseCase(NewInMemoryStorage(map[string]string{}))
		original, err := uc.Handle("https://example.com/123")

		assert.Error(t, err)
		assert.Equal(t, "", original)
	})

	t.Run("When original found by short ID, then return it", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{})
		lnk, _ := newLink("https://example.com")
		_ = storage.Store(lnk)

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Handle("https://example.com/link/" + lnk.ShortID().String())

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", original)
	})
}

func TestGetOriginalByIDUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When no originals found by short ID, then return an error", func(t *testing.T) {
		uc := NewGetOriginalForRedirectUseCase(NewInMemoryStorage(map[string]string{}))
		original, err := uc.Handle("123")

		assert.Error(t, err)
		assert.Equal(t, "", original)
	})

	t.Run("When original found by short ID, then return it", func(t *testing.T) {
		storage := NewInMemoryStorage(map[string]string{})
		lnk, _ := newLink("https://example.com")
		_ = storage.Store(lnk)

		uc := NewGetOriginalForRedirectUseCase(storage)
		original, err := uc.Handle(lnk.shortID.String())

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", original)
	})
}
