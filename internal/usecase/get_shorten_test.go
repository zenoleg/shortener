package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/usecase/mocks"
)

func TestGetShortenUseCase_Do(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.AssertNumberOfCalls(t, "Store", 0)
		idGenerator.AssertNumberOfCalls(t, "Generate", 0)

		uc := NewGetShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(context.TODO(), NewGetShortURLQuery(false, "localhost", ""))

		assert.Empty(t, destination)
		assert.Error(t, err)
	})

	t.Run("When passed URL value is valid but not found, then return error", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.
			On("GetOriginalURL", context.TODO(), mock.Anything).
			Return(domain.URL(""), assert.AnError).
			Once()

		idGenerator.
			On("Generate", mock.Anything).
			Return(domain.ID("t92YuUGbw1WY4V2LvoDc0RHa")).
			Once()

		uc := NewGetShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(context.TODO(), NewGetShortURLQuery(false, "localhost", "https://example.com"))

		assert.Empty(t, destination)
		assert.Error(t, err)
	})

	t.Run("When passed URL value is valid and found, then return destination URL", func(t *testing.T) {
		storage := mocks.NewReadOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.
			On("GetOriginalURL", context.TODO(), mock.Anything).
			Return(domain.URL("https://example.com"), nil).
			Once()

		idGenerator.
			On("Generate", mock.Anything).
			Return(domain.ID("t92YuUGbw1WY4V2LvoDc0RHa")).
			Once()

		uc := NewGetShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(context.TODO(), NewGetShortURLQuery(true, "localhost", "https://example.com"))

		assert.NoError(t, err)
		assert.Equal(t, "https://localhost/link/t92YuUGbw1WY4V2LvoDc0RHa", destination.String())
	})
}
