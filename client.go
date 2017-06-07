package gong

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/go-prompt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type Client interface {
	GetAuthFields() map[string]bool
	GetName() string
	FormatField(fieldName string, value string) string
	Authenticate(fields map[string]string) bool
	Start(issueType string, issueId string) (branchName string, err error)
}

func Start(client Client, issueType, issueId string) (string, error) {
	return client.Start(issueType, issueId)
}

func NewClient(clientName string) (Client, error) {
	if clientName == "jira" {
		return NewJiraClient(), nil
	}

	return nil, errors.New(fmt.Sprintf("Could not find client: %v", clientName))
}

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

	return nil, errors.New("Could not load authenticated client")
}

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

func GetUserHomeOrDefault() string {
	usr, err := user.Current()

	if err != nil {
		return "./"
	}

	return usr.HomeDir
}

func GetFileLocation() string {
	dir := GetUserHomeOrDefault()
	return filepath.Join(dir, ".gong.json")
}

func Load() (map[string]string, error) {
	fileLocation := GetFileLocation()
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

func Save(values map[string]string) error {
	fileLocation := GetFileLocation()
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
