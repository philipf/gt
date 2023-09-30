// Implements the ClientGateway by reading from a file
package file

import (
	"encoding/json"
	"os"
	"path"

	"github.com/philipf/gt/internal/settings"
	"github.com/philipf/gt/internal/toggl"
)

type FileClientGateway struct {
}

const filename = "toggl-clients.json"

func (f *FileClientGateway) Get(filter string) (toggl.TogglClients, error) {
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
		return toggl.TogglClients{}, nil
	}

	// Read the contents of the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a TogglClients slice
	var clients toggl.TogglClients
	err = json.Unmarshal(fileContent, &clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}
