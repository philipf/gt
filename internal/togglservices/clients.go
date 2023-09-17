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

func GetClients(filter string) (ToggleClients, error) {
	uri := fmt.Sprintf("%s/workspaces/%s/clients", BASE_URI, getWorkspaceId())

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

	clients, err := UnmarshalToggleClient(body)

	if err != nil {
		return nil, err
	}

	sort.Sort(ByName(clients))

	return clients, nil
}

type ToggleClients []ToggleClientElement

func UnmarshalToggleClient(data []byte) (ToggleClients, error) {
	var r ToggleClients
	err := json.Unmarshal(data, &r)
	return r, err
}

type ToggleClientElement struct {
	ID       int64  `json:"id"`
	Wid      int64  `json:"wid"`
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
	At       string `json:"at"`
}

// ByName implements sort.Interface based on the Name field of ToggleClientElement.
type ByName ToggleClients

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
