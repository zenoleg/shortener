package shortener

import (
	"net/url"
	"strings"

	"emperror.dev/errors"
	"github.com/jxskiss/base62"
)

type (
	originalURL struct {
		original string
	}

	shortURL struct {
		encodedValue string
	}

	link struct {
		original originalURL
		short    shortURL
	}
)

func newOriginalURL(original string) (originalURL, error) {
	original = strings.TrimSpace(original)

	if original == "" {
		return originalURL{}, errors.New("original url is empty")
	}

	_, err := url.ParseRequestURI(original)

	if err != nil {
		return originalURL{}, err
	}

	return originalURL{original: original}, nil
}

func newShortURL(originalURL originalURL) shortURL {
	encoded := base62.Encode([]byte(originalURL.String()))

	return shortURL{
		encodedValue: string(encoded),
	}
}

func newLink(original string) (link, error) {
	originalValue, err := newOriginalURL(original)
	if err != nil {
		return link{}, errors.Errorf("url '%s' is invalid", original)
	}

	return link{
		original: originalValue,
		short:    newShortURL(originalValue),
	}, nil
}

func (u originalURL) String() string {
	return u.original
}

func (u shortURL) String() string {
	return u.encodedValue
}

func (l link) Original() string {
	return l.original.String()
}

func (l link) Short() string {
	return l.short.String()
}
