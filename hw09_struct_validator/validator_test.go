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
		wantErr     bool
		expectedErr error
	}{
		{
			in: User{
				ID:     "12345",
				Name:   "Dart Weider",
				Age:    51,
				Email:  "eniken-empire.loc",
				Role:   "amin",
				Phones: []string{"012345678911", "01234"},
				meta:   []byte{},
			},
			wantErr: true,
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
					Field: "Email",
					Err:   valuerror.ErrValidateFieldByRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   valuerror.ErrValidateIn{TrueLimit: "admin,stuff", ActualValue: "amin"},
				},
				ValidationError{
					Field: "Phones",
					Err:   valuerror.ErrValidateLen{TrueLimit: 11, ActualValue: 12},
				},
			},
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Name:   "Dart Weider",
				Age:    50,
				Email:  "eniken@empire.loc",
				Role:   "admin",
				Phones: []string{"01234567891", "12312312311"},
				meta:   []byte{},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Name:   "Dart Weider",
				Age:    16,
				Email:  "eniken@empire.loc",
				Role:   "admin",
				Phones: []string{"01234567891", "12312312311"},
				meta:   []byte{},
			},
			wantErr: true,
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   valuerror.ErrValidateMin{TrueLimit: 18, ActualValue: 16},
				},
			},
		},
		{
			in: App{
				Version: "v.1.2.0",
			},
			wantErr: true,
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   valuerror.ErrValidateLen{TrueLimit: 5, ActualValue: 7},
				},
			},
		},
		{
			in: App{
				Version: "v.1.2",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("some header"),
				Payload:   []byte("some payload"),
				Signature: []byte("some signature"),
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ValidationErrors
			actualErr := Validate(tt.in)
			switch {
			case tt.wantErr:
				require.Error(t, actualErr)
				require.ErrorAs(t, actualErr, &valErr)
				fmt.Println(valErr)
				require.Equal(t, tt.expectedErr, actualErr)
			default:
				require.NoError(t, actualErr)
			}
		})
	}
}
