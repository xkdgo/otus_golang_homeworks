package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"sort"
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
		{
			in: Response{
				Code: 200,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
			},
			wantErr: true,
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   valuerror.ErrValidateIn{TrueLimit: "200,404,500", ActualValue: "201"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ValidationErrors
			var expectedErr ValidationErrors
			actualErr := Validate(tt.in)
			switch {
			case tt.wantErr:
				require.Error(t, actualErr)
				require.ErrorAs(t, actualErr, &valErr)
				require.ErrorAs(t, tt.expectedErr, &expectedErr)
				sort.Slice(valErr, func(i, j int) bool {
					if valErr[i].Field == valErr[j].Field {
						return valErr[i].Err.Error() < valErr[j].Err.Error()
					}
					return valErr[i].Field > valErr[j].Field
				})
				sort.Slice(expectedErr, func(i, j int) bool {
					if expectedErr[i].Field == expectedErr[j].Field {
						return expectedErr[i].Err.Error() < expectedErr[j].Err.Error()
					}
					return expectedErr[i].Field > expectedErr[j].Field
				})
				require.Equal(t, expectedErr, valErr)
			default:
				require.NoError(t, actualErr)
			}
		})
	}
}

func TestValidateSlices(t *testing.T) {
	type Slices struct {
		StringsCheckIn  []string `validate:"in:az,buki,vedi"`
		StringsCheckLen []string `validate:"len:4"`
		Emails          []string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Ages            []int    `validate:"min:10|max:16"`
		Limits          []int    `validate:"in:4,5,6"`
	}
	tests := []struct {
		in          interface{}
		wantErr     bool
		expectedErr error
	}{

		{
			in: Slices{
				StringsCheckIn:  []string{"buki", "az", "vedi", "az"},
				StringsCheckLen: []string{"1234", "4321", "abcd"},
				Emails:          []string{"good1@dot.com", "good2@dot.com", "good3@dot.com"},
				Ages:            []int{11, 12, 13, 14},
				Limits:          []int{6, 6, 4, 4, 5},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Slices{
				StringsCheckIn:  []string{"buk", "a", "ved", "az"},
				StringsCheckLen: []string{"123", "432", "abc"},
				Emails:          []string{"bad1-dot.com", "bad2-dot.com", "bad3-dot.com"},
				Ages:            []int{8, 18},
				Limits:          []int{1, 2, 3},
			},
			wantErr: true,
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "StringsCheckIn", Err: valuerror.ErrValidateIn{TrueLimit: "az,buki,vedi", ActualValue: "buk"},
				},
				ValidationError{
					Field: "StringsCheckLen", Err: valuerror.ErrValidateLen{TrueLimit: 4, ActualValue: 3},
				},
				ValidationError{
					Field: "Emails", Err: valuerror.ErrValidateFieldByRegexp,
				},
				ValidationError{
					Field: "Ages", Err: valuerror.ErrValidateMin{TrueLimit: 10, ActualValue: 8},
				},
				ValidationError{
					Field: "Ages", Err: valuerror.ErrValidateMax{TrueLimit: 16, ActualValue: 18},
				},
				ValidationError{
					Field: "Limits", Err: valuerror.ErrValidateIn{TrueLimit: "4,5,6", ActualValue: "1"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ValidationErrors
			var expectedErr ValidationErrors
			actualErr := Validate(tt.in)
			switch {
			case tt.wantErr:
				require.Error(t, actualErr)
				require.ErrorAs(t, actualErr, &valErr)
				require.ErrorAs(t, tt.expectedErr, &expectedErr)
				sort.Slice(valErr, func(i, j int) bool {
					if valErr[i].Field == valErr[j].Field {
						return valErr[i].Err.Error() < valErr[j].Err.Error()
					}
					return valErr[i].Field > valErr[j].Field
				})
				sort.Slice(expectedErr, func(i, j int) bool {
					if expectedErr[i].Field == expectedErr[j].Field {
						return expectedErr[i].Err.Error() < expectedErr[j].Err.Error()
					}
					return expectedErr[i].Field > expectedErr[j].Field
				})
				require.Equal(t, expectedErr, valErr)
			default:
				require.NoError(t, actualErr)
			}
		})
	}
}
