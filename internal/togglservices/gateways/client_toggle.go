package gateways

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"

	"github.com/philipf/gt/internal/togglservices"
	"github.com/spf13/viper"
)

type ToggleClientGateway struct {
}

func NewToggleClientGateway() ClientGateway {
	return &ToggleClientGateway{}
}

func (t *ToggleClientGateway) GetClients(filter string) (togglservices.TogglClients, error) {
	uri, err := getApiClientsList()
	if err != nil {
		return nil, err
	}

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

	var r togglservices.TogglClients
	err = json.Unmarshal(body, &r)

	if err != nil {
		return nil, err
	}

	sort.Sort(togglservices.ClientsByName(r))

	return r, nil
}
