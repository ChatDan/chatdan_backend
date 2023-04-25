package fiberx

import (
	"ChatDanBackend/common/schemax"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MyErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	httpError := schemax.Response{
		Code:     500,
		ErrorMsg: err.Error(),
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError.Code = 404
	} else {
		switch e := err.(type) {
		case *schemax.Response:
			httpError = *e
		case *fiber.Error:
			httpError.Code = e.Code
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
