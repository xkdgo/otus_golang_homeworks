package hw09structvalidator

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrValidateLen           = errors.New("invalid field len")
	ErrValidateMax           = errors.New("greater than max")
	ErrValidateFieldByRegexp = errors.New("doesn`t match regular expression")
)

type ErrValidateLen struct {
	trueLen   int
	actualLen int
}

func (v ErrValidateLen) Error() string {
	b := strings.Builder{}
	b.WriteString("invalid actual field len is ")
	b.WriteString(strconv.Itoa(v.trueLen))
	b.WriteString("should be equal ")
	b.WriteString(strconv.Itoa(v.actualLen))
	return b.String()
}

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
	// Place your code here.
	return errs
}
