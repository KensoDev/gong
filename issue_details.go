package gong

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

func GetBranchName(jiraClient *jira.Client, issueId string, issueType string) string {
	issue, _, _ := jiraClient.Issue.Get(issueId, nil)

	issueTitleSlug := SlugifyTitle(issue.Fields.Summary)
	return fmt.Sprintf("%s/%s-%s", issueType, issueId, issueTitleSlug)
}
