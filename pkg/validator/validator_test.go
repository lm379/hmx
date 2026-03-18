package validator

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		// Valid emails
		{
			name:  "Valid simple email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "Valid email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "Valid email with plus sign",
			email: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "Valid email with dots",
			email: "first.last@example.com",
			want:  true,
		},
		{
			name:  "Valid email with numbers",
			email: "user123@example.com",
			want:  true,
		},
		{
			name:  "Valid email with hyphen",
			email: "user-name@example.com",
			want:  true,
		},
		{
			name:  "Valid email with underscore",
			email: "user_name@example.com",
			want:  true,
		},
		{
			name:  "Valid email with percentage",
			email: "user%name@example.com",
			want:  true,
		},

		// Invalid emails
		{
			name:  "Empty email",
			email: "",
			want:  false,
		},
		{
			name:  "No @ symbol",
			email: "invalidemail.com",
			want:  false,
		},
		{
			name:  "No domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "No local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "No TLD",
			email: "user@example",
			want:  false,
		},
		{
			name:  "Multiple @ symbols",
			email: "user@@example.com",
			want:  false,
		},
		{
			name:  "Invalid TLD",
			email: "user@example.c",
			want:  false,
		},
		{
			name:  "Spaces in email",
			email: "user @example.com",
			want:  false,
		},
		{
			name:  "Special characters at start",
			email: ".user@example.com",
			want:  true,
		},
		{
			name:  "Special characters at end",
			email: "user.@example.com",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		// Valid Chinese phone numbers
		{
			name:  "Valid phone number with 13",
			phone: "13123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 14",
			phone: "14123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 15",
			phone: "15123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 16",
			phone: "16123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 17",
			phone: "17123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 18",
			phone: "18123456789",
			want:  true,
		},
		{
			name:  "Valid phone number with 19",
			phone: "19123456789",
			want:  true,
		},

		// Invalid phone numbers
		{
			name:  "Empty phone",
			phone: "",
			want:  false,
		},
		{
			name:  "Not starting with 1",
			phone: "23123456789",
			want:  false,
		},
		{
			name:  "Second digit is 2",
			phone: "12123456789",
			want:  false,
		},
		{
			name:  "Second digit is 0",
			phone: "10123456789",
			want:  false,
		},
		{
			name:  "Too short",
			phone: "1312345678",
			want:  false,
		},
		{
			name:  "Too long",
			phone: "131234567890",
			want:  false,
		},
		{
			name:  "Contains letters",
			phone: "1312345678a",
			want:  false,
		},
		{
			name:  "Contains special characters",
			phone: "1312345678!",
			want:  false,
		},
		{
			name:  "With spaces",
			phone: "131 2345 6789",
			want:  false,
		},
		{
			name:  "With hyphen",
			phone: "131-2345-6789",
			want:  false,
		},
		{
			name:  "With parentheses",
			phone: "(131)23456789",
			want:  false,
		},
		{
			name:  "All zeros",
			phone: "10000000000",
			want:  false,
		},
		{
			name:  "Only 9 digits",
			phone: "123456789",
			want:  false,
		},
		{
			name:  "Negative number",
			phone: "-13123456789",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePhone(tt.phone); got != tt.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestValidateEmailEdgeCases(t *testing.T) {
	// Test edge cases for email validation
	edgeCases := []struct {
		email string
		want  bool
	}{
		{"a@b.cd", true},
		{"a@b.c", false},
		{"a@b.cde", true},
		{"a@b.cdef", true},
		{"1@2.34", false},
		{"test@domain.co.uk", true},
		{"test@domain.info", true},
		{"test@domain.travel", true},
		{"test@domain.museum", true},
		{"test@domain.aero", true},
	}

	for _, tc := range edgeCases {
		t.Run(tc.email, func(t *testing.T) {
			if got := ValidateEmail(tc.email); got != tc.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tc.email, got, tc.want)
			}
		})
	}
}

func TestValidatePhoneEdgeCases(t *testing.T) {
	// Test edge cases for phone validation
	edgeCases := []struct {
		phone string
		want  bool
	}{
		{"13000000000", true},
		{"19999999999", true},
		{"130000000001", false},
		{"0", false},
		{"12345", false},
		{"11111111111", false},
		{"18999999999", true},
		{"10000000001", false},
	}

	for _, tc := range edgeCases {
		t.Run(tc.phone, func(t *testing.T) {
			if got := ValidatePhone(tc.phone); got != tc.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tc.phone, got, tc.want)
			}
		})
	}
}
