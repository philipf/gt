package togglservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"

	"github.com/spf13/viper"
)

type CreateProjectRequest struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
	IsActive  bool   `json:"active"`
	ClientId  int64  `json:"cid"`
}

func CreateProject(projectName string, clientId int64) error {

	project := CreateProjectRequest{
		Name:      projectName,
		IsPrivate: true,
		IsActive:  true,
		ClientId:  clientId,
	}

	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/workspaces/%s/projects", BASE_URI, getWorkspaceId())
	log.Println("URI:", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d failed to create project: %s", resp.StatusCode, body)
	}

	fmt.Println("Project created successfully:", string(body))
	return nil
}

type GetProjectsOpts struct {
	Name string
	//Active bool
}

func GetProjects(filter *GetProjectsOpts) (ToggleProjects, error) {
	uri := fmt.Sprintf("%s/workspaces/%s/projects", BASE_URI, getWorkspaceId())

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if filter != nil {
		if filter.Name != "" {
			q.Add("name", filter.Name)
		}

		// if !filter.Active {
		// 	q.Add("active", "false")
		// }
	}

	q.Add("active", "true")
	u.RawQuery = q.Encode()

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

	projects, err := unmarshalToggleProject(body)

	if err != nil {
		return nil, err
	}

	sort.Sort(ProjectsByName(projects))

	return projects, nil
}

type ToggleProjects []ToggleProjectElement

func unmarshalToggleProject(data []byte) (ToggleProjects, error) {
	var r ToggleProjects
	err := json.Unmarshal(data, &r)
	return r, err
}

type ToggleProjectElement struct {
	ID                  int64       `json:"id"`
	WorkspaceID         int64       `json:"workspace_id"`
	ClientID            int64       `json:"client_id"`
	Name                string      `json:"name"`
	IsPrivate           bool        `json:"is_private"`
	Active              bool        `json:"active"`
	At                  string      `json:"at"`
	CreatedAt           string      `json:"created_at"`
	ServerDeletedAt     interface{} `json:"server_deleted_at"`
	Color               string      `json:"color"`
	Billable            interface{} `json:"billable"`
	Template            interface{} `json:"template"`
	AutoEstimates       interface{} `json:"auto_estimates"`
	EstimatedHours      interface{} `json:"estimated_hours"`
	Rate                interface{} `json:"rate"`
	RateLastUpdated     interface{} `json:"rate_last_updated"`
	Currency            interface{} `json:"currency"`
	Recurring           bool        `json:"recurring"`
	RecurringParameters interface{} `json:"recurring_parameters"`
	CurrentPeriod       interface{} `json:"current_period"`
	FixedFee            interface{} `json:"fixed_fee"`
	ActualHours         int64       `json:"actual_hours"`
}

// ProjectsByName implements sort.Interface based on the Name field of ToggleClientElement.
type ProjectsByName ToggleProjects

func (a ProjectsByName) Len() int           { return len(a) }
func (a ProjectsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProjectsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
