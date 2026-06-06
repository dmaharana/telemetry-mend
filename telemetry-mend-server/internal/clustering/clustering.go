package clustering

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

var (
	// Regex patterns to strip dynamic content
	timestampRegex = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})?`)
	uuidRegex      = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	ipv4Regex      = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	hexRegex       = regexp.MustCompile(`0x[0-9a-fA-F]+`)
	digitRegex     = regexp.MustCompile(`\b\d+\b`)
)

func SanitizeLog(log string) string {
	// Order matters: more specific to more general
	s := timestampRegex.ReplaceAllString(log, "<TIMESTAMP>")
	s = uuidRegex.ReplaceAllString(s, "<UUID>")
	s = ipv4Regex.ReplaceAllString(s, "<IP>")
	s = hexRegex.ReplaceAllString(s, "<HEX>")
	s = digitRegex.ReplaceAllString(s, "<NUM>")
	
	// Normalize whitespace
	s = strings.Join(strings.Fields(s), " ")
	
	return strings.TrimSpace(s)
}

func GenerateFingerprint(sanitizedLog string) string {
	hash := sha256.Sum256([]byte(sanitizedLog))
	return fmt.Sprintf("%x", hash)
}
