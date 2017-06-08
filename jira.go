package gong

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andygrunwald/go-jira"
)

type JiraClient struct {
	client *jira.Client
	config map[string]string
}

func NewJiraClient() *JiraClient {
	return &JiraClient{}
}

func (j *JiraClient) GetName() string {
	return "jira"
}

func GetIssueId(branchName string) string {
	re := regexp.MustCompile(`([A-Z]+-[\d]+)`)
	matches := re.FindAllString(branchName, -1)

	if len(matches) == 0 {
		return ""
	}

	return matches[0]
}

func (j *JiraClient) Comment(branchName, comment string) error {
	fmt.Println(branchName)
	issueId := GetIssueId(branchName)
	fmt.Println(issueId)

	jiraComment := &jira.Comment{
		Body: comment,
	}
	_, _, err := j.client.Issue.AddComment(issueId, jiraComment)

	return err
}

func (j *JiraClient) Browse(branchName string) (string, error) {
	issueId := GetIssueId(branchName)

	domain, ok := j.config["domain"]

	if !ok {
		return "", errors.New("Could not locate domain in config")
	}

	if issueId == "" {
		return "", errors.New("Could not find issue id in the branch name")
	}

	url := fmt.Sprintf("%s/browse/%s", domain, issueId)

	return url, nil
}

func (j *JiraClient) GetBranchName(issueType string, issueId string) (string, error) {
	issue, _, err := j.client.Issue.Get(issueId, nil)

	if err != nil {
		return "", err
	}

	issueTitleSlug := SlugifyTitle(issue.Fields.Summary)
	return fmt.Sprintf("%s/%s-%s", issueType, issueId, issueTitleSlug), nil
}

func indexOf(status string, data []string) int {
	for k, v := range data {
		if status == v {
			return k
		}
	}
	return -1
}

func (j *JiraClient) Start(issueType string, issueId string) (string, error) {
	allowed := []string{"Ready", "Start"}

	fmt.Println(issueId)

	transitions, response, err := j.client.Issue.GetTransitions(issueId)

	if err != nil {
		fmt.Println(err)
		fmt.Println(response.Body)
		return "", err
	}

	nextTransition := transitions[0]

	if indexOf(nextTransition.Name, allowed) > -1 {
		_, err := j.client.Issue.DoTransition(issueId, nextTransition.ID)

		if err != nil {
			return "", err
		}

		_, _ = j.Start(issueType, issueId)
	}

	branchName, err := j.GetBranchName(issueType, issueId)

	if err != nil {
		return "", err
	}

	return branchName, nil
}

func (j *JiraClient) FormatField(fieldName string, value string) string {
	if fieldName == "domain" {
		return fmt.Sprintf("https://%s", value)
	}

	return value
}

func (j *JiraClient) GetAuthFields() map[string]bool {
	return map[string]bool{
		"username":       false,
		"domain":         false,
		"password":       true,
		"project_prefix": false,
	}
}

func (j *JiraClient) Authenticate(fields map[string]string) bool {
	jiraClient, err := jira.NewClient(nil, fields["domain"])

	if err != nil {
		return false
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(fields["username"], fields["password"])

	if err != nil || res == false {
		return false
	}

	j.client = jiraClient
	j.config = fields

	return true
}
