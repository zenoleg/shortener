package transport

import validation "github.com/go-ozzo/ozzo-validation/v4"

type (
	ShortenRequest struct {
		URL string `json:"url" query:"url"`
	}

	OriginalRequest struct {
		ShortID string `query:"short_id"`
	}
)

func (r ShortenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required),
	)
}

func (r OriginalRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ShortID, validation.Required),
	)
}
