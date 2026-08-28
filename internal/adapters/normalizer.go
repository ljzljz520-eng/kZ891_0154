package adapters

import (
	"regexp"
	"strings"
)

var serialCleaner = regexp.MustCompile(`[^A-Z0-9-]`)

func CanonicalPhone(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	if strings.HasPrefix(value, "+86") {
		value = strings.TrimPrefix(value, "+86")
	}
	return value
}

func CanonicalSerial(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	return serialCleaner.ReplaceAllString(value, "")
}

func IsLikelyPhone(value string) bool {
	value = CanonicalPhone(value)
	return len(value) == 11 && value[0] == '1'
}
func IsLikelySerial(value string) bool {
	value = CanonicalSerial(value)
	return len(value) >= 6 && len(value) <= 32
}
func CanonicalCategory(value string) string { return strings.TrimSpace(strings.ToLower(value)) }
