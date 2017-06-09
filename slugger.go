package gong

import (
	"regexp"
	"strings"
)

// SlugifyTitle : Make sure to clean up the branch names
func SlugifyTitle(ticketTitle string) string {
	re := regexp.MustCompile("[^a-z0-9]+")

	return strings.Trim(re.ReplaceAllString(strings.ToLower(ticketTitle), "-"), "-")
}
