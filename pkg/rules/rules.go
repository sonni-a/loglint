package rules

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func CheckLowercase(msg string) (string, string) {
	if msg == "" {
		return "", ""
	}
	r, size := utf8.DecodeRuneInString(msg)
	if unicode.IsUpper(r) {
		fixed := string(unicode.ToLower(r)) + msg[size:]
		return "must start with lowercase letter", fixed
	}
	return "", ""
}

func isStrippable(r rune) bool {
	switch r {
	case '!', '?':
		return true
	default:
		return unicode.Is(unicode.So, r)
	}
}

func CheckSpecialChars(msg string) (string, string) {
	var b strings.Builder
	found := false

	for _, r := range msg {
		if isStrippable(r) {
			found = true
			continue
		}
		b.WriteRune(r)
	}

	if !found {
		return "", ""
	}

	fixed := strings.Join(strings.Fields(b.String()), " ")
	if fixed == "" || fixed == msg {
		return "contains special characters", msg
	}
	return "contains special characters", fixed
}

func CheckEnglish(msg string) (string, string) {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return "must be in English", msg
		}
	}
	return "", ""
}

func CheckSensitive(msg string, userPatterns []string) (string, string) {
	fixed := msg
	found := false

	for _, pattern := range userPatterns {
		cleanPattern := strings.TrimSpace(pattern)
		if cleanPattern == "" {
			continue
		}

		next := redactAll(fixed, cleanPattern)
		if next == fixed {
			continue
		}

		found = true
		fixed = next
	}

	if found {
		return "may contain sensitive data", fixed
	}
	return "", ""
}

func redactAll(s, pattern string) string {
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(pattern))
	return re.ReplaceAllString(s, "redacted")
}
