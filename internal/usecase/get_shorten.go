package usecase

import "github.com/zenoleg/shortener/internal/domain"

type (
	GetShortURLQuery struct {
		isSSL       bool
		host        string
		originalURL string
	}

	GetShortUseCase struct {
		storage     ReadOnlyStorage
		idGenerator IDGenerator
	}
)

func NewGetShortURLQuery(isSSL bool, host string, originalURL string) GetShortURLQuery {
	return GetShortURLQuery{
		isSSL:       isSSL,
		host:        host,
		originalURL: originalURL,
	}
}

func NewGetShortenUseCase(storage ReadOnlyStorage, idGenerator IDGenerator) GetShortUseCase {
	return GetShortUseCase{
		storage:     storage,
		idGenerator: idGenerator,
	}
}

func (uc GetShortUseCase) Do(query GetShortURLQuery) (DestinationURL, error) {
	url, err := domain.NewURL(query.originalURL)
	if err != nil {
		return "", err
	}

	id := uc.idGenerator.Generate(url)

	_, err = uc.storage.GetOriginalURL(id)
	if err != nil {
		return "", err
	}

	return NewDestinationURL(query.isSSL, query.host, id.String()), nil
}
