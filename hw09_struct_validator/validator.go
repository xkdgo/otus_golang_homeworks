package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/xkdgo/otus_golang_homeworks/hw09_struct_validator/valuerror"
)

var validatorKey = "validate"

var ErrExtractValidator = errors.New("some extractor error")

type ErrInvalidValidator struct {
	message string
}

func (e ErrInvalidValidator) Error() string {
	return e.message
}

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	b := strings.Builder{}
	b.WriteString(v.Field)
	b.WriteString(": ")
	b.WriteString(v.Err.Error())
	return b.String()
}

type validationFunc func(string, string, reflect.Value, reflect.Type) error

var validationFuncMap = map[string]validationFunc{
	"len":    validateLen,
	"min":    validateMin,
	"max":    validateMax,
	"in":     validateIn,
	"regexp": validateRegexp,
}

func validateLen(limit string, structFieldName string, valueToCheck reflect.Value, valueType reflect.Type) error {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ErrInvalidValidator{fmt.Sprintf("%s should be integer", limit)}
	}
	switch {
	case valueType.String() == "string":
		if len(valueToCheck.String()) != limitInt {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateLen{TrueLimit: limitInt, ActualValue: len(valueToCheck.String())},
			}
		}
	case valueType.String() == "[]string":
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]string)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("Could not assert %#v to string slice", elemInterface)}
		}
		for _, elem := range elemSlice {
			if len(elem) != limitInt {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateLen{TrueLimit: limitInt, ActualValue: len(elem)},
				}
			}
		}
	default:
		return ErrInvalidValidator{fmt.Sprintf(
			"Field %s unsupported type %s for method \"len\"", structFieldName, valueType)}
	}
	return nil
}

func validateMin(limit string, structFieldName string, valueToCheck reflect.Value, valueType reflect.Type) error {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ErrInvalidValidator{fmt.Sprintf("%s should be integer", limit)}
	}
	switch {
	case valueType.String() == "int":
		if int(valueToCheck.Int()) < limitInt {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateMin{TrueLimit: limitInt, ActualValue: int(valueToCheck.Int())},
			}
		}
	case valueType.String() == "[]int":
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]int)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("Could not assert %#v to int slice", elemInterface)}
		}
		for _, elem := range elemSlice {
			if elem < limitInt {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateMin{TrueLimit: limitInt, ActualValue: elem},
				}
			}
		}
	default:
		return ErrInvalidValidator{fmt.Sprintf(
			"Field %s unsupported type %s for method \"min\"", structFieldName, valueType)}
	}
	return nil
}

func validateMax(limit string, structFieldName string, valueToCheck reflect.Value, valueType reflect.Type) error {
	return nil
}

func validateIn(limit string, structFieldName string, valueToCheck reflect.Value, valueType reflect.Type) error {
	return nil
}

func validateRegexp(limit string, structFieldName string, valueToCheck reflect.Value, valueType reflect.Type) error {
	return nil
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
	var valerr ValidationError
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
		for key, limit := range extrValMap {
			validationFn, ok := validationFuncMap[key]
			if !ok {
				log.Fatalf("validator %s not implemented", key)
			}
			err = validationFn(limit, fieldName, fieldValue, fieldType)
			if err != nil {
				ok = errors.As(err, &valerr)
				if !ok {
					return err
				}
				errs = append(errs, valerr)
			}
		}
	}
	if len(errs) == 0 {
		return errs // TODO return nil
	}
	return errs
}

func extractValidators(val string) (map[string]string, error) {
	if !strings.Contains(val, ":") {
		return nil, ErrInvalidValidator{""}
	}
	validateCandidates := strings.Split(val, "|")
	extractedMap := make(map[string]string, len(validateCandidates))
	for _, candidate := range validateCandidates {
		keyWithVal := strings.Split(candidate, ":")
		if len(keyWithVal) != 2 {
			return nil, ErrInvalidValidator{message: fmt.Sprintf("validator should be func:limit but got %s", val)}
		}
		if _, ok := extractedMap[strings.Trim(keyWithVal[0], " ")]; ok {
			return nil, ErrInvalidValidator{message: fmt.Sprintf(
				"Duplicate key \"%s\" in the same validator: %s", strings.Trim(keyWithVal[0], " "), val)}
		}
		extractedMap[strings.Trim(keyWithVal[0], " ")] = keyWithVal[1]
	}
	return extractedMap, nil
}
