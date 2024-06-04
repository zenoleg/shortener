package transport

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}

	ShortenResponse struct {
		Destination string `json:"destination"`
	}
)

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Message: msg}
}

func NewShortenResponse(shortURL string) ShortenResponse {
	return ShortenResponse{Destination: shortURL}
}
