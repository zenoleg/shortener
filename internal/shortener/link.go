package shortener

import (
	"net/url"
	"strings"

	"emperror.dev/errors"
	"github.com/jxskiss/base62"
)

type (
	shortID struct {
		encoded string
	}

	originalURL struct {
		original string
	}

	link struct {
		shortID  shortID
		original originalURL
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

func newShortID(originalURL originalURL) shortID {
	encoded := base62.Encode([]byte(originalURL.String()))

	return shortID{
		encoded: string(encoded),
	}
}

func newLink(original string) (link, error) {
	originalValue, err := newOriginalURL(original)
	if err != nil {
		return link{}, errors.Errorf("url '%s' is invalid", original)
	}

	return link{
		shortID:  newShortID(originalValue),
		original: originalValue,
	}, nil
}

func (u originalURL) String() string {
	return u.original
}

func (u shortID) String() string {
	return u.encoded
}

func (l link) Original() string {
	return l.original.String()
}

func (l link) Short() string {
	return l.shortID.String()
}
