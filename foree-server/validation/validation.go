package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"xue.io/go-pay/server/transport"
)

var validate = validator.New()

func ValidateStruct(s any, errMsg string) *transport.BadRequestError {
	ret := transport.NewFormError(errMsg)
	if err := validate.Struct(s); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			ret.AddDetails(e.Field(), fmt.Sprintf("Invalid %s", e.Field()))
		}
	}
	return ret
}
