package usecase

import (
	"context"

	"github.com/zenoleg/shortener/internal/domain"
)

type (
	GetShortURLQuery struct {
		host        string
		originalURL string
		isSSL       bool
	}

	GetShortUseCase struct {
		storage     readOnlyStorage
		idGenerator idGenerator
	}
)

func NewGetShortURLQuery(isSSL bool, host string, originalURL string) GetShortURLQuery {
	return GetShortURLQuery{
		isSSL:       isSSL,
		host:        host,
		originalURL: originalURL,
	}
}

func NewGetShortenUseCase(storage readOnlyStorage, idGenerator idGenerator) GetShortUseCase {
	return GetShortUseCase{
		storage:     storage,
		idGenerator: idGenerator,
	}
}

func (uc GetShortUseCase) Do(ctx context.Context, query GetShortURLQuery) (DestinationURL, error) {
	url, err := domain.NewURL(query.originalURL)
	if err != nil {
		return "", err
	}

	id := uc.idGenerator.Generate(url)

	_, err = uc.storage.GetOriginalURL(ctx, id)
	if err != nil {
		return "", err
	}

	return NewDestinationURL(query.isSSL, query.host, id.String()), nil
}
