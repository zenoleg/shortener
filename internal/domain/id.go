package domain

import "github.com/jxskiss/base62"

type (
	Base62IDGenerator struct{}
)

func NewBase62IDGenerator() Base62IDGenerator {
	return Base62IDGenerator{}
}

func (g Base62IDGenerator) Generate(originalURL URL) (ID, error) {
	return NewID(string(base62.Encode([]byte(originalURL.String()))))
}
