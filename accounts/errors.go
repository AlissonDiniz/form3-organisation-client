package accounts

import (
	"fmt"
)

type error interface {
	Error() string
}

type ReadResponseBodyError struct {
	ErrorMessage string
	StackTrace string
}
func (e *ReadResponseBodyError) Error() string {
	return fmt.Sprintf("Read Response Body Error: %v StackTrace: %v", e.ErrorMessage, e.StackTrace)
}

type JsonParseError struct {
	ErrorMessage string
	StackTrace string
}
func (e *JsonParseError) Error() string {
	return fmt.Sprintf("Json Parse Error: %v StackTrace: %v", e.ErrorMessage, e.StackTrace)
}

type NotFoundError struct {
	ErrorMessage    string `json:"error_message"`
}
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found Error: %v", e.ErrorMessage)
}

type ConflictError struct {
	ErrorMessage    string `json:"error_message"`
}
func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict Error: %v", e.ErrorMessage)
}

type InternalServerError struct {
	ErrorMessage    string `json:"error_message"`
}
func (e *InternalServerError) Error() string {
	return fmt.Sprintf("Internal Server Error: %v", e.ErrorMessage)
}
