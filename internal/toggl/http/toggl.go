package http

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	UriBase        = "https://api.track.toggl.com/api/v9"
	UriClients     = "%s/workspaces/%s/clients"
	UriTimeEntries = "%s/me/time_entries"
	UriProject     = "%s/workspaces/%s/projects"
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

func getApiClientsListUri() (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriClients, UriBase, workspaceID)
	return uri, nil
}

func getTimeEntriesUri() (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriClients, UriBase, workspaceID)
	return uri, nil
}

func getCreateProjectUri() (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriProject, UriBase, workspaceID)
	return uri, nil
}
