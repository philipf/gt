package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/viper"
)

type ToggleProjectGateway struct {
}

type CreateProjectRequest struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
	IsActive  bool   `json:"active"`
	ClientID  int64  `json:"cid"`
}

func (t *ToggleProjectGateway) GetProjects() (toggl.TogglProjects, error) {
	return toggl.TogglProjects{}, errors.New("not implemented")
}

func (t *ToggleProjectGateway) CreateProject(projectName string, clientID int64) error {
	uri, err := getCreateProjectUri()
	if err != nil {
		return err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	project := CreateProjectRequest{
		Name:      projectName,
		IsPrivate: true,
		IsActive:  true,
		ClientID:  clientID,
	}

	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(data))
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

func GetProjects(filter *toggl.GetProjectsOpts) (toggl.TogglProjects, error) {
	uri, err := getCreateProjectUri()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if filter != nil {
		if filter.Name != "" {
			q.Add("name", filter.Name)
		}
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

	var projects toggl.TogglProjects
	err = json.Unmarshal(body, &projects)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
