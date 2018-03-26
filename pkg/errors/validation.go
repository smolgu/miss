package errors

import (
	"fmt"
	"log"

	"github.com/go-playground/validator"
)

// FormatValidatorError check user input and returns human readable errors if
// exists
func FormatValidatorError(err error) error {
	if err == nil {
		return nil
	}
	log.Printf("%+v", err)
	if verr, ok := err.(validator.ValidationErrors); ok {
		if len(verr) != 1 {
			return verr
		}
		return GetValidationError(verr[0].Field(), verr[0].Tag())
	} else {
		return fmt.Errorf("cannot cast error %+v", err)
	}
	return err
}

func GetValidationError(fieldName string, tag string) error {
	return fmt.Errorf(errMap[tag], russianFieldTranslate(fieldName))
}

var errMap = map[string]string{
	"required": "отсувствует %v",
}

func russianFieldTranslate(tagName string) string {
	if russianName, found := russianTags[tagName]; found {
		return russianName
	}
	return tagName
}

var russianTags = map[string]string{
	"Token": "token",
}
