package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zenoleg/shortener/internal/domain"
	"github.com/zenoleg/shortener/internal/usecase/mocks"
)

func TestShortenUseCase_Do(t *testing.T) {
	t.Parallel()

	t.Run("When passed URL value is invalid, then return an error", func(t *testing.T) {
		storage := mocks.NewWriteOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.AssertNumberOfCalls(t, "Store", 0)
		idGenerator.AssertNumberOfCalls(t, "Generate", 0)

		uc := NewShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(NewShortenQuery(false, "localhost", ""))

		assert.Empty(t, destination)
		assert.Error(t, err)
	})

	t.Run("When passed URL value is valid, then save URL into a storage and return nil", func(t *testing.T) {
		storage := mocks.NewWriteOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.
			On("Store", mock.Anything).
			Return(nil).
			Once()

		idGenerator.
			On("Generate", mock.Anything).
			Return(domain.ID("t92YuUGbw1WY4V2LvoDc0RHa")).
			Once()

		uc := NewShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(NewShortenQuery(false, "localhost", "https://example.com"))

		assert.NoError(t, err)
		assert.Equal(t, "http://localhost/link/t92YuUGbw1WY4V2LvoDc0RHa", destination.String())
	})

	t.Run("When Storage return an error, then return it with empty destination URL", func(t *testing.T) {
		storage := mocks.NewWriteOnlyStorage(t)
		idGenerator := mocks.NewIDGenerator(t)

		storage.
			On("Store", mock.Anything).
			Return(errors.New("error")).
			Once()

		idGenerator.
			On("Generate", mock.Anything).
			Return(domain.ID("t92YuUGbw1WY4V2LvoDc0RHa")).
			Once()

		uc := NewShortenUseCase(storage, idGenerator)

		destination, err := uc.Do(NewShortenQuery(false, "localhost", "https://example.com"))

		assert.Error(t, err)
		assert.Empty(t, destination)
	})
}
