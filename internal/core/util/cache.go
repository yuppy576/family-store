package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// GenerateCacheKey generates a cache key based on the input parameters
func GenerateCacheKey(prefix string, params any) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}

// GenerateCacheParams generates a cache params based on the input parameters
func GenerateCacheKeyParams(params ...any) string {
	var str string

	for i, param := range params {
		str += fmt.Sprintf("%v", param)

		last := len(params) - 1
		if i != last {
			str += "-"
		}
	}

	return str
}

// Serialize marshals the input data into an array of bytes
func Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

// Deserialize unmarshals the input data into the output interface
func Deserialize(data []byte, output any) error {
	return json.Unmarshal(data, output)
}

var chineseRegex = regexp.MustCompile(`[\p{Han}]+`)

func Slugify(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if chineseRegex.MatchString(s) {
		result := ""
		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				result += strings.ToLower(string(r))
			} else {
				result += "-"
			}
		}
		s = result
	} else {
		re := regexp.MustCompile(`[^a-z0-9]+`)
		s = re.ReplaceAllString(s, "-")
	}

	s = strings.Trim(s, "-")
	re := regexp.MustCompile(`-{2,}`)
	s = re.ReplaceAllString(s, "-")

	if len(s) > 30 {
		s = s[:30]
		s = strings.TrimSuffix(s, "-")
	}

	if s == "" {
		s = "store"
	}

	return s
}
