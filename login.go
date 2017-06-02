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

func GetDefaultLocationLoginDetails() LoginDetails {
	return LoginDetails{
		Directory: getUserHomeDirOrDefault(),
	}
}

func GetAuthenticatedClient() (*jira.Client, error) {
	loginDetails := GetDefaultLocationLoginDetails()
	loginDetails, err := Load(loginDetails)

	if err != nil {
		return nil, err
	}

	return loginDetails.GetClient()
}

func Load(loginDetails LoginDetails) (LoginDetails, error) {
	fileLocation := loginDetails.GetLoginDetailsFileLocation()
	cfg, err := ini.InsensitiveLoad(fileLocation)

	if err != nil {
		return LoginDetails{}, err
	}

	section, err := cfg.GetSection("")
	if err != nil {
		return LoginDetails{}, err
	}

	username, err := section.GetKey("username")
	if err != nil {
		return LoginDetails{}, err
	}

	password, err := section.GetKey("password")
	if err != nil {
		return LoginDetails{}, err
	}

	domain, err := section.GetKey("domain")
	if err != nil {
		return LoginDetails{}, err
	}

	loginDetails = LoginDetails{
		Username: username.String(),
		Password: password.String(),
		Domain:   domain.String(),
	}

	return loginDetails, nil
}

func (loginDetails *LoginDetails) GetLoginDetailsFileLocation() string {
	return fmt.Sprintf("%s/.gong.ini", loginDetails.Directory)
}

func getUserHomeDirOrDefault() string {
	usr, err := user.Current()

	if err != nil {
		return "./"
	}

	return usr.HomeDir
}

func (loginDetails *LoginDetails) GetClient() (*jira.Client, error) {
	jiraClient, err := jira.NewClient(nil, loginDetails.Domain)

	if err != nil {
		return nil, err
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(loginDetails.Username, loginDetails.Password)

	if err != nil || res == false {
		return nil, err
	}

	return jiraClient, nil
}

func (loginDetails *LoginDetails) Save() error {
	cfg := ini.Empty()
	fileLocation := loginDetails.GetLoginDetailsFileLocation()

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
