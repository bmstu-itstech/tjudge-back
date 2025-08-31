package tjudge_test

import (
	"testing"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/stretchr/testify/require"
)

func TestNewUserPrototype(t *testing.T) {
	type testCase struct {
		testname    string
		username    string
		fullname    string
		password    string
		wantErr     bool
		expectedErr string
	}
	tests := []testCase{
		{
			testname: "Valid user",
			username:  "anna",
			fullname:  "Anna Ivanova",
			password: "12345678",
			wantErr:  false,
		},
		{
			testname: "Valid user with long password",
			username: "bob",
			fullname: "Bob Boska",
			password: "verylongpassword123!@#",
			wantErr:  false,
		},
		{
			testname:    "Empty username",
			username:     "",
			fullname:    "No No",
			password:    "password123",
			wantErr:     true,
			expectedErr: "invalid input: expected not empty name",
		},
		{
			testname:    "Short password",
			username:    "Sasha",
			fullname:    "Alex Titov",
			password:    "1234567",
			wantErr:     true,
			expectedErr: "the password contains less than 8 characters",
		},
		{
			testname:    "Empty password",
			username:    "rick_sun",
			fullname:    "Rick Rock",
			password:    "",
			wantErr:     true,
			expectedErr: "the password contains less than 8 characters",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			got, err := tjudge.NewUserPrototype(tt.username, tt.fullname, tt.password, true)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.username, got.Username())
				require.Equal(t, tt.fullname, got.Fullname())
			}
		})
	}
}
