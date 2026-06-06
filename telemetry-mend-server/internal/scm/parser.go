package scm

import (
	"regexp"
	"strconv"
)

type CodeLocation struct {
	FilePath string
	Line     int
}

var (
	// Matches /path/to/file.go:123
	goStackTraceRegex = regexp.MustCompile(`([a-zA-Z0-9._/-]+\.go):(\d+)`)
)

func ParseStackTrace(log string) []CodeLocation {
	matches := goStackTraceRegex.FindAllStringSubmatch(log, -1)
	var locations []CodeLocation
	for _, match := range matches {
		line, _ := strconv.Atoi(match[2])
		locations = append(locations, CodeLocation{
			FilePath: match[1],
			Line:     line,
		})
	}
	return locations
}
