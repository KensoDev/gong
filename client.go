package gong

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kensodev/gong/clients"
	"github.com/segmentio/go-prompt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type Client interface {
	GetAuthFields() map[string]bool
	GetName() string
	FormatField(fieldName string, value string) string
}

func NewClient(clientName string) (Client, error) {
	if clientName == "jira" {
		return jiraapi.NewJiraClient(), nil
	}

	return nil, errors.New(fmt.Sprintf("Could not find client: %v", clientName))
}

func Login(client Client) error {
	fields := map[string]string{
		"client": client.GetName(),
	}

	for k, v := range client.GetAuthFields() {
		message := fmt.Sprintf("Please enter your jira %v", k)
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
		return err
	}

	return nil
}

func GetLoginDetailsFileLocation() string {
	usr, err := user.Current()

	if err != nil {
		return "./"
	}

	return usr.HomeDir
}

func Save(values map[string]string) error {
	dir := GetLoginDetailsFileLocation()
	fileLocation := filepath.Join(dir, ".gong.json")

	fmt.Println(fileLocation)
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
