package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/philipf/gt/internal/cache"
	"github.com/philipf/gt/internal/settings"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/viper"
)

type CreateProjectRequest struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
	IsActive  bool   `json:"active"`
	ClientID  int64  `json:"cid"`
}

type TogglProjectGateway struct {
}

func (t *TogglProjectGateway) Get(filter *toggl.GetProjectsOpts) (toggl.TogglProjects, error) {
	cache := cache.JsonFileCache[toggl.TogglProjects, toggl.GetProjectsOpts]{}

	cacheDir, err := settings.GetGtConfigPath()
	if err != nil {
		return nil, err
	}

	cacheFilePath := path.Join(cacheDir, "toggl-projects.json")
	cacheMaxAge := time.Duration(time.Minute * 5)

	if filter == nil {
		filter = &toggl.GetProjectsOpts{}
	}

	cr, err := cache.Get(*filter, cacheFilePath, cacheMaxAge)

	if err != nil {
		log.Println("Cache miss")
	} else {
		log.Println("Cache hit")
		return *cr, nil
	}

	r, err := t.getFromToggl(filter)
	if err != nil {
		return nil, err
	}

	err = cache.Save(*filter, cacheFilePath, &r, cacheMaxAge)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (t *TogglProjectGateway) getFromToggl(filter *toggl.GetProjectsOpts) (toggl.TogglProjects, error) {
	uri, err := getCreateProjectUri()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if filter == nil {
		q.Add("active", "true")
	} else {
		if !filter.IncludeArchived {
			q.Add("active", "true")
		}

		if filter.Name != "" {
			q.Add("name", filter.Name)
		}

		// if len(filter.ClientIDs) > 0 {
		// 	for _, clientID := range filter.ClientIDs {
		// 		q.Add("client_ids", fmt.Sprintf("%d", clientID))
		// 	}
		// }

		if filter.ClientID != 0 {
			q.Add("client_ids", fmt.Sprintf("%d", filter.ClientID))
		}
	}

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

func (t *TogglProjectGateway) Create(projectName string, clientID int64) error {
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
