package Validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := &Validator{validator: validator.New()}
	//if err := v.validator.RegisterValidation("c-password", validatePassword); err != nil {
	//	panic("Error create c-password validator")
	//}
	//if err := v.validator.RegisterValidation("c-country", validateCountry, true); err != nil {
	//	panic("Error create c-country validator")
	//}
	return v
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
