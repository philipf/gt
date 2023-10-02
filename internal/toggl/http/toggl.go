package http

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	UriBase            = "https://api.track.toggl.com/api/v9"
	UriClients         = "%s/workspaces/%s/clients"
	UriTimeEntries     = "%s/me/time_entries"
	UriTimeEntriesStop = "%s/workspaces/%s/time_entries/%d/stop"
	UriProject         = "%s/workspaces/%s/projects"
)

func getAPIToken() string {
	r := viper.GetString("toggl.apiKey")

	return r
}

func getWorkspaceID() (string, error) {
	r := viper.GetString("toggl.workspace")

	if r == "" {
		return "", errors.New("no workspace ID found")
	}

	return r, nil
}

func getApiClientsListUri() (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriClients, UriBase, workspaceID)
	return uri, nil
}

func getTimeEntriesUri() (string, error) {
	uri := fmt.Sprintf(UriTimeEntries, UriBase)
	return uri, nil
}

func getTimeEntriesStopUri(entryID int64) (string, error) {
	workspaceID, err := getWorkspaceID()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(UriTimeEntriesStop, UriBase, workspaceID, entryID)
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
