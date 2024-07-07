package storage

import (
	"sync"

	"emperror.dev/errors"
	"github.com/zenoleg/shortener/internal/domain"
)

var ErrURLNotFound = errors.New("Original URL not found")

type (
	InMemoryStorage struct {
		links map[string]string
		mu    sync.RWMutex
	}
)

func NewInMemoryStorage(store map[string]string) *InMemoryStorage {
	return &InMemoryStorage{
		links: store,
		mu:    sync.RWMutex{},
	}
}

func (s *InMemoryStorage) Store(shortenURL domain.ShortenURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.links[shortenURL.ID()] = shortenURL.OriginalURL()

	return nil
}

func (s *InMemoryStorage) GetOriginalURL(id domain.ID) (domain.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	originalURL, ok := s.links[id.String()]

	if !ok {
		return "", ErrURLNotFound
	}

	return domain.URL(originalURL), nil
}
