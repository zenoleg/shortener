package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type alwaysErrorFakeStorage struct{}

func (s alwaysErrorFakeStorage) Store(l link) error {
	return assert.AnError
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
