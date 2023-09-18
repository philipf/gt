package toggl

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

var projectTileRegEx = regexp.MustCompile(`^\[([^|]+)\|([P|S])\|(\d+)(\|(T\d{8}\.\d+))?\] (.+)$`)

type ProjectService interface {
	CreateProject(projectName string, clientID int64) error
	GetProjects(filter *GetProjectsOpts) (TogglProjects, error)

	ParseProjectTitle(project string) (ProjectTitle, error)
	ValidProjectName(name string) bool
}

type ProjectServiceImplementation struct {
	ProjectGateway ProjectGateway
}

func (p *ProjectServiceImplementation) CreateProject(projectName string, clientID int64) error {
	return p.ProjectGateway.CreateProject(projectName, clientID)

}

func (p *ProjectServiceImplementation) GetProjects(filter *GetProjectsOpts) (TogglProjects, error) {
	projects, err := p.ProjectGateway.GetProjects()
	if err != nil {
		return TogglProjects{}, err
	}

	sort.Sort(ProjectsByName(projects))
	return projects, nil
}

func (p *ProjectServiceImplementation) ParseProjectTitle(project string) (ProjectTitle, error) {
	if project == "" {
		return ProjectTitle{}, errors.New("project cannot be empty")
	}

	matches := projectTileRegEx.FindStringSubmatch(project)

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

func (p *ProjectServiceImplementation) ValidProjectName(name string) bool {
	return projectTileRegEx.MatchString(name)
}
