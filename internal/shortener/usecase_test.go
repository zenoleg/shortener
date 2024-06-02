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
