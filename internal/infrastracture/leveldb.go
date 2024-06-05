package infrastracture

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
)

type Config struct {
	Path string
}

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
