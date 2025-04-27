package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHasPassword(t *testing.T) {
	cases := []struct {
		input string
		expected error
	}{
		{
			input: "password",
			expected: nil,
		},
		{
			input: "admin",
			expected: nil,
		},
		{
			input: "TestingPassword123",
			expected: nil,
		},
	}

	for _, c := range cases {
		hash, _ := HashPassword(c.input)
		if err := CheckPasswordHash(hash, c.input); err != nil {
			t.Errorf("Expected correct password to match")
		}
		if err := CheckPasswordHash(hash, "wrongpassword"); err == nil {
			t.Errorf("Expected wrong password to NOT match")
		}
	}
}

func TestJWT(t *testing.T) {
	tokenSecret1 := "SecretTokenString123"
	tokenSecret2 := "IamSecretLookAtMe"
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()

	jwtString1, _ := MakeJWT(uuid1, tokenSecret1, time.Second*30)
	jwtString2, _ := MakeJWT(uuid2, tokenSecret2, time.Second*30)
	expiredJWT, _ := MakeJWT(uuid3, tokenSecret1, -1*time.Second)

	tests := []struct {
		name string
		uuid uuid.UUID
		jwtString string
		tokenSecret string
		wantErr bool
	}{
		{
			name: "Correct jwt",
			uuid: uuid1,
			jwtString: jwtString1,
			tokenSecret: tokenSecret1,
			wantErr: false,
		},
		{
			name: "Wrong jwt",
			uuid: uuid1,
			jwtString: jwtString2,
			tokenSecret: tokenSecret1,
			wantErr: true,
		},
		{
			name: "Empty JWT",
			uuid: uuid2,
			jwtString: "",
			tokenSecret: tokenSecret2,
			wantErr: true,
		},
		{
			name: "Expired JWT",
			uuid: uuid3,
			jwtString: expiredJWT,
			tokenSecret: tokenSecret1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateJWT(tt.jwtString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && id != tt.uuid {
				t.Errorf("ValidateJWT() error = %v, want id %v", err, tt.uuid)
			}
		})
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name string
		h http.Header
		expectedToken string
		wantErr bool
	}{
		{
			name: "Correct set header",
			h: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer token123")
				return h
			}(), 
			expectedToken: "token123",
			wantErr: false,
		},
		{
			name: "Correct header 1 space too much",
			h: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer  blah")
				return h
			}(), 
			expectedToken: "blah",
			wantErr: false,
		},
		{
			name: "Invalid header, no token",
			h: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "Bearer")
				return h
			}(), 
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Invalid header, no Authorization set",
			h: func() http.Header {
				h := http.Header{}
				return h
			}(), 
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Invalid spelled header",
			h: func() http.Header {
				h := http.Header{}
				h.Set("Authorization", "InvalidBearer token123")
				return h
			}(), 
			expectedToken: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.h)
			if err != nil && !tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token != tt.expectedToken {
				t.Errorf("GetBearerToken() error = %v, token =  %v", err, tt.expectedToken)
			}
		})
	}
}