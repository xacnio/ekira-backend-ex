package errs

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type (
	Error struct {
		StatusCode  int
		Code        uint16
		Description string
		Message     string
		Error       error
		DError      error
	}
)

const ERR_NOTFOUND uint16 = 2000
const ERR_FORBIDDEN uint16 = 2001

const ERR_DATABASE uint16 = 1000
const ERR_PARSE uint16 = 1001
const ERR_VALIDATE uint16 = 1002
const ERR_DATABASEQUERY uint16 = 1003
const ERR_JWT uint16 = 1004
const ERR_LOGIN uint16 = 1005
const ERR_REGISTER uint16 = 1006
const ERR_NOAUTH uint16 = 1007
const ERR_INVALIDIMAGE uint16 = 1008
const ERR_IMAGEUPLOAD uint16 = 1009
const ERR_IMAGEPROC uint16 = 1010

func (e *Error) SetMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) Set(msg string) *Error {
	e.Description = msg
	return e
}

func (e *Error) SetDError(error error) *Error {
	e.DError = error
	return e
}

func NewError(statusCode int, code uint16, msgs ...string) *Error {
	err := &Error{StatusCode: statusCode, Code: code, Error: errors.New(msgs[0])}
	if len(msgs) >= 1 {
		err.Description = msgs[0]
	}
	if len(msgs) >= 2 {
		err.Message = msgs[1]
	}
	return err
}

var (
	ErrDatabaseConnection = NewError(fiber.StatusInternalServerError, ERR_DATABASE, "not connected to database")
	ErrDatabaseQuery      = NewError(fiber.StatusInternalServerError, ERR_DATABASEQUERY, "database query error")
	ErrBadRequest         = NewError(fiber.StatusBadRequest, ERR_PARSE, "bad request")
	ErrForbidden          = NewError(fiber.StatusForbidden, ERR_FORBIDDEN, "forbidden")
	ErrValidate           = NewError(fiber.StatusBadRequest, ERR_VALIDATE, "validation error")
	ErrJwt                = NewError(fiber.StatusInternalServerError, ERR_JWT, "jwt error")
	ErrLogin              = NewError(fiber.StatusOK, ERR_LOGIN, "invalid email or password")
	ErrRegister           = NewError(fiber.StatusOK, ERR_REGISTER, "email is already registered")
	ErrNoAuth             = NewError(fiber.StatusUnauthorized, ERR_NOAUTH, "no authorization")
	ErrNotFound           = NewError(fiber.StatusNotFound, ERR_NOTFOUND, "not found")
	ErrInvalidImage       = NewError(fiber.StatusBadRequest, ERR_INVALIDIMAGE, "invalid image")
	ErrUploadImage        = NewError(fiber.StatusInternalServerError, ERR_IMAGEUPLOAD, "file error")
	ErrImageProc          = NewError(fiber.StatusInternalServerError, ERR_IMAGEPROC, "file operation error")
)
