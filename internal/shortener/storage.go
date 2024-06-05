package shortener

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
)

var ErrNotFound *NotFoundError = &NotFoundError{}

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

	LevelDBStorage struct {
		connection *leveldb.DB
		logger     zerolog.Logger
	}

	NotFoundError struct {
		msg string
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

	s.links[lnk.ShortID().String()] = lnk.Original()

	return nil
}

func (s *InMemoryStorage) GetOriginalURL(short shortID) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	original, ok := s.links[short.String()]
	if !ok {
		return "", ErrNotFound
	}

	return original, nil
}

func NewLevelDBStorage(connection *leveldb.DB, logger zerolog.Logger) *LevelDBStorage {
	return &LevelDBStorage{
		connection: connection,
		logger:     logger,
	}
}

func (s *LevelDBStorage) Store(lnk link) error {
	return nil
}

func (s *LevelDBStorage) GetOriginalURL(short shortID) (string, error) {
	return "", nil
}

func (e NotFoundError) Error() string {
	return e.msg
}
