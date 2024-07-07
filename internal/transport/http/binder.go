package http

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type binder struct {
	db echo.DefaultBinder
}

func NewValidationBinder() echo.Binder {
	return &binder{}
}

func (b *binder) Bind(i interface{}, c echo.Context) error {
	err := b.db.Bind(i, c)
	if err != nil {
		return err
	}

	if val, ok := i.(validation.Validatable); ok {
		return val.Validate()
	}

	return nil
}
