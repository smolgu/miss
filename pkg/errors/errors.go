package errors

import (
	"fmt"
	"reflect"
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
	if e, ok := err.(*Error); ok {
		return reflect.DeepEqual(e.Typed(), typed)
	}
	return reflect.DeepEqual(err, typed)
}

// Sprint check is err Error and print dev msg
func Sprint(err error) string {
	if e, ok := err.(*Error); ok && e.devMsg != nil {
		return e.devMsg.Error()
	}
	return err.Error()
}

// func IsZeroOfUnderlyingType(x interface{}) bool {
// 	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
// }
