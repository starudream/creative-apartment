package ivalidator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	validateOnce sync.Once
)

func V() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()
	})
	return validate
}

func Struct(v any) error {
	return wrapError(V().Struct(v))
}

func Var(v any, name, rule string) error {
	return wrapError(V().Struct(newStruct(v, name, rule)))
}

func newStruct(v any, name, rule string) any {
	f := reflect.StructField{
		Name:  name,
		Type:  reflect.TypeOf(v),
		Tag:   reflect.StructTag(fmt.Sprintf(`validate:"%s"`, rule)),
		Index: []int{0},
	}
	nv := reflect.New(reflect.StructOf([]reflect.StructField{f}))
	nv.Elem().Field(0).Set(reflect.ValueOf(v))
	return nv.Interface()
}

type Errors struct {
	es validator.ValidationErrors
}

var _ error = (*Errors)(nil)

func (e Errors) Error() string {
	sb := strings.Builder{}
	for i := 0; i < len(e.es); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("validation for '%s' failed on the '%s' tag", e.es[i].Field(), e.es[i].Tag()))
	}
	return strings.TrimSpace(sb.String())
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	switch v := err.(type) {
	case validator.ValidationErrors:
		err = &Errors{es: v}
	}
	return err
}
