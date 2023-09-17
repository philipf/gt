package togglservices

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"

	"github.com/spf13/viper"
)

func GetClients(filter string) (TogglClients, error) {
	uri := fmt.Sprintf("%s/api/v9/workspaces/%s/clients", BASE_URI, getWorkspaceID())

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if filter != "" {
		q := u.Query()
		q.Add("name", filter)
		u.RawQuery = q.Encode()
	}

	log.Println("URI:", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clients, err := unmarshalTogglClient(body)

	if err != nil {
		return nil, err
	}

	sort.Sort(ClientsByName(clients))

	return clients, nil
}

func unmarshalTogglClient(data []byte) (TogglClients, error) {
	var r TogglClients
	err := json.Unmarshal(data, &r)
	return r, err
}
