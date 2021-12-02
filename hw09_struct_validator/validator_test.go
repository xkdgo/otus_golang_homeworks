package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw09_struct_validator/valuerror"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:" in:admin,stuff"`
		Phones []string `validate:" len :11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "12345",
				Name:   "Dart Weider",
				Age:    51,
				Email:  "eniken-empire.loc",
				Role:   "admin",
				Phones: []string{"012345678911", "01234"},
				meta:   []byte{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   valuerror.ErrValidateLen{TrueLimit: 36, ActualValue: 5},
				},
				ValidationError{
					Field: "Age",
					Err:   valuerror.ErrValidateMax{TrueLimit: 50, ActualValue: 51},
				},
				ValidationError{
					Field: "Phones",
					Err:   valuerror.ErrValidateLen{TrueLimit: 11, ActualValue: 12},
				},
				ValidationError{
					Field: "Email",
					Err:   valuerror.ErrValidateFieldByRegexp,
				},
			},
		},
		{
			in: User{
				ID:     "12345",
				Name:   "Dart Weider",
				Age:    16,
				Email:  "eniken-empire.loc",
				Role:   "admin",
				Phones: []string{"012345678911", "01234"},
				meta:   []byte{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   valuerror.ErrValidateLen{TrueLimit: 36, ActualValue: 5},
				},
				ValidationError{
					Field: "Age",
					Err:   valuerror.ErrValidateMin{TrueLimit: 18, ActualValue: 16},
				},
				ValidationError{
					Field: "Phones",
					Err:   valuerror.ErrValidateLen{TrueLimit: 11, ActualValue: 12},
				},
				ValidationError{
					Field: "Email",
					Err:   valuerror.ErrValidateFieldByRegexp,
				},
			},
		},
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ValidationErrors
			actualErr := Validate(tt.in)
			require.Error(t, actualErr)
			require.ErrorAs(t, actualErr, &valErr)
			require.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
