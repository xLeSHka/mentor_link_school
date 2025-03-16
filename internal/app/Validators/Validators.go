package Validators

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"unicode"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := &Validator{validator: validator.New()}
	if err := v.validator.RegisterValidation("c-password", validatePassword); err != nil {
		panic("Error create c-password validator")
	}
	if err := v.validator.RegisterValidation("c-role", validateRole); err != nil {
		panic("Error create c-role validator")
	}
	//if err := v.validator.RegisterValidation("c-country", validateCountry, true); err != nil {
	//	panic("Error create c-country validator")
	//}
	return v
}
func validateRole(field validator.FieldLevel) bool {
	role := field.Field().String()
	if role != "student" && role != "owner" && role != "mentor" {
		return false
	}
	return true
}
func validatePassword(field validator.FieldLevel) bool {
	fmt.Println("tryValidate: " + field.Field().String())
	fmt.Println("Type: " + fmt.Sprintf("%T", field.Field().String()))
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, char := range field.Field().String() {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.In(char, unicode.Symbol, unicode.Punct):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}
func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

func (v *Validator) ShouldBindJSON(c *gin.Context, i any) error {
	if err := c.ShouldBindBodyWith(i, binding.JSON); err != nil {
		return err
	}
	return v.Validate(i)
}
func (v *Validator) ShouldBindQuery(c *gin.Context, i any) error {
	if err := c.ShouldBindQuery(i); err != nil {
		return err
	}
	return v.Validate(i)
}
func (v *Validator) ShouldBindUri(c *gin.Context, i any) error {
	if err := c.ShouldBindUri(i); err != nil {
		return err
	}
	return v.Validate(i)
}
