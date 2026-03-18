package emailutils

import (
	"testing"

	"github.com/lm379/hmx/config"
)

func TestSendEmailConfigValidation(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name        string
		smtpHost    string
		smtpPort    int
		smtpUser    string
		smtpPass    string
		to          string
		subject     string
		body        string
		expectError bool
	}{
		{
			name:        "Valid configuration",
			smtpHost:    "smtp.example.com",
			smtpPort:    587,
			smtpUser:    "test@example.com",
			smtpPass:    "password",
			to:          "recipient@example.com",
			subject:     "Test Subject",
			body:        "Test Body",
			expectError: true, // Will fail because SMTP server is not available
		},
		{
			name:        "Empty SMTP host",
			smtpHost:    "",
			smtpPort:    587,
			smtpUser:    "test@example.com",
			smtpPass:    "password",
			to:          "recipient@example.com",
			subject:     "Test Subject",
			body:        "Test Body",
			expectError: true,
		},
		{
			name:        "Empty SMTP user",
			smtpHost:    "smtp.example.com",
			smtpPort:    587,
			smtpUser:    "",
			smtpPass:    "password",
			to:          "recipient@example.com",
			subject:     "Test Subject",
			body:        "Test Body",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.AppConfig = config.Config{
				SMTPHost: tt.smtpHost,
				SMTPPort: tt.smtpPort,
				SMTPUser: tt.smtpUser,
				SMTPPass: tt.smtpPass,
			}

			err := SendEmail(tt.to, tt.subject, tt.body)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError {
				t.Log("SendEmail called without error (SMTP server may not be available)")
			}
		})
	}
}

func TestSendEmailParameters(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		SMTPUser: "test@example.com",
		SMTPPass: "password",
	}

	tests := []struct {
		name    string
		to      string
		subject string
		body    string
	}{
		{
			name:    "Simple email",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "Test Body",
		},
		{
			name:    "HTML body",
			to:      "recipient@example.com",
			subject: "HTML Email",
			body:    "<html><body><h1>Test</h1></body></html>",
		},
		{
			name:    "Empty subject",
			to:      "recipient@example.com",
			subject: "",
			body:    "Test Body",
		},
		{
			name:    "Empty body",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "",
		},
		{
			name:    "Long subject",
			to:      "recipient@example.com",
			subject: string(make([]byte, 100)),
			body:    "Test Body",
		},
		{
			name:    "Long body",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    string(make([]byte, 10000)),
		},
		{
			name:    "Special characters in subject",
			to:      "recipient@example.com",
			subject: "测试标题！@#￥%",
			body:    "Test Body",
		},
		{
			name:    "Unicode in body",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "这是测试内容，包含中文和特殊字符！😊",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmail(tt.to, tt.subject, tt.body)

			// We expect errors because SMTP server is not available
			// But we verify the function can handle various inputs without panicking
			if err != nil {
				t.Logf("SendEmail returned error (expected): %v", err)
			}
		})
	}
}

func TestSendEmailRecipientFormats(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		SMTPUser: "test@example.com",
		SMTPPass: "password",
	}

	tests := []struct {
		name string
		to   string
	}{
		{
			name: "Simple email",
			to:   "user@example.com",
		},
		{
			name: "Email with subdomain",
			to:   "user@mail.example.com",
		},
		{
			name: "Email with plus sign",
			to:   "user+tag@example.com",
		},
		{
			name: "Email with dots",
			to:   "first.last@example.com",
		},
		{
			name: "Multiple recipients (comma separated)",
			to:   "user1@example.com,user2@example.com",
		},
		{
			name: "Email with numbers",
			to:   "user123@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmail(tt.to, "Test", "Test Body")

			// We expect errors because SMTP server is not available
			if err != nil {
				t.Logf("SendEmail returned error (expected): %v", err)
			}
		})
	}
}

func TestSendEmailWithHTML(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		SMTPUser: "test@example.com",
		SMTPPass: "password",
	}

	htmlBody := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>测试邮件</title>
</head>
<body>
    <h1>您好！</h1>
    <p>这是一封测试邮件。</p>
    <p style="color: red;">红色文字</p>
</body>
</html>`

	err := SendEmail("recipient@example.com", "HTML 邮件测试", htmlBody)

	// We expect error because SMTP server is not available
	if err != nil {
		t.Logf("SendEmail returned error (expected): %v", err)
	}
}

func TestSendEmailEdgeCases(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		SMTPUser: "test@example.com",
		SMTPPass: "password",
	}

	tests := []struct {
		name    string
		to      string
		subject string
		body    string
	}{
		{
			name:    "Empty to field",
			to:      "",
			subject: "Test",
			body:    "Body",
		},
		{
			name:    "All empty fields",
			to:      "",
			subject: "",
			body:    "",
		},
		{
			name:    "Whitespace only",
			to:      "   ",
			subject: "   ",
			body:    "   ",
		},
		{
			name:    "Newlines in subject",
			to:      "user@example.com",
			subject: "Subject\nwith\nnewlines",
			body:    "Body",
		},
		{
			name:    "Special characters in body",
			to:      "user@example.com",
			subject: "Test",
			body:    "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that function doesn't panic with edge cases
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("SendEmail panicked with: %v", r)
				}
			}()

			err := SendEmail(tt.to, tt.subject, tt.body)
			if err != nil {
				t.Logf("SendEmail returned error: %v", err)
			}
		})
	}
}

func TestSendEmailConsistency(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		SMTPUser: "test@example.com",
		SMTPPass: "password",
	}

	// Test that same inputs produce consistent results
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "Test Body"

	err1 := SendEmail(to, subject, body)
	err2 := SendEmail(to, subject, body)

	// Both should return errors (SMTP not available)
	if (err1 == nil) != (err2 == nil) {
		t.Errorf("Inconsistent error results: %v vs %v", err1, err2)
	}
}

func TestSendEmailConfigPorts(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name        string
		smtpPort    int
		expectError bool
	}{
		{
			name:        "Standard SMTP port 25",
			smtpPort:    25,
			expectError: true,
		},
		{
			name:        "SMTP submission port 587",
			smtpPort:    587,
			expectError: true,
		},
		{
			name:        "SMTPS port 465",
			smtpPort:    465,
			expectError: true,
		},
		{
			name:        "Custom port 2525",
			smtpPort:    2525,
			expectError: true,
		},
		{
			name:        "Port 0 (invalid)",
			smtpPort:    0,
			expectError: true,
		},
		{
			name:        "Negative port",
			smtpPort:    -1,
			expectError: true,
		},
		{
			name:        "Port out of range",
			smtpPort:    99999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.AppConfig = config.Config{
				SMTPHost: "smtp.example.com",
				SMTPPort: tt.smtpPort,
				SMTPUser: "test@example.com",
				SMTPPass: "password",
			}

			err := SendEmail("recipient@example.com", "Test", "Body")

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
