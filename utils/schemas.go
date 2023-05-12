package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PageRequest struct {
	PageNum  int `json:"page_num" query:"page_num" validate:"required,min=1"`
	PageSize int `json:"page_size" query:"page_size" validate:"required,min=1,max=100"`
}

func (q PageRequest) QuerySet(tx *gorm.DB) *gorm.DB {
	return tx.Offset((q.PageNum - 1) * q.PageSize).Limit(q.PageSize)
}

type Response struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"error_msg"`
	Data     any    `json:"data,omitempty"`
}

func (r Response) Error() string {
	return r.ErrorMsg
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(200).JSON(Response{
		Code: 200,
		Data: data,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(201).JSON(Response{
		Code: 201,
		Data: data,
	})
}
