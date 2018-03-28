package errors

import "fmt"

// ErrNotFound describe not found error
var ErrNotFound = fmt.Errorf("not found")

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
	return e.devMsg.Error()
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

//
// import (
// 	"fmt"
//
// )
//
// var (
// 	notFoundFmt = "объект %v не найден"
// )
//
// func New(objectType int, errorCode int, args ...interface{}) Error {
// 	return Error{
// 		code:       errorCode,
// 		objectType: objectType,
// 		args:       args,
// 	}
// }
//
// func NewNotFound(object interface{}) Error {
// 	switch v := object.(type) {
// 	case types.User:
// return Error{
//   objectType:
//   return Error{}
// }
// 	}
// }
//
// type Error struct {
// 	code       int
// 	objectType int
// 	args       []interface{}
// }
//
// func (e Error) Error() string {
// 	switch e.code {
// 	case int(ErrorCode_NotFound):
// 		switch e.objectType {
// 		case int(ObjectType_ObjectUser):
// 			return fmt.Sprintf("Пользователь %v не найден", e.args...)
// 		}
// 	default:
// 		return "неизвестная ошибка"
// 	}
// 	return "unknown error"
// }
