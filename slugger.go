package gong

import (
	"regexp"
	"strings"
)

func SlugifyTitle(ticketTitle string) string {
	re := regexp.MustCompile("[^a-z0-9]+")

	return strings.Trim(re.ReplaceAllString(strings.ToLower(ticketTitle), "-"), "-")
}
