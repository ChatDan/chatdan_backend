package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/juju/errors"
	"gorm.io/gorm"
)

func BadRequest(messages ...string) *Response {
	message := "Bad Request"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &Response{
		Code:     400,
		ErrorMsg: message,
	}
}

func Unauthorized(messages ...string) *Response {
	message := "Invalid JWT Token"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &Response{
		Code:     401,
		ErrorMsg: message,
	}
}

func Forbidden(messages ...string) *Response {
	message := "您没有权限进行此操作"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &Response{
		Code:     403,
		ErrorMsg: message,
	}
}

func NotFound(messages ...string) *Response {
	message := "Not Found"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &Response{
		Code:     404,
		ErrorMsg: message,
	}
}

func InternalServerError(messages ...string) *Response {
	message := "Unknown Error"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &Response{
		Code:     500,
		ErrorMsg: message,
	}
}

func MyErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	httpError := Response{
		Code:     500,
		ErrorMsg: err.Error(),
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError.Code = 404
	} else {
		switch e := err.(type) {
		case *Response:
			httpError = *e
		case *fiber.Error:
			httpError.Code = e.Code
		case *ErrorDetail:
			httpError.Code = 400
			httpError.Data = e
		case fiber.MultiError:
			httpError.Code = 400
			httpError.ErrorMsg = ""
			for _, err = range e {
				httpError.ErrorMsg += err.Error() + "\n"
			}
		}
	}

	return c.Status(httpError.Code).JSON(&httpError)
}
