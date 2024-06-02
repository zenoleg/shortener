package transport

import "os"

type Config struct {
	Address string
}

func NewConfig() Config {
	address, exists := os.LookupEnv("HTTP_ADDRESS")
	if !exists {
		address = ":0"
	}

	return Config{
		Address: address,
	}
}
