package usecase

import (
	"context"

	"github.com/zenoleg/shortener/internal/domain"
)

type GetOriginalForRedirectUseCase struct {
	storage ReadOnlyStorage
}

func NewGetOriginalForRedirectUseCase(storage ReadOnlyStorage) GetOriginalForRedirectUseCase {
	return GetOriginalForRedirectUseCase{storage: storage}
}

func (uc GetOriginalForRedirectUseCase) Do(ctx context.Context, id domain.ID) (domain.URL, error) {
	return uc.storage.GetOriginalURL(ctx, id)
}
