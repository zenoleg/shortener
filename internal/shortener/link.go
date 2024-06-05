package shortener

import (
	"net/url"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/jxskiss/base62"
)

type (
	shortID struct {
		encoded string
	}

	shortLink struct {
		link string
	}

	originalURL struct {
		original string
	}

	link struct {
		shortID  shortID
		original originalURL
	}
)

func newShortID(originalURL originalURL) shortID {
	encoded := base62.Encode([]byte(originalURL.String()))

	return shortID{
		encoded: string(encoded),
	}
}

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

func newShortLink(link string) (shortLink, error) {
	link = strings.TrimSpace(link)
	link = strings.TrimSuffix(link, "/")

	pattern := regexp.MustCompile(`/link/[[:word:]]+$`)
	if !pattern.MatchString(link) {
		return shortLink{}, errors.New("short link is invalid")
	}

	return shortLink{
		link: link,
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

func (l link) ShortID() shortID {
	return l.shortID
}

func (l shortLink) shortID() shortID {
	split := strings.Split(l.link, "/")
	id := split[len(split)-1]

	return shortID{encoded: id}
}
