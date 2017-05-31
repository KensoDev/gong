package gong

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/go-ini/ini"
	"os/user"
)

type LoginDetails struct {
	Username  string
	Domain    string
	Password  string
	Directory string
}

func getUserHomeDirOrDefault() string {
	usr, err := user.Current()

	if err != nil {
		return "./"
	}

	return usr.HomeDir
}

func (loginDetails *LoginDetails) Verify() error {
	jiraClient, err := jira.NewClient(nil, loginDetails.Domain)

	if err != nil {
		return err
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(loginDetails.Username, loginDetails.Password)

	if err != nil || res == false {
		return err
	}

	return nil
}

func (loginDetails *LoginDetails) Save() error {
	cfg := ini.Empty()
	fileLocation := fmt.Sprintf("%s/.gong.ini", loginDetails.Directory)

	section := cfg.Section("")

	section.NewKey("username", loginDetails.Username)
	section.NewKey("password", loginDetails.Password)
	section.NewKey("domain", loginDetails.Domain)

	return cfg.SaveTo(fileLocation)
}

func NewLoginDetails(username string, password string, domain string) LoginDetails {
	return LoginDetails{
		Username:  username,
		Password:  password,
		Domain:    fmt.Sprintf("https://%s", domain),
		Directory: getUserHomeDirOrDefault(),
	}
}
