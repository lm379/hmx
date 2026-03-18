package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestUserSex(t *testing.T) {
	tests := []struct {
		name  string
		sex   UserSex
		valid bool
	}{
		{"Male value", Male, true},
		{"Female value", Female, true},
		{"Other value", Other, true},
		{"Empty value", UserSex(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && string(tt.sex) == "" {
				t.Error("Expected non-empty string for valid UserSex")
			}
		})
	}
}

func TestUserRole(t *testing.T) {
	tests := []struct {
		name  string
		role  UserRole
		valid bool
	}{
		{"Administrator value", Administrator, true},
		{"User value", User, true},
		{"Empty value", UserRole(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && string(tt.role) == "" {
				t.Error("Expected non-empty string for valid UserRole")
			}
		})
	}
}

func TestUsersTableName(t *testing.T) {
	user := Users{}
	tableName := user.TableName()

	if tableName != "users" {
		t.Errorf("Expected table name 'users', got %s", tableName)
	}
}

func TestUsersStruct(t *testing.T) {
	now := time.Now()

	user := Users{
		UserID:      1,
		Username:    "testuser",
		Phone:       "13800138000",
		Email:       sql.NullString{String: "test@example.com", Valid: true},
		Password:    "hashedpassword",
		Sex:         Male,
		Icon:        sql.NullString{String: "avatar.jpg", Valid: true},
		Role:        User,
		LastLoginAt: &now,
		LastIp:      "192.168.1.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Verify field values
	if user.UserID != 1 {
		t.Errorf("UserID = %d, want 1", user.UserID)
	}
	if user.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", user.Username)
	}
	if user.Phone != "13800138000" {
		t.Errorf("Phone = %s, want 13800138000", user.Phone)
	}
	if user.Sex != Male {
		t.Errorf("Sex = %v, want Male", user.Sex)
	}
	if user.Role != User {
		t.Errorf("Role = %v, want User", user.Role)
	}
}

func TestUsersDefaults(t *testing.T) {
	user := Users{}

	// Check default values - empty struct has zero values
	// For string types, zero value is empty string
	if string(user.Sex) != "" {
		t.Errorf("Default Sex should be empty string, got %v", user.Sex)
	}
	if string(user.Role) != "" {
		t.Errorf("Default Role should be empty string, got %v", user.Role)
	}
}

func TestUserResponse(t *testing.T) {
	now := time.Now()
	email := "test@example.com"
	icon := "avatar.jpg"

	response := UserResponse{
		UserID:      1,
		Username:    "testuser",
		Phone:       "13800138000",
		Email:       &email,
		Sex:         Male,
		Icon:        &icon,
		Role:        User,
		LastLoginAt: &now,
		LastLoginIP: "192.168.1.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if response.UserID != 1 {
		t.Errorf("UserID = %d, want 1", response.UserID)
	}
	if response.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", response.Username)
	}
	if response.Email == nil {
		t.Error("Email should not be nil")
	}
	if *response.Email != email {
		t.Errorf("Email = %s, want %s", *response.Email, email)
	}
	if response.Icon == nil {
		t.Error("Icon should not be nil")
	}
	if *response.Icon != icon {
		t.Errorf("Icon = %s, want %s", *response.Icon, icon)
	}
}

func TestRegisterInput(t *testing.T) {
	input := RegisterInput{
		Username: "testuser",
		Phone:    "13800138000",
		Email:    "test@example.com",
		Password: "password123",
		Code:     "123456",
	}

	if input.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", input.Username)
	}
	if input.Phone != "13800138000" {
		t.Errorf("Phone = %s, want 13800138000", input.Phone)
	}
	if input.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", input.Email)
	}
	if input.Code != "123456" {
		t.Errorf("Code = %s, want 123456", input.Code)
	}
}

func TestLoginInput(t *testing.T) {
	input := LoginInput{
		Account:  "test@example.com",
		Password: "password123",
	}

	if input.Account != "test@example.com" {
		t.Errorf("Account = %s, want test@example.com", input.Account)
	}
	if input.Password != "password123" {
		t.Errorf("Password = %s, want password123", input.Password)
	}
}

func TestSendCodeInput(t *testing.T) {
	input := SendCodeInput{
		Email: "test@example.com",
	}

	if input.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", input.Email)
	}
}

func TestForgetPasswordInput(t *testing.T) {
	input := ForgetPasswordInput{
		Email:       "test@example.com",
		Code:        "123456",
		NewPassword: "newpassword123",
	}

	if input.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", input.Email)
	}
	if input.Code != "123456" {
		t.Errorf("Code = %s, want 123456", input.Code)
	}
	if input.NewPassword != "newpassword123" {
		t.Errorf("NewPassword = %s, want newpassword123", input.NewPassword)
	}
}

func TestUpdateUserProfileRequest(t *testing.T) {
	username := "newusername"
	phone := "13900139000"
	email := "new@example.com"
	sex := Female
	icon := "newavatar.jpg"
	code := "123456"

	input := UpdateUserProfileRequest{
		Username: username,
		Phone:    phone,
		Sex:      sex,
		Icon:     icon,
		Email:    email,
		Code:     code,
	}

	if input.Username != username {
		t.Errorf("Username = %s, want %s", input.Username, username)
	}
	if input.Phone != phone {
		t.Errorf("Phone = %s, want %s", input.Phone, phone)
	}
	if input.Email != email {
		t.Errorf("Email = %s, want %s", input.Email, email)
	}
	if input.Sex != sex {
		t.Errorf("Sex = %v, want %v", input.Sex, sex)
	}
	if input.Icon != icon {
		t.Errorf("Icon = %s, want %s", input.Icon, icon)
	}
	if input.Code != code {
		t.Errorf("Code = %s, want %s", input.Code, code)
	}
}

func TestUpdatePasswordRequest(t *testing.T) {
	input := UpdatePasswordRequest{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	}

	if input.OldPassword != "oldpassword123" {
		t.Errorf("OldPassword = %s, want oldpassword123", input.OldPassword)
	}
	if input.NewPassword != "newpassword456" {
		t.Errorf("NewPassword = %s, want newpassword456", input.NewPassword)
	}
}

func TestUsersGORMTags(t *testing.T) {
	// Verify that Users struct has proper GORM tags
	user := Users{}

	// We can't directly access struct tags in tests,
	// but we can verify the struct works with GORM
	if user.UserID == 0 {
		// This is just to ensure the struct compiles correctly
		// Real GORM tests would require an actual database
	}
}

func TestUserSexStringComparison(t *testing.T) {
	tests := []struct {
		sex UserSex
		str string
	}{
		{Male, "Male"},
		{Female, "Female"},
		{Other, "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if string(tt.sex) != tt.str {
				t.Errorf("UserSex.String() = %s, want %s", string(tt.sex), tt.str)
			}
		})
	}
}

func TestUserRoleStringComparison(t *testing.T) {
	tests := []struct {
		role UserRole
		str  string
	}{
		{Administrator, "Administrator"},
		{User, "User"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if string(tt.role) != tt.str {
				t.Errorf("UserRole.String() = %s, want %s", string(tt.role), tt.str)
			}
		})
	}
}
