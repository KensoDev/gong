package gong

import (
	"fmt"
	"bytes"
	"text/template"
	"strings"
)

type Branch struct {
	pattern              string
	replacementCharacter string
}

type BranchName struct {
	IssueType  string
	IssueID    string
	IssueTitle string
}

var issueTypeTemplate = "{{.IssueType}}"
var issueIDTemplate = "{{.IssueID}}"
var issueTitleTemplate = "{{.IssueTitle}}"
var defaultPattern = fmt.Sprintf(`%s/%s-%s`, issueTypeTemplate, issueIDTemplate, issueTitleTemplate)

func NewBranch(configuration map[string]string) (Branch, error) {
	pattern := fieldOrDefault(configuration, "branch_pattern", defaultPattern)
	if err := validatePattern(pattern); err != nil {
		return Branch{}, err
	}
	replacementCharacter := fieldOrDefault(configuration, "branch_replacement_character", "-")
	return Branch{pattern, replacementCharacter}, nil
}

func fieldOrDefault(configuration map[string]string, field string, defaultValue string) string {
	value := configuration[field]
	if value == "" {
		return defaultValue
	}
	return value
}

func validatePattern(pattern string) error {
	if strings.Contains(pattern, issueIDTemplate) && strings.Contains(pattern, issueTitleTemplate) {
		return nil
	}
	return fmt.Errorf("branch_pattern should have both %s and %s", issueIDTemplate, issueTitleTemplate)
}

func (b *Branch) Name(issueType string, issueID string, issueTitle string) (string, error) {
	branchTemplate := template.New("template")
	branchTemplate, err := branchTemplate.Parse(b.pattern)
	if err != nil {
		return "", err
	}
	formattedTitle := SlugifyTitle(issueTitle, b.replacementCharacter)
	branchName := BranchName{issueType, issueID, formattedTitle}
	var result bytes.Buffer
	if err := branchTemplate.Execute(&result, branchName); err != nil {
		return "", err
	}	
	fmt.Printf(result.String())
	return result.String(), nil
}
