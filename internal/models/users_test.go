package models

import (
	"testing"

	"github.com/joshdcuneo/go-ui/internal/assert"
)

const (
	validName     = "Test User"
	validEmail    = "test@example.com"
	validPassword = "password"
)

func TestUserModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := newTestDB(t)
	user := UserModel{DB: db}

	err := user.Insert(validName, validEmail, validPassword)
	assert.NilError(t, err)

	err = user.Insert(validName, validEmail, validPassword)
	assert.Equal(t, err, ErrDuplicateEmail)
}

func TestUserModelAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		password  string
		expect    int
		expectErr error
	}{
		{
			name:      "Valid credentials",
			email:     validEmail,
			password:  validPassword,
			expect:    1,
			expectErr: nil,
		},
		{
			name:      "Invalid email",
			email:     "invalid@example.com",
			password:  validPassword,
			expect:    0,
			expectErr: ErrInvalidCredentials,
		},
		{
			name:      "Invalid password",
			email:     validEmail,
			password:  "wrongpassword",
			expect:    0,
			expectErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := newTestDB(t)
			user := UserModel{DB: db}
			user.Insert(validName, validEmail, validPassword)

			actual, err := user.Authenticate(tt.email, tt.password)
			assert.Equal(t, actual, tt.expect)
			assert.Equal(t, err, tt.expectErr)
		})
	}
}

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		expect bool
	}{
		{
			name:   "User exists",
			userID: 1,
			expect: true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			expect: false,
		},
		{
			name:   "Non-existent ID",
			userID: 999,
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			user := UserModel{DB: db}
			user.Insert(validName, validEmail, validPassword)
			actual, err := user.Exists(tt.userID)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, actual, tt.expect)
			assert.NilError(t, err)
		})
	}
}
