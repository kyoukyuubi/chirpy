package auth

import "testing"

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