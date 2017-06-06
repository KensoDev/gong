package jiraapi

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

type JiraClient struct{}

func NewJiraClient() JiraClient {
	return JiraClient{}
}

func (j JiraClient) GetName() string {
	return "jira"
}

func (j JiraClient) FormatField(fieldName string, value string) string {
	if fieldName == "domain" {
		return fmt.Sprintf("https://%s", value)
	}

	return value
}

func (j JiraClient) GetAuthFields() map[string]bool {
	return map[string]bool{
		"username":       false,
		"domain":         false,
		"password":       true,
		"project_prefix": false,
	}
}

func (j JiraClient) Authenticate(fields map[string]string) bool {
	jiraClient, err := jira.NewClient(nil, fields["domain"])

	if err != nil {
		return false
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(fields["username"], fields["password"])

	if err != nil || res == false {
		return false
	}

	return true
}
