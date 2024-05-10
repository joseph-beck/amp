package binding

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Binder interface {
	Name() string
	Bind(*http.Request, any) error
}

var (
	JSON = jsonBinding{}
)

func readBody(request *http.Request) (*bytes.Buffer, error) {
	if request.Body == nil {
		return nil, errors.New("invalid request")
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(body)
	request.Body = io.NopCloser(buff)
	new := bytes.NewBuffer(body)
	return new, nil
}

type StructValidator interface {
	ValidateStruct(any) error
	Engine() any
}

var Validator StructValidator = &defaultValidator{}

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}

type SliceValidationError []error

func (err SliceValidationError) Error() string {
	n := len(err)

	switch n {
	case 0:
		return ""
	default:
		var builder strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&builder, "[%d]: %s", 0, err[0].Error())
		}

		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					builder.WriteString("\n")
					fmt.Fprintf(&builder, "[%d]: %s", i, err[i].Error())
				}
			}
		}

		return builder.String()
	}
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyInit()
	return v.validate.Struct(obj)
}

func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(val.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := val.Len()
		validate := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(val.Index(i).Interface()); err != nil {
				validate = append(validate, err)
			}
		}

		if len(validate) == 0 {
			return nil
		}

		return validate
	default:
		return nil
	}
}

func (v *defaultValidator) Engine() any {
	v.lazyInit()
	return v.validate
}
