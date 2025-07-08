package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
	}{
		{name: "simple password", password: "password"},
		{name: "numeric password", password: "123456"},
		{name: "qwerty password", password: "qwerty"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tc.password)
			if err != nil {
				t.Errorf("HashPassword returned error: %v", err)
			}
			if hashedPassword == tc.password {
				t.Errorf("hashed password should not equal original password")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		name     string
		password string
	}{
		{name: "simple password", password: "password"},
		{name: "numeric password", password: "123456"},
		{name: "qwerty password", password: "qwerty"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tc.password)
			if err != nil {
				t.Errorf("HashPassword returned error: %v", err)
			}
			if hashedPassword == tc.password {
				t.Errorf("hashed password should not equal original password")
			}
			err = CheckPasswordHash(tc.password, hashedPassword)
			if err != nil {
				t.Errorf("VerifyPassword returned error: %v", err)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	cases := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "valid token",
			userID:      uuid.New(),
			tokenSecret: "mysecretkey",
			expiresIn:   1 * time.Hour,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "empty secret",
			userID:      uuid.New(),
			tokenSecret: "",
			expiresIn:   1 * time.Hour,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "expired token time",
			userID:      uuid.New(),
			tokenSecret: "mysecretkey",
			expiresIn:   -1 * time.Hour,
			expectError: false,
			expectEmpty: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := MakeJWT(tc.userID, tc.tokenSecret, tc.expiresIn)

			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tc.expectEmpty && token != "" {
				t.Errorf("expected empty token but got %q", token)
			}
			if !tc.expectEmpty && token == "" {
				t.Errorf("expected non-empty token but got empty")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()

	validToken, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to create valid token: %v", err)
	}

	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	invalidToken := validToken + "invalid"

	cases := []struct {
		name        string
		tokenString string
		secret      string
		expectError bool
		expectUser  uuid.UUID
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			secret:      secret,
			expectError: false,
			expectUser:  userID,
		},
		{
			name:        "expired token",
			tokenString: expiredToken,
			secret:      secret,
			expectError: true,
		},
		{
			name:        "invalid token string",
			tokenString: invalidToken,
			secret:      secret,
			expectError: true,
		},
		{
			name:        "wrong secret",
			tokenString: validToken,
			secret:      "wrongsecret",
			expectError: true,
		},
		{
			name:        "empty token",
			tokenString: "",
			secret:      secret,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := ValidateJWT(tc.tokenString, tc.secret)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if user != tc.expectUser {
				t.Errorf("expected userID %v, got %v", tc.expectUser, user)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name        string
		headers     http.Header
		expectError bool
		expectToken string
	}{
		{
			name:        "valid token",
			headers:     http.Header{"Authorization": {"Bearer token"}},
			expectError: false,
			expectToken: "token",
		},
		{
			name:        "missing token",
			headers:     http.Header{"Authorization": {}},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.headers)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if token != tc.expectToken {
				t.Errorf("expected token %q, got %q", tc.expectToken, token)
			}
		})
	}
}
