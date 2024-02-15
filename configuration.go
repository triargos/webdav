package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	AdminUserName       string `json:"adminUserName"`
	CredentialsFilePath string `json:"credentialsFilePath"`
	PermissionsFilePath string `json:"permissionsFilePath"`
	Realm               string `json:"realm"`
	UsersRoot           string `json:"usersRoot"`
	DataPath            string `json:"dataPath"`
	CertPath            string `json:"certPath"`
	KeyPath             string `json:"keyPath"`
}

func readConfiguration() (*Configuration, error) {
	file, err := os.Open("conf.json")
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)
	decoder := json.NewDecoder(file)
	configuration := &Configuration{}
	err = decoder.Decode(configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
