package models

import "fmt"

var NotFound = fmt.Errorf("not found")

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
