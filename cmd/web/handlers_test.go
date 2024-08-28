package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/joshdcuneo/go-ui/internal/assert"
	"github.com/joshdcuneo/go-ui/internal/models/mocks"
)

// TODO figure out how to test authentication works
func TestHealthcheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/healthcheck")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCsrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = mocks.ValidUserPassword
		validEmail    = mocks.ValidUserEmail
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name          string
		userName      string
		userEmail     string
		userPassword  string
		csrfToken     string
		expectCode    int
		expectFormTag string
	}{
		{
			name:          "Valid submission",
			userName:      validName,
			userEmail:     validEmail,
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusSeeOther,
			expectFormTag: "",
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			expectCode:   http.StatusBadRequest,
		},
		{
			name:          "Empty name",
			userName:      "",
			userEmail:     validEmail,
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty email",
			userName:      validName,
			userEmail:     "",
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty password",
			userName:      validName,
			userEmail:     validEmail,
			userPassword:  "",
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Short password",
			userName:      validName,
			userEmail:     validEmail,
			userPassword:  "short",
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Duplicate email",
			userName:      validName,
			userEmail:     "dupe@example.com",
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.expectCode)

			if tt.expectFormTag != "" {
				assert.StringContains(t, body, tt.expectFormTag)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCsrfToken := extractCSRFToken(t, body)

	const (
		validEmail    = mocks.ValidUserEmail
		validPassword = mocks.ValidUserPassword
		formTag       = "<form action='/user/login' method='POST' novalidate>"
	)

	tests := []struct {
		name          string
		userEmail     string
		userPassword  string
		csrfToken     string
		expectCode    int
		expectFormTag string
	}{
		{
			name:          "Valid submission",
			userEmail:     validEmail,
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusSeeOther,
			expectFormTag: "",
		},
		{
			name:          "Invalid email",
			userEmail:     "invalid@example.com",
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Invalid password",
			userEmail:     validEmail,
			userPassword:  "invalid",
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty email",
			userEmail:     "",
			userPassword:  validPassword,
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty password",
			userEmail:     validEmail,
			userPassword:  "",
			csrfToken:     validCsrfToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:         "Invalid CSRF token",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			expectCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/login", form)

			assert.Equal(t, code, tt.expectCode)

			if tt.expectFormTag != "" {
				assert.StringContains(t, body, tt.expectFormTag)
			}

		})
	}
}

func TestUserLogout(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, authenticateAs(app, 1))
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCsrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", validCsrfToken)

	code, header, _ := ts.postForm(t, "/user/logout", form)
	assert.Equal(t, code, http.StatusSeeOther)
	assert.Equal(t, header.Get("Location"), "/")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/")
	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, "Go UI")
	assert.StringContains(t, body, "Login")
	assert.StringContains(t, body, "Signup")
	assert.StringNotContains(t, body, "Logout")

	ts = newTestServer(t, authenticateAs(app, 1))
	defer ts.Close()

	code, _, body = ts.get(t, "/")
	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, "Go UI")
	assert.StringNotContains(t, body, "Login")
	assert.StringNotContains(t, body, "Signup")
	assert.StringContains(t, body, "Logout")
}
