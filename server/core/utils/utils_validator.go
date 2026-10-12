package utils

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

// Validatable is a rich type that knows its own legal values; the `valid` rule
// asks it.
type Validatable interface {
	Valid() bool
}

var once sync.Once

// RegisterValidators adds Ivory's rules to gin's own validator, the one every bound
// request and every MCP tool input runs through. It is safe to call repeatedly.
func RegisterValidators() {
	once.Do(func() {
		engine := binding.Validator.Engine().(*validator.Validate)
		// it does the same as `required` plus trims spaces (so "   " such string will be treated as invalid)
		mustRegister(engine.RegisterValidation("notblank", validators.NotBlank))
		// it requires for the `type` implementation of Validatable interface
		mustRegister(engine.RegisterValidation("valid", isValid))
	})
}

func isValid(fl validator.FieldLevel) bool {
	field := fl.Field()
	if !field.CanInterface() {
		return false
	}
	rich, ok := field.Interface().(Validatable)
	if !ok {
		panic(fmt.Sprintf("utils: %s has the valid rule but %s has no Valid method", fl.StructFieldName(), field.Type()))
	}
	return rich.Valid()
}

func mustRegister(err error) {
	if err != nil {
		panic(err)
	}
}
