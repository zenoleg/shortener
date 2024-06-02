package transport

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}

	ShortenResponse struct {
		ShortURL string `json:"short_url"`
	}
)

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Message: msg}
}

func NewShortenResponse(shortURL string) ShortenResponse {
	return ShortenResponse{ShortURL: shortURL}
}
