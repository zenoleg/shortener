package usecase

import (
	"context"

	"github.com/zenoleg/shortener/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=writeOnlyStorage
type writeOnlyStorage interface {
	Store(ctx context.Context, shortenURL domain.ShortenURL) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=readOnlyStorage
type readOnlyStorage interface {
	GetOriginalURL(ctx context.Context, id domain.ID) (domain.URL, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=idGenerator
type idGenerator interface {
	Generate(originalURL domain.URL) domain.ID
}
