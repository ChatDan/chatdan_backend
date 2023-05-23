package utils

import (
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type ErrorDetailElement struct {
	Field    string       `json:"field"`
	Tag      string       `json:"tag"`
	TypeKind reflect.Kind `json:"-"`
	Value    string       `json:"value"`
	ErrorMsg string       `json:"error_msg"`
}

func (e *ErrorDetailElement) Error() string {
	if e.ErrorMsg != "" {
		return e.ErrorMsg
	}
	switch e.Tag {
	case "min":
		if e.TypeKind == reflect.String {
			e.ErrorMsg = e.Field + "至少" + e.Value + "字符"
		} else {
			e.ErrorMsg = e.Field + "最小值为" + e.Value
		}
	case "max":
		if e.TypeKind == reflect.String {
			e.ErrorMsg = e.Field + "限长" + e.Value + "字符"
		} else {
			e.ErrorMsg = e.Field + "最大值为" + e.Value
		}
	case "required":
		e.ErrorMsg = e.Field + "不能为空"
	case "email":
		e.ErrorMsg = e.Field + "格式不正确"
	case "modify":
		e.ErrorMsg = "请求体不能为空"
	}
	return e.ErrorMsg
}

type ErrorDetail []*ErrorDetailElement

func (e *ErrorDetail) Error() string {
	return "Validation Error"
}

var Validate = validator.New()

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateStruct(model any) error {
	errors := Validate.Struct(model)
	if errors != nil {
		var errorDetail ErrorDetail
		for _, err := range errors.(validator.ValidationErrors) {
			detail := ErrorDetailElement{
				Field:    err.Field(),
				Tag:      err.Tag(),
				TypeKind: err.Type().Kind(),
				Value:    err.Param(),
			}
			errorDetail = append(errorDetail, &detail)
		}
		return &errorDetail
	}
	return nil
}

func ValidateQuery(c *fiber.Ctx, model any) error {
	if err := c.QueryParser(model); err != nil {
		return err
	}
	if err := defaults.Set(model); err != nil {
		return err
	}
	return ValidateStruct(model)
}

// ValidateBody supports json only
func ValidateBody(c *fiber.Ctx, model any) error {
	body := c.Body()
	if len(body) == 0 || string(body) == "{}" {
		return BadRequest("Body is empty")
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	if err := defaults.Set(model); err != nil {
		return err
	}
	return ValidateStruct(model)
}
