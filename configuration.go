package gong

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

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