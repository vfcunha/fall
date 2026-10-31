package fall

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func FormDecoder[T any](r *http.Request) Result[*T] {
	if err := r.ParseForm(); err != nil {
		return NewResult((*T)(nil), fmt.Errorf("failed to parse form: %w", err))
	}

	var dto T
	val := reflect.ValueOf(&dto).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		formKey := field.Tag.Get("form")
		if formKey == "" {
			formKey = strings.ToLower(field.Name)
		}

		formValue := r.FormValue(formKey)
		if formValue == "" {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(formValue)
		case reflect.Bool:
			fieldVal.SetBool(formValue == "on" || formValue == "true" || formValue == "1")
		case reflect.Int, reflect.Int64:
			if v, err := strconv.ParseInt(formValue, 10, 64); err == nil {
				fieldVal.SetInt(v)
			} else {
				return NewResult((*T)(nil), fmt.Errorf("invalid int value for field %s: %s", field.Name, formValue))
			}
		case reflect.Uint, reflect.Uint64:
			if v, err := strconv.ParseUint(formValue, 10, 64); err == nil {
				fieldVal.SetUint(v)
			} else {
				return NewResult((*T)(nil), fmt.Errorf("invalid uint value for field %s: %s", field.Name, formValue))
			}
		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(formValue, 64); err == nil {
				fieldVal.SetFloat(v)
			} else {
				return NewResult((*T)(nil), fmt.Errorf("invalid float value for field %s: %s", field.Name, formValue))
			}
		default:
			return NewResult((*T)(nil), fmt.Errorf("unsupported field type for field %s: %s", field.Name, fieldVal.Kind()))
		}
	}

	instance := any(&dto)
	if validator, ok := instance.(validator); ok {
		if err := validator.Validate(); err != nil {
			return NewResult((*T)(nil), err)
		}
	}

	return NewResult(&dto, nil)
}
