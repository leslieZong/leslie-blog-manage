package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate 是整个项目共用的 validator 实例。
//
// validator 本身是可以复用的，所以项目启动以后
// 不需要每次请求都重新创建一个 validator。
var Validate = validator.New()

// FormatErrors 将 validator 返回的错误
// 转换成更容易返回给前端的字符串。
//
// 例如：
//
//	Username string `validate:"required,min=3,max=50"`
//
// 如果 username 为空，validator 可能返回比较底层的错误信息。
//
// 我们最终希望前端看到：
//
//	username is required
//
// 而不是一大堆 validator 内部信息。
func FormatErrors(err error) string {

	// 如果不是 validator.ValidationErrors，
	// 说明可能是其他类型的错误。
	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return err.Error()
	}

	messages := make([]string, 0, len(validationErrors))

	for _, fieldError := range validationErrors {

		fieldName := fieldError.Field()

		switch fieldError.Tag() {

		case "required":
			messages = append(
				messages,
				fieldName+" is required",
			)

		case "min":
			messages = append(
				messages,
				fieldName+" is too short",
			)

		case "max":
			messages = append(
				messages,
				fieldName+" is too long",
			)

		case "email":
			messages = append(
				messages,
				fieldName+" must be a valid email",
			)

		default:
			messages = append(
				messages,
				fieldName+" is invalid",
			)
		}
	}

	return strings.Join(messages, "; ")
}
