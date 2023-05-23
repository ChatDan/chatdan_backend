package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"reflect"
)

type PageRequest struct {
	PageNum  int `json:"page_num" query:"page_num" validate:"required,min=1"`
	PageSize int `json:"page_size" query:"page_size" validate:"required,min=1,max=100"`
	Version  int `json:"version" query:"version" validate:"omitempty,min=0"` // 分页版本号，一个时间戳，用于保证分页查询的一致性和正确性。不填默认使用最新版本时间戳
}

func (q PageRequest) QuerySet(tx *gorm.DB) *gorm.DB {
	return tx.Offset((q.PageNum - 1) * q.PageSize).Limit(q.PageSize)
}

type CanPostprocess interface {
	Postprocess(c *fiber.Ctx) error
}

type Response[T any] struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"error_msg"`
	Data     *T     `json:"data,omitempty"`
}

type RespForSwagger struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"error_msg"`
	Data     any    `json:"data,omitempty"`
}

func (r Response[T]) Error() string {
	return r.ErrorMsg
}

func Success[T any](c *fiber.Ctx, data *T) error {
	err := postprocess(c, data)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(Response[T]{
		Code: 200,
		Data: data,
	})
}

func Created[T any](c *fiber.Ctx, data *T) error {
	err := postprocess(c, data)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(Response[T]{
		Code: 201,
		Data: data,
	})
}

func postprocess[T any](c *fiber.Ctx, data *T) error {
	modelType := reflect.TypeOf(data)

	modelValue := reflect.ValueOf(data)

	method, ok := modelType.MethodByName("Postprocess")
	if !ok {
		return nil
	}

	returns := method.Func.Call([]reflect.Value{modelValue, reflect.ValueOf(c)})

	if len(returns) == 0 || returns[0].IsNil() {
		return nil
	}

	return returns[0].Interface().(error)
}
