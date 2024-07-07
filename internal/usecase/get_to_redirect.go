package usecase

import "github.com/zenoleg/shortener/internal/domain"

type GetOriginalForRedirectUseCase struct {
	storage ReadOnlyStorage
}

func NewGetOriginalForRedirectUseCase(storage ReadOnlyStorage) GetOriginalForRedirectUseCase {
	return GetOriginalForRedirectUseCase{storage: storage}
}

func (uc GetOriginalForRedirectUseCase) Do(id domain.ID) (domain.URL, error) {
	return uc.storage.GetOriginalURL(id)
}
