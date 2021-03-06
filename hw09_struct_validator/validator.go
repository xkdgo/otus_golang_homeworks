package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/xkdgo/otus_golang_homeworks/hw09_struct_validator/valuerror"
)

var validatorKey = "validate"

var ErrExtractValidator = errors.New("some extractor error")

const (
	stringType      = "string"
	stringSliceType = "[]string"
	intType         = "int"
	intSliceType    = "[]int"
)

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

type validationFunc func(tagKeyMethod string,
	validationLimit string,
	structFieldName string,
	valueToCheck reflect.Value,
	valueType reflect.Type) error

var validationFuncMap = map[string]validationFunc{
	"len":    validateLen,
	"min":    validateMinMax,
	"max":    validateMinMax,
	"in":     validateIn,
	"regexp": validateRegexp,
}

func validateLen(tagKeyMethod string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type) error {
	limitInt, err := strconv.Atoi(validationLimit)
	if err != nil {
		return ErrInvalidValidator{fmt.Sprintf("%s should be integer", validationLimit)}
	}
	switch {
	case valueType.Kind().String() == stringType:
		if len(valueToCheck.String()) != limitInt {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateLen{TrueLimit: limitInt, ActualValue: len(valueToCheck.String())},
			}
		}
	case valueType.String() == stringSliceType:
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]string)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("could not assert %#v to string slice", elemInterface)}
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
			"field %s unsupported type %s for method \"%s\"", structFieldName, valueType, tagKeyMethod)}
	}
	return nil
}

func validateMinMax(tagKeyMethod string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type) error {
	limitInt, err := strconv.Atoi(validationLimit)
	if err != nil {
		return ErrInvalidValidator{fmt.Sprintf("%s should be integer", validationLimit)}
	}
	switch {
	case valueType.Kind().String() == intType:
		switch tagKeyMethod {
		case "min":
			if int(valueToCheck.Int()) < limitInt {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateMin{TrueLimit: limitInt, ActualValue: int(valueToCheck.Int())},
				}
			}
		case "max":
			if int(valueToCheck.Int()) > limitInt {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateMax{TrueLimit: limitInt, ActualValue: int(valueToCheck.Int())},
				}
			}
		}

	case valueType.String() == intSliceType:
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]int)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("could not assert %#v to int slice", elemInterface)}
		}
		for _, elem := range elemSlice {
			switch tagKeyMethod {
			case "min":
				if elem < limitInt {
					return ValidationError{
						Field: structFieldName,
						Err:   valuerror.ErrValidateMin{TrueLimit: limitInt, ActualValue: elem},
					}
				}
			case "max":
				if elem > limitInt {
					return ValidationError{
						Field: structFieldName,
						Err:   valuerror.ErrValidateMax{TrueLimit: limitInt, ActualValue: elem},
					}
				}
			}
		}
	default:
		return ErrInvalidValidator{fmt.Sprintf(
			"field %s unsupported type %s for method \"%s\"", structFieldName, valueType, tagKeyMethod)}
	}
	return nil
}

func convSliceAtoiMap(sliceWithStrings []string) (map[int]struct{}, error) {
	resultIntMap := make(map[int]struct{}, len(sliceWithStrings))
	for _, strValue := range sliceWithStrings {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return nil, err
		}
		resultIntMap[intValue] = struct{}{}
	}
	return resultIntMap, nil
}

func convSliceStringMap(sliceWithStrings []string) map[string]struct{} {
	resultIntMap := make(map[string]struct{}, len(sliceWithStrings))
	for _, strValue := range sliceWithStrings {
		resultIntMap[strValue] = struct{}{}
	}
	return resultIntMap
}

