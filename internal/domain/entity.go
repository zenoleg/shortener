package domain

import (
	"errors"
	"net/url"
	"strings"
)

type (
	ID  string
	URL string

	ShortenURL struct {
		id          ID
		originalURL URL
	}
)

func NewID(id string) (ID, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return "", errors.New("ID must be not empty")
	}

	return ID(id), nil
}

func NewURL(originalURL string) (URL, error) {
	originalURL = strings.TrimSpace(originalURL)

	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("URL is not valid")
	}

	return URL(originalURL), nil
}

func NewShortenURL(id ID, originalURL URL) ShortenURL {
	return ShortenURL{
		id:          id,
		originalURL: originalURL,
	}
}

func (i ID) String() string {
	return string(i)
}

func (u URL) String() string {
	return string(u)
}

func (s ShortenURL) ID() string {
	return s.id.String()
}

func (s ShortenURL) OriginalURL() string {
	return s.originalURL.String()
}
