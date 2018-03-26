package errors

import (
	"fmt"
)

var (
	// ErrNotFound describe not found error
	ErrNotFound = fmt.Errorf("not found")
	// ErrPublicUnknownError will be displayed for end user. Its hide dev info
	ErrPublicUnknownError = fmt.Errorf("неизвестная ошибка")
)

// New create new Error with description for user and developer
func New(userMsg, devMsg error, typeds ...error) error {
	err := &Error{
		userMsg: userMsg,
		devMsg:  devMsg,
	}
	if len(typeds) > 0 {
		err.typed = typeds[0]
	}

	return err
}

// Error main object of this package
type Error struct {
	userMsg error
	devMsg  error
	typed   error
}

func (e Error) Error() string {
	return e.userMsg.Error()
}

// Typed returns error with type if has. Else return himself
func (e Error) Typed() error {
	if e.typed != nil {
		return e.typed
	}
	return e
}

// CheckTyped cast err to Error and check Error.Typed == typed
func CheckTyped(err error, typed error) bool {
	if e, ok := err.(Error); ok {
		return e.Typed() == typed
	}
	return err == typed
}
