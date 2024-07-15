package usecase

import (
	"context"

	"github.com/zenoleg/shortener/internal/domain"
)

type (
	ShortenQuery struct {
		host        string
		originalURL string
		isSSL       bool
	}

	ShortenUseCase struct {
		storage     WriteOnlyStorage
		idGenerator IDGenerator
	}
)

func NewShortenQuery(isSSL bool, host string, originalURL string) ShortenQuery {
	return ShortenQuery{
		isSSL:       isSSL,
		host:        host,
		originalURL: originalURL,
	}
}

func NewShortenUseCase(storage WriteOnlyStorage, idGenerator IDGenerator) ShortenUseCase {
	return ShortenUseCase{
		storage:     storage,
		idGenerator: idGenerator,
	}
}

func (uc ShortenUseCase) Do(ctx context.Context, query ShortenQuery) (DestinationURL, error) {
	url, err := domain.NewURL(query.originalURL)
	if err != nil {
		return "", err
	}

	id := uc.idGenerator.Generate(url)

	shortenURL := domain.NewShortenURL(id, url)

	err = uc.storage.Store(ctx, shortenURL)
	if err != nil {
		return "", err
	}

	return NewDestinationURL(query.isSSL, query.host, id.String()), nil
}
