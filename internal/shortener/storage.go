package shortener

import (
	"sync"

	"emperror.dev/errors"
)

type (
	Storage interface {
		WriteOnlyStorage
		ReadOnlyStorage
	}

	WriteOnlyStorage interface {
		Store(lnk link) error
	}

	ReadOnlyStorage interface {
		GetOriginalURL(short shortID) (string, error)
	}

	InMemoryStorage struct {
		links map[string]string
		mx    sync.RWMutex
	}
)

func NewInMemoryStorage(store map[string]string) *InMemoryStorage {
	return &InMemoryStorage{
		links: store,
		mx:    sync.RWMutex{},
	}
}

func (s *InMemoryStorage) Store(lnk link) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.links[lnk.Short()] = lnk.Original()

	return nil
}

func (s *InMemoryStorage) GetOriginalURL(short shortID) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	original, ok := s.links[short.String()]
	if !ok {
		return "", errors.Errorf("original url not found by: %s", short.String())
	}

	return original, nil
}
