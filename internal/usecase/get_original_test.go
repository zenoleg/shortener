package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/usecase/mocks"
)

func TestGetOriginalUseCase_Do(t *testing.T) {
	t.Parallel()

	t.Run("When no originals found by ID, then return an error", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)

		storage.
			On("GetOriginalURL", context.TODO(), domain.ID("123")).
			Return(domain.URL(""), assert.AnError).
			Once()

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Do(context.TODO(), "http://localhost/link/123")

		assert.Error(t, err)
		assert.Empty(t, original.String())
	})

	t.Run("When original found by short ID, then return it", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)

		storage.
			On("GetOriginalURL", context.TODO(), domain.ID("123")).
			Return(domain.URL("https://example.com"), nil).
			Once()

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Do(context.TODO(), "http://localhost/link/123")

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", original.String())
	})

	t.Run("When short link has empty ID, then return error", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)
		storage.AssertNumberOfCalls(t, "GetOriginalURL", 0)

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Do(context.TODO(), "http://localhost/link/")

		assert.Error(t, err)
		assert.Empty(t, original.String())
	})
}
