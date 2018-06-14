package gong

import (
	"regexp"
	"strings"
)

// SlugifyTitle : Make sure to clean up the branch names
func SlugifyTitle(ticketTitle string, replacementCharacter string) string {
	re := regexp.MustCompile("[^a-z0-9]+");
	lowerCaseTitle := strings.ToLower(ticketTitle)
	formattedTitle := re.ReplaceAllString(lowerCaseTitle, replacementCharacter)
	return strings.Trim(formattedTitle, replacementCharacter)
}
