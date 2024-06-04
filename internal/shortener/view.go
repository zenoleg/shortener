package shortener

import "fmt"

type (
	DestinationURL struct {
		url string
	}
)

func newDestinationURL(isSSL bool, host string, short shortID) DestinationURL {
	protocol := "http"
	if isSSL {
		protocol = "https"
	}

	return DestinationURL{
		url: fmt.Sprintf("%s://%s/link/%s", protocol, host, short.String()),
	}
}

func (u DestinationURL) String() string {
	return u.url
}
