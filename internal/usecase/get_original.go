package usecase

import (
	"context"
	"strings"

	"github.com/zenoleg/shortener/internal/domain"
)

type GetOriginalUseCase struct {
	storage readOnlyStorage
}

func NewGetOriginalUseCase(storage readOnlyStorage) GetOriginalUseCase {
	return GetOriginalUseCase{storage: storage}
}

func (uc GetOriginalUseCase) Do(ctx context.Context, shortURL domain.URL) (domain.URL, error) {
	url := shortURL.String()
	split := strings.Split(url, "/")

	id, err := domain.NewID(split[len(split)-1])
	if err != nil {
		return "", err
	}

	original, err := uc.storage.GetOriginalURL(ctx, id)
	if err != nil {
		return "", err
	}

	return original, nil
}
