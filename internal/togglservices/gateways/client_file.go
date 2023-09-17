// Implements the ClientGateway by reading from a file
package gateways

import (
	"encoding/json"
	"os"
	"path"

	"github.com/philipf/gt/internal/settings"
	"github.com/philipf/gt/internal/togglservices"
)

type FileClientGateway struct {
}

func NewFileClientGateway() ClientGateway {
	return &FileClientGateway{}
}

const filename = "toggle-clients.config"

func (f *FileClientGateway) GetClients(filter string) (togglservices.TogglClients, error) {
	d, err := settings.GetGtConfigPath()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(d, filename)

	// Check if the file exists
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if os.IsNotExist(err) {
		// If the file does not exist, return an empty list
		return togglservices.TogglClients{}, nil
	}

	// Read the contents of the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a TogglClients slice
	var clients togglservices.TogglClients
	err = json.Unmarshal(fileContent, &clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}
