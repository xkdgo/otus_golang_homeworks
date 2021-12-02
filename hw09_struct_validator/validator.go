package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var validatorKey = "validate"

var (
	ErrInvalidValidator = errors.New("validator should be func:limit")
	ErrExtractValidator = errors.New("some extractor error")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

// ValidationErrors.Error converts slice of errors into single string.
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

// This function validate fields of struct, if they have `validatorKey flag.
// By default this validatorKey="validate".
// ...
func Validate(v interface{}) error {
	var errs ValidationErrors
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Struct {
		return nil
	}
	structRval := rVal.Type()
	errs = make(ValidationErrors, 0, structRval.NumField())
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
		val, ok := fieldTag.Lookup(validatorKey)
		if !ok {
			continue
		}
		if val == "" {
			continue
		}
		fmt.Println("TagValue= ", val)
		extrValMap, err := extractValidators(val)
		if err != nil {
			return err // TODO make correct validation errors
		}
		fmt.Printf("%#v\n", extrValMap)
	}
	if len(errs) == 0 {
		return errs // TODO return nil
	}
	return errs
}

func extractValidators(val string) (map[string]string, error) {
	if !strings.Contains(val, ":") {
		return nil, ErrInvalidValidator
	}
	validateCandidates := strings.Split(val, "|")
	extractedMap := make(map[string]string, len(validateCandidates))
	for _, candidate := range validateCandidates {
		keyWithVal := strings.Split(candidate, ":")
		if len(keyWithVal) != 2 {
			return nil, ErrInvalidValidator
		}
		extractedMap[strings.Trim(keyWithVal[0], " ")] = keyWithVal[1]
	}
	return extractedMap, nil
}