func validateIn(tagKeyMethod string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type) error {
	limitSlice := strings.Split(validationLimit, ",")
	switch {
	case valueType.Kind().String() == intType:
		limitIntMap, err := convSliceAtoiMap(limitSlice)
		if err != nil {
			return ErrInvalidValidator{fmt.Sprintf("for %s limit %s should be integer, %s",
				valueToCheck, validationLimit, err.Error())}
		}
		if _, ok := limitIntMap[int(valueToCheck.Int())]; !ok {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateIn{TrueLimit: validationLimit, ActualValue: fmt.Sprintf("%d", valueToCheck.Int())},
			}
		}
	case valueType.String() == intSliceType:
		limitIntMap, err := convSliceAtoiMap(limitSlice)
		if err != nil {
			return ErrInvalidValidator{fmt.Sprintf("for %s limit %s should be integer, %s",
				valueToCheck, validationLimit, err.Error())}
		}
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]int)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("could not assert %#v to int slice", elemInterface)}
		}
		for _, elem := range elemSlice {
			if _, ok := limitIntMap[elem]; !ok {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateIn{TrueLimit: validationLimit, ActualValue: fmt.Sprintf("%d", elem)},
				}
			}
		}
	case valueType.Kind().String() == stringType:
		limitStringMap := convSliceStringMap(limitSlice)
		if _, ok := limitStringMap[valueToCheck.String()]; !ok {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateIn{TrueLimit: validationLimit, ActualValue: valueToCheck.String()},
			}
		}
	case valueType.String() == stringSliceType:
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]string)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("could not assert %#v to string slice", elemInterface)}
		}
		limitStringMap := convSliceStringMap(limitSlice)
		for _, elem := range elemSlice {
			if _, ok = limitStringMap[elem]; !ok {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateIn{TrueLimit: validationLimit, ActualValue: elem},
				}
			}
		}
	default:
		return ErrInvalidValidator{fmt.Sprintf(
			"field %s unsupported type %s for method \"%s\"", structFieldName, valueType, tagKeyMethod)}
	}
	return nil
}

func validateRegexp(tagKeyMethod string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type) error {
	regex, err := regexp.Compile(validationLimit)
	if err != nil {
		return ErrInvalidValidator{fmt.Sprintf("could not compile regex %#v", validationLimit)}
	}
	switch {
	case valueType.Kind().String() == stringType:
		if !regex.MatchString(valueToCheck.String()) {
			return ValidationError{
				Field: structFieldName,
				Err:   valuerror.ErrValidateFieldByRegexp,
			}
		}
	case valueType.String() == stringSliceType:
		elemInterface := valueToCheck.Interface()
		elemSlice, ok := elemInterface.([]string)
		if !ok {
			return ErrInvalidValidator{fmt.Sprintf("could not assert %#v to string slice", elemInterface)}
		}
		for _, elem := range elemSlice {
			if !regex.MatchString(elem) {
				return ValidationError{
					Field: structFieldName,
					Err:   valuerror.ErrValidateFieldByRegexp,
				}
			}
		}
	default:
		return ErrInvalidValidator{fmt.Sprintf(
			"field %s unsupported type %s for method \"%s\"", structFieldName, valueType, tagKeyMethod)}
	}
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
		if len(v) != 0 && i < len(v)-1 {
			b.WriteString(" \\\n")
		}
	}
	return b.String()
}

// This function validate fields of struct, if they have "validatorKey" flag.
// By default this validatorKey="validate".
// ...
func Validate(v interface{}) error {
	var errs ValidationErrors
	var valerr ValidationError
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Struct {
		return ErrInvalidValidator{"Validate support for structs only"}
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
		tagValue, ok := fieldTag.Lookup(validatorKey)
		if !ok {
			continue
		}
		if tagValue == "" {
			continue
		}
		extrValMap, err := extractValidators(tagValue)
		if err != nil {
			return err
		}
		for key, limit := range extrValMap {
			validationFn, ok := validationFuncMap[key]
			if !ok {
				log.Fatalf("validator %s not implemented", key)
			}
			err = validationFn(key, limit, fieldName, fieldValue, fieldType)
			if err != nil {
				if ok = errors.As(err, &valerr); !ok {
					return err
				}
				errs = append(errs, valerr)
			}
		}
	}
	if len(errs) == 0 {
		return nil
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
		keyWithVal := strings.SplitN(candidate, ":", 2)
		if len(keyWithVal) != 2 {
			return nil, ErrInvalidValidator{message: fmt.Sprintf("validator should be func:limit but got %s", val)}
		}
		if _, ok := extractedMap[strings.Trim(keyWithVal[0], " ")]; ok {
			return nil, ErrInvalidValidator{message: fmt.Sprintf(
				"duplicate key \"%s\" in the same validator: %s", strings.Trim(keyWithVal[0], " "), val)}
		}
		extractedMap[strings.Trim(keyWithVal[0], " ")] = keyWithVal[1]
	}
	return extractedMap, nil
}
