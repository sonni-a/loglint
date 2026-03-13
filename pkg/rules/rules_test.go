package rules

import (
	"testing"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMsg string
		wantFix string
	}{
		{"Already lowercase", "starting server", "", ""},
		{"Starts with upper", "Starting server", "must start with lowercase letter", "starting server"},
		{"Empty string", "", "", ""},
		{"Single upper letter", "A", "must start with lowercase letter", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotFix := CheckLowercase(tt.input)
			if gotMsg != tt.wantMsg || gotFix != tt.wantFix {
				t.Errorf("CheckLowercase() = (%q, %q), want (%q, %q)", gotMsg, gotFix, tt.wantMsg, tt.wantFix)
			}
		})
	}
}

func TestCheckSpecialChars(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMsg string
		wantFix string
	}{
		{"Clean string", "server started", "", ""},
		{"With exclamation", "error!!!", "contains special characters", "error"},
		{"With percent and space", "value %d is set", "contains special characters", "value d is set"},
		{"Emoji and dots", "done 🚀...", "contains special characters", "done"},
		{"Mixed symbols", "user_id: 123", "contains special characters", "userid 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotFix := CheckSpecialChars(tt.input)
			if gotMsg != tt.wantMsg || gotFix != tt.wantFix {
				t.Errorf("CheckSpecialChars() = (%q, %q), want (%q, %q)", gotMsg, gotFix, tt.wantMsg, tt.wantFix)
			}
		})
	}
}

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMsg string
	}{
		{"English string", "connection failed", ""},
		{"Cyrillic string", "ошибка сервера", "must be in English"},
		{"Mixed languages", "server error: ошибка", "must be in English"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, _ := CheckEnglish(tt.input)
			if gotMsg != tt.wantMsg {
				t.Errorf("CheckEnglish() = %q, want %q", gotMsg, tt.wantMsg)
			}
		})
	}
}

func TestCheckSensitive(t *testing.T) {
	patterns := []string{"password", "api_key", "token"}
	tests := []struct {
		name    string
		input   string
		wantMsg string
		wantFix string
	}{
		{"No secrets", "process started", "", ""},
		{"Contains password", "my password is 123", "may contain sensitive data", "my [redacted] is 123"},
		{"Case insensitive", "API_KEY=hidden", "may contain sensitive data", "[redacted]=hidden"},
		{"Multiple secrets", "token: abc, password: 123", "may contain sensitive data", "[redacted]: abc, [redacted]: 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotFix := CheckSensitive(tt.input, patterns)
			if gotMsg != tt.wantMsg || gotFix != tt.wantFix {
				t.Errorf("CheckSensitive() = (%q, %q), want (%q, %q)", gotMsg, gotFix, tt.wantMsg, tt.wantFix)
			}
		})
	}
}
