package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

// ValidationErrors.Error converts slice of errors into single string
func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for i, err := range v {
		b.WriteString(err.Field)
		b.WriteString(": ")
		b.WriteString(err.Err.Error())
		if len(v) != 0 && i == len(v)-1 {
			b.WriteString(" | ")
		}

	}
	return b.String()
}

func Validate(v interface{}) error {
	var errs ValidationErrors
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Struct {
		return nil
	}
	structRval := rVal.Type()
	for i := 0; i < structRval.NumField(); i++ {
		fld := structRval.Field(i)
		var (
			fieldName  = fld.Name
			fieldType  = fld.Type
			fieldTag   = fld.Tag
			fieldValue = rVal.Field(i)
		)
		fmt.Println(
			"Fieldname: ", fieldName,
			"\nFieldValue: ", fieldValue,
			"\nType: ", fieldType,
			"\nTag: ", fieldTag,
		)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
