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
		if unicode.IsLetter(r) && !unicode.In(r, unicode.Latin) {
			return "must be in English", msg
		}
		if r > unicode.MaxASCII {
			return "must be in English", msg
		}
	}
	return "", ""
}

func CheckSensitive(msg string, userPatterns []string) (string, string) {
	lowerMsg := strings.ToLower(msg)
	found := false
	fixed := msg

	for _, pattern := range userPatterns {
		cleanPattern := strings.TrimSpace(pattern)
		if cleanPattern == "" {
			continue
		}

		if strings.Contains(lowerMsg, strings.ToLower(cleanPattern)) {
			found = true
			idx := strings.Index(strings.ToLower(fixed), strings.ToLower(cleanPattern))
			for idx != -1 {
				fixed = fixed[:idx] + "[redacted]" + fixed[idx+len(cleanPattern):]
				idx = strings.Index(strings.ToLower(fixed), strings.ToLower(cleanPattern))
			}
		}
	}

	if found {
		return "may contain sensitive data", fixed
	}
	return "", ""
}
