package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type alwaysErrorStorage struct{}

func (s alwaysErrorStorage) Store(l link) error {
	return assert.AnError
}

func TestShortenUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		uc := NewShortenUseCase(NewInMemoryStorage())
		err := uc.Handle(" ")

		assert.Error(t, err)
	})

	t.Run("When passed URL value is valid, then save URL into a storage and return nil", func(t *testing.T) {
		storage := NewInMemoryStorage()

		uc := NewShortenUseCase(storage)
		err := uc.Handle("http://example.com")

		assert.NoError(t, err)
		assert.Len(t, storage.links, 1)
	})

	t.Run("When storage return an error, then use case must return it", func(t *testing.T) {
		storage := alwaysErrorStorage{}

		uc := NewShortenUseCase(storage)
		err := uc.Handle("http://example.com")

		assert.Error(t, err)
	})
}

func TestGenerateShortenUseCase_Handle(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		uc := NewGenerateShortenUseCase()
		short, err := uc.Handle("localhost", " ")

		assert.Error(t, err)
		assert.Equal(t, "", short.String())
	})

	t.Run("When passed URL value is valid, then return shorten url", func(t *testing.T) {
		uc := NewGenerateShortenUseCase()
		short, err := uc.Handle("localhost", "https://google.com")

		assert.NoError(t, err)
		assert.Equal(t, "localhost/link/t92YuUGbn92bn9yL6MHc0RHa", short.String())
	})
}
