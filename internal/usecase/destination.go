package usecase

import "fmt"

type DestinationURL string

func NewDestinationURL(isSSL bool, host string, id string) DestinationURL {
	protocol := "http"
	if isSSL {
		protocol = "https"
	}

	return DestinationURL(fmt.Sprintf("%s://%s/link/%s", protocol, host, id))
}

func (u DestinationURL) String() string {
	return string(u)
}
