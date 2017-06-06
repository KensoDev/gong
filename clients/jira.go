package jiraapi

import (
	"fmt"
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
