package storage

import (
	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zenoleg/shortener/internal/domain"
)

type (
	LevelDBStorage struct {
		connection *leveldb.DB
		logger     zerolog.Logger
	}
)

func NewLevelDBStorage(connection *leveldb.DB, logger zerolog.Logger) *LevelDBStorage {
	return &LevelDBStorage{
		connection: connection,
		logger:     logger,
	}
}

func (s LevelDBStorage) Store(shortenURL domain.ShortenURL) error {
	key := shortenURL.ID()
	value := shortenURL.OriginalURL()

	return s.connection.Put([]byte(key), []byte(value), nil)
}

func (s LevelDBStorage) GetOriginalURL(id domain.ID) (domain.URL, error) {
	res, err := s.connection.Get([]byte(id.String()), nil)
	if err != nil {
		return "", ErrURLNotFound
	}

	return domain.URL(res), nil
}
