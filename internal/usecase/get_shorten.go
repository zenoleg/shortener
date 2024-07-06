package usecase

import "github.com/zenoleg/shortener/internal/domain"

type (
	GetShortenQuery struct {
		isSSL       bool
		host        string
		originalURL string
	}

	GetShortenUseCase struct {
		storage     ReadOnlyStorage
		idGenerator IDGenerator
	}
)

func NewGetShortenQuery(isSSL bool, host string, originalURL string) GetShortenQuery {
	return GetShortenQuery{
		isSSL:       isSSL,
		host:        host,
		originalURL: originalURL,
	}
}

func NewGetShortenUseCase(storage ReadOnlyStorage, idGenerator IDGenerator) GetShortenUseCase {
	return GetShortenUseCase{
		storage:     storage,
		idGenerator: idGenerator,
	}
}

func (uc GetShortenUseCase) Do(query ShortenQuery) (DestinationURL, error) {
	url, err := domain.NewURL(query.originalURL)
	if err != nil {
		return "", err
	}

	id := uc.idGenerator.Generate(url)

	return NewDestinationURL(query.isSSL, query.host, id.String()), nil
}
