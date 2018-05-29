package gong

import (
	"fmt"
	"github.com/segmentio/go-prompt"
)

// Client : Public interface for the generic client
type Client interface {
	GetAuthFields() map[string]bool
	GetName() string
	FormatField(fieldName string, value string) string
	Authenticate(fields map[string]string) bool
	Start(issueType string, issueID string) (branchName string, err error)
	Browse(branchName string) (string, error)
	Comment(branchName, comment string) error
	PrepareCommitMessage(branchName, commitMessage string) string
	Create() (string, error)
}

func Create(client Client) (string, error) {
	return client.Create()
}

// PrepareCommitMessage : Prepares the commit message and returns a new commit message
func PrepareCommitMessage(client Client, branchName, commitMessage string) string {
	return client.PrepareCommitMessage(branchName, commitMessage)
}
	
// Comment : Comment on an issue
func Comment(client Client, branchName, comment string) error {
	return client.Comment(branchName, comment)
}

// Browse : Open a browser on the issue related to the branch
func Browse(client Client, branchName string) (string, error) {
	return client.Browse(branchName)
}

// Start : Start working on an issue
func Start(client Client, issueType, issueID string) (string, error) {
	return client.Start(issueType, issueID)
}

// NewClient : Return a new client that matches the name passed in
func NewClient(clientName string) (Client, error) {
	if clientName == "jira" {
		return NewJiraClient(), nil
	}

	if clientName == "pivotal" {
		return NewPivotalClient(), nil
	}

	return nil, fmt.Errorf("Could not find client: %v", clientName)
}
	
// NewAuthenticatedClient : Return a new client authenticated
func NewAuthenticatedClient() (Client, error) {
	fields, err := Load()

	if err != nil {
		return nil, err
	}

	client, err := NewClient(fields["client"])

	if err != nil {
		return nil, err
	}

	authenticated := client.Authenticate(fields)

	if authenticated {
		return client, nil
	}

	return nil, fmt.Errorf("Could not load authenticated client")
}

// Login : Logs the user into the specified client
func Login(client Client) (bool, error) {
	clientName := client.GetName()

	fields := map[string]string{
		"client": clientName,
	}

	for k, v := range client.GetAuthFields() {
		message := fmt.Sprintf("Please enter your %v %v", clientName, k)
		promptValue := ""
		if v {
			promptValue = prompt.PasswordMasked(message)
		} else {
			promptValue = prompt.String(message)
		}

		fields[k] = client.FormatField(k, promptValue)
	}

	err := Save(fields)

	if err != nil {
		return false, err
	}

	fields, err = Load()

	if err != nil {
		return false, err
	}

	authenticated := client.Authenticate(fields)

	if authenticated {
		return true, nil
	}

	return false, fmt.Errorf("Cloud not login")
}
