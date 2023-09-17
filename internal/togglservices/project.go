package togglservices

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/spf13/viper"
)

func CreateProject(projectName string, clientID int64) error {

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

	uri := fmt.Sprintf("%s/api/v9/workspaces/%s/projects", BASE_URI, getWorkspaceID())
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

func GetProjects(filter *GetProjectsOpts) (TogglProjects, error) {
	uri := fmt.Sprintf("%s/api/v9/workspaces/%s/projects", BASE_URI, getWorkspaceID())

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

	projects, err := unmarshalTogglProject(body)

	if err != nil {
		return nil, err
	}

	sort.Sort(ProjectsByName(projects))

	return projects, nil
}

func unmarshalTogglProject(data []byte) (TogglProjects, error) {
	var r TogglProjects
	err := json.Unmarshal(data, &r)
	return r, err
}

func ParseProjectTitle(project string) (ProjectTitle, error) {
	if project == "" {
		return ProjectTitle{}, errors.New("project cannot be empty")
	}

	matches := ProjectTileRegEx.FindStringSubmatch(project)

	if matches == nil {
		return ProjectTitle{}, errors.New("project does not match the naming convention")
	}

	taskID, err := strconv.Atoi(matches[3])
	if err != nil {
		return ProjectTitle{}, fmt.Errorf("failed to convert TaskID to integer: %v", err)
	}

	return ProjectTitle{
		Client:   matches[1],
		IsTask:   (matches[2] == "S"),
		TaskID:   taskID,
		TicketID: matches[5], // Note: [4] would be the entire optional group including the pipe
		Project:  matches[6],
	}, nil
}

func ValidProjectName(name string) bool {
	return ProjectTileRegEx.MatchString(name)
}
