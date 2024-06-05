package shortener

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
)

var ErrNotFound = &NotFoundError{}

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

	LoggedStorage struct {
		storage Storage
		logger  zerolog.Logger
	}

	InMemoryStorage struct {
		links map[string]string
		mx    sync.RWMutex
	}

	LevelDBStorage struct {
		connection *leveldb.DB
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

func NewLevelDBStorage(connection *leveldb.DB, logger zerolog.Logger) *LoggedStorage {
	return &LoggedStorage{
		storage: &LevelDBStorage{
			connection: connection,
		},
		logger: logger,
	}
}

func (s *LevelDBStorage) Store(lnk link) error {
	return s.connection.Put([]byte(lnk.ShortID().String()), []byte(lnk.Original()), nil)
}

func (s *LevelDBStorage) GetOriginalURL(short shortID) (string, error) {
	res, err := s.connection.Get([]byte(short.String()), nil)
	if err != nil {
		return "", ErrNotFound
	}

	return string(res), nil
}

func (s *LoggedStorage) Store(lnk link) error {
	err := s.storage.Store(lnk)
	if err != nil {
		s.logger.Err(err).Msg("failed to store link")

		return err
	}

	s.logger.Info().
		Str("short_id", lnk.ShortID().String()).
		Str("original", lnk.Original()).
		Msg("link stored")

	return nil
}

func (s *LoggedStorage) GetOriginalURL(short shortID) (string, error) {
	original, err := s.storage.GetOriginalURL(short)
	if err != nil {
		s.logger.Err(err).Msg("failed to fetch original link")

		return original, err
	}

	s.logger.Info().
		Str("short_id", short.String()).
		Msg("link fetched")

	return original, err
}

func (e NotFoundError) Error() string {
	return e.msg
}
