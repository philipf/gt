package project

import (
	"fmt"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

var filter toggl.GetProjectsOpts

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive projects",
	Long: `Archives projects using the Toggl API

Example:
gt toggl project archive
`,

	Run: func(cmd *cobra.Command, _ []string) {
		clientId, err := cmd.Flags().GetInt64("clientId")
		cobra.CheckErr(err)

		name, err := cmd.Flags().GetString("name")
		cobra.CheckErr(err)

		filter = toggl.GetProjectsOpts{
			ClientID: clientId,
			Name:     name,
		}

		items, err := projectService.Get(&filter)
		cobra.CheckErr(err)

		m := model{
			projects: items,
			selected: make(map[int]struct{}),
		}

		p := tea.NewProgram(m)
		_, err = p.Run()
		cobra.CheckErr(err)

		// print all the selected projects
		for i := range m.selected {
			err := projectService.Archive(m.projects[i].ID)
			cobra.CheckErr(err)
			fmt.Println("Archived project", m.projects[i].Name)
		}
	},
}

func fetchProjects() tea.Msg {
	return func() tea.Msg {
		projects, err := projectService.Get(&filter)
		if err != nil {
			return errMsg{err}
		}
		return projects
	}
}

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func init() {
	projectCmd.AddCommand(archiveCmd)

	// filter
	archiveCmd.Flags().Int64P("clientId", "c", 0, "Filter by client ID")
	archiveCmd.Flags().StringP("name", "n", "", "Filter by project name")
}

type model struct {
	projects toggl.TogglProjects
	cursor   int
	selected map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return fetchProjects
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "w":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.projects) - 1
			}
		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Select projects to archive [Q]uit/[W]rite:\n\n"

	for i, p := range m.projects {
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, p.Name)
	}

	return s
}
