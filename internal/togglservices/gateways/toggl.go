package gateways

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	UriBase       = "https://api.track.toggl.com/api/v9"
	UriClientsGet = "%s/api/v9/workspaces/%s/clients"
)

func getAPIToken() (string, error) {
	r := viper.GetString("toggl.apiKey")

	if r == "" {
		return "", errors.New("No API token found")
	}

	return r, nil
}

func getWorkspaceID() (string, error) {
	r := viper.GetString("toggl.workspace")

	if r == "" {
		return "", errors.New("No workspace ID found")
	}

	return r, nil
}

//uri := fmt.Sprintf("%s/api/v9/workspaces/%s/clients", BASE_URI, getWorkspaceID())

func getApiClientsList() (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriClientsGet, UriBase, workspaceID)
	return uri, nil
}
