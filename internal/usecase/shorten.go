package usecase

import "github.com/zenoleg/shortener/internal/domain"

type (
	ShortenQuery struct {
		isSSL       bool
		host        string
		originalURL string
	}

	ShortenUseCase struct {
		storage     WriteOnlyStorage
		idGenerator IDGenerator
	}

	WriteOnlyStorage interface {
		Store(shortenURL domain.ShortenURL) error
	}

	IDGenerator interface {
		Generate(originalURL domain.URL) domain.ID
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

func (uc ShortenUseCase) Do(query ShortenQuery) (string, error) {
	url, err := domain.NewURL(query.originalURL)
	if err != nil {
		return "", err
	}

	id := uc.idGenerator.Generate(url)

	shortenURL := domain.NewShortenURL(id, url)

	err = uc.storage.Store(shortenURL)
	if err != nil {
		return "", err
	}

	return shortenURL.ID(), nil
}
