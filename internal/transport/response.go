package transport

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}

	DestinationResponse struct {
		Destination string `json:"destination"`
	}
)

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Message: msg}
}

func NewShortenResponse(shortURL string) DestinationResponse {
	return DestinationResponse{Destination: shortURL}
}
