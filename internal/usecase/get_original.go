package usecase

import (
	"strings"

	"github.com/zenoleg/shortener/internal/domain"
)

type GetOriginalUseCase struct {
	storage ReadOnlyStorage
}

func NewGetOriginalUseCase(storage ReadOnlyStorage) GetOriginalUseCase {
	return GetOriginalUseCase{storage: storage}
}

func (uc GetOriginalUseCase) Do(shortURL domain.URL) (domain.URL, error) {
	url := shortURL.String()
	split := strings.Split(url, "/")

	id, err := domain.NewID(split[len(split)-1])
	if err != nil {
		return "", err
	}

	original, err := uc.storage.GetOriginalURL(id)
	if err != nil {
		return "", err
	}

	return original, nil
}
