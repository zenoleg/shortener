package storage

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zenoleg/shortener/internal/domain"
)

type (
	Config struct {
		Path string
	}

	LevelDBStorage struct {
		connection *leveldb.DB
		logger     zerolog.Logger
	}
)

func NewConfig() Config {
	path, exists := os.LookupEnv("LEVEL_DB_PATH")
	if !exists {
		path = "/tmp"
	}

	return Config{
		Path: path,
	}
}

func NewLevelDBConnection(cfg Config, logger zerolog.Logger) (*leveldb.DB, func(), error) {
	path := fmt.Sprintf("%s/links", cfg.Path)

	db, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, func() {
			logger.Error().Err(err).Msg("can not open level db")
			_ = db.Close()
		}, err
	}

	logger.Info().Msgf("level db successfully opened on '%s'", path)

	return db, func() {
		_ = db.Close()

		logger.Debug().Msg("level db closed")
	}, nil
}

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
