package usecase

import "github.com/zenoleg/shortener/internal/domain"

type (
	WriteOnlyStorage interface {
		Store(shortenURL domain.ShortenURL) error
	}

	ReadOnlyStorage interface {
		GetOriginalURL(id domain.ID) (domain.URL, error)
	}

	IDGenerator interface {
		Generate(originalURL domain.URL) domain.ID
	}
)
