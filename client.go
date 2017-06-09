package gong

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/go-prompt"
	"io/ioutil"
	"os/user"
	"path/filepath"
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

	return client.Authenticate(fields), nil
}

func getUserHomeOrDefault() string {
	usr, err := user.Current()

	if err != nil {
		return "./"
	}

	return usr.HomeDir
}

func getFileLocation() string {
	dir := getUserHomeOrDefault()
	return filepath.Join(dir, ".gong.json")
}

// Load : Load the configuration from a file
func Load() (map[string]string, error) {
	fileLocation := getFileLocation()
	var c = map[string]string{}

	file, err := ioutil.ReadFile(fileLocation)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &c)

	if err != nil {
		return nil, err
	}

	return c, nil
}

// Save : saves the configuration to a file
func Save(values map[string]string) error {
	fileLocation := getFileLocation()
	loginDetails, err := json.Marshal(values)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileLocation, loginDetails, 0644)

	if err != nil {
		return err
	}

	return nil
}
