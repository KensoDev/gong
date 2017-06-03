package gong

import (
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"regexp"
)

func GetIssueID(branchName string) (string, error) {
	re := regexp.MustCompile(`([A-Z]+-[\d]+)`)
	matches := re.FindAllString(branchName, -1)

	if len(matches) == 0 {
		return "", errors.New("No matches found in the branch name")
	}

	return matches[0], nil
}

func AddComment(jiraClient *jira.Client, issueID string, commentBody string) error {
	comment := &jira.Comment{
		Body: commentBody,
	}
	_, _, err := jiraClient.Issue.AddComment(issueID, comment)

	return err
}

func GetBranchName(jiraClient *jira.Client, issueId string, issueType string) string {
	issue, _, _ := jiraClient.Issue.Get(issueId, nil)

	issueTitleSlug := SlugifyTitle(issue.Fields.Summary)
	return fmt.Sprintf("%s/%s-%s", issueType, issueId, issueTitleSlug)
}

func indexOf(status string, data []string) int {
	for k, v := range data {
		if status == v {
			return k
		}
	}
	return -1
}

func StartIssue(jiraClient *jira.Client, issueId string) error {
	allowed := []string{"Ready", "Start"}

	transitions, _, _ := jiraClient.Issue.GetTransitions(issueId)
	nextTransition := transitions[0]

	if indexOf(nextTransition.Name, allowed) > -1 {
		_, err := jiraClient.Issue.DoTransition(issueId, nextTransition.ID)

		if err != nil {
			return err
		}

		_ = StartIssue(jiraClient, issueId)
	}

	return nil
}
