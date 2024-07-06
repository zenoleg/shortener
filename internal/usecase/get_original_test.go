package usecase

import (
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/usecase/mocks"
)

func TestGetOriginalUseCase_Do(t *testing.T) {
	t.Parallel()

	t.Run("When no originals found by ID, then return an error", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)

		storage.
			On("GetOriginalURL", domain.ID("123")).
			Return(domain.URL(""), errors.New("nothing")).
			Once()

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Do("http://localhost/link/123")

		assert.Error(t, err)
		assert.Empty(t, original.String())
	})

	t.Run("When original found by short ID, then return it", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)

		storage.
			On("GetOriginalURL", domain.ID("123")).
			Return(domain.URL("https://example.com"), nil).
			Once()

		uc := NewGetOriginalUseCase(storage)
		original, err := uc.Do("http://localhost/link/123")

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", original.String())
	})
}
