package valuerror

import (
	"errors"
	"strconv"
	"strings"
)

var ErrValidateFieldByRegexp = errors.New("doesn`t match regular expression")

type ErrValidateLen struct {
	TrueLimit   int
	ActualValue int
}

func (e ErrValidateLen) Error() string {
	b := strings.Builder{}
	b.WriteString("invalid actual field len is ")
	b.WriteString(strconv.Itoa(e.ActualValue))
	b.WriteString(" should be equal ")
	b.WriteString(strconv.Itoa(e.TrueLimit))
	return b.String()
}

type ErrValidateMax struct {
	TrueLimit   int
	ActualValue int
}

func (e ErrValidateMax) Error() string {
	b := strings.Builder{}
	b.WriteString("value is ")
	b.WriteString(strconv.Itoa(e.ActualValue))
	b.WriteString(" should be less than or equal to ")
	b.WriteString(strconv.Itoa(e.TrueLimit))
	return b.String()
}

type ErrValidateMin struct {
	TrueLimit   int
	ActualValue int
}

func (e ErrValidateMin) Error() string {
	b := strings.Builder{}
	b.WriteString("value is ")
	b.WriteString(strconv.Itoa(e.ActualValue))
	b.WriteString(" should be greater than or equal to ")
	b.WriteString(strconv.Itoa(e.TrueLimit))
	return b.String()
}

type ErrValidateIn struct {
	TrueLimit   string
	ActualValue string
}

func (e ErrValidateIn) Error() string {
	b := strings.Builder{}
	b.WriteString("value is ")
	b.WriteString(e.ActualValue)
	b.WriteString(" should be one of ")
	b.WriteString(e.TrueLimit)
	return b.String()
}
