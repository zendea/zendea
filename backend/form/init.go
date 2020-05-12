package form

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/go-playground/validator/v10"
	"strings"
)

func init() {
	// Register custom validate methods
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//_ = v.RegisterValidation("customValidate", customValidate)
		// _ = v.RegisterValidation("pwdValidate", pwdValidate)
	// }
}

// Bind : bind request dto and auto verify parameters
func Bind(c *gin.Context, obj interface{}) error {
	_ = c.ShouldBindUri(obj)
	if err := c.ShouldBind(obj); err != nil {
		var tagErrorMsg []string
		for _, e := range err.(validator.ValidationErrors) {
			if _, has := ValidateErrorMessage[e.Tag()]; has {
				tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(ValidateErrorMessage[e.Tag()], e.Field(), e.Value()))
			} else {
				tagErrorMsg = append(tagErrorMsg,fmt.Sprintf(ValidateErrorMessage["default"], e.Tag(), e.Field(), e.Value()))
			}
		}
		return errors.New(strings.Join(tagErrorMsg, ","))
	}

	return nil
}

//ValidateErrorMessage : customize error messages
var ValidateErrorMessage = map[string]string{
	"default":        "%s - %s is invalid(%s)",
	"customValidate": "%s can not be %s",
	"required":       "%s is required,got empty %#v",
	"pwdValidate":    "%s is not a valid password",
}
