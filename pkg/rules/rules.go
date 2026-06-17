package rules

import (
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

func CheckSpecialChars(msg string) (string, string) {
	var fixed strings.Builder
	foundSpecial := false

	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			fixed.WriteRune(r)
		} else {
			foundSpecial = true
		}
	}

	if foundSpecial {
		result := strings.Join(strings.Fields(fixed.String()), " ")
		return "contains special characters", result
	}

	return "", ""
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

		lowerPattern := strings.ToLower(cleanPattern)
		lowerFixed := strings.ToLower(fixed)
		if !strings.Contains(lowerFixed, lowerPattern) {
			continue
		}

		found = true
		fixed = redactAll(fixed, lowerFixed, lowerPattern)
	}

	if found {
		return "may contain sensitive data", fixed
	}
	return "", ""
}

func redactAll(s, lowerS, lowerPattern string) string {
	var b strings.Builder
	last := 0
	patternLen := len(lowerPattern)

	for i := 0; i <= len(lowerS)-patternLen; i++ {
		if lowerS[i:i+patternLen] == lowerPattern {
			b.WriteString(s[last:i])
			b.WriteString("[redacted]")
			last = i + patternLen
			i += patternLen - 1
		}
	}

	b.WriteString(s[last:])
	return b.String()
}
