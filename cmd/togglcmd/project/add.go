package project

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Long: `Adds a project using a specific naming convention and client association
Example:
gt toggl project add --name "Project Name" --clientID 12345

To obtain the clientID, run the following command:
gt toggl client list
`,

	Run: func(cmd *cobra.Command, args []string) {

		interactive, _ := cmd.Flags().GetBool("interactive")

		if interactive {
			doInteractive(cmd)
		} else {
			addProject(cmd)
		}

	},
}

func doInteractive(cmd *cobra.Command) {
	fmt.Println("Interactive mode")

	// Step 1: Get project type
	var projectType string
	projectTypeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Task Type").
				Options(
					huh.NewOption("Service Desk", "S"),
					huh.NewOption("Project", "P"),
				).
				Value(&projectType),
		),
	)
	
	err := projectTypeForm.Run()
	cobra.CheckErr(err)
	
	if projectType != "S" && projectType != "P" {
		fmt.Println("Invalid project type")
		return
	}
	
	// Get client filter from flags
	clientFilter, _ := cmd.Flags().GetString("clientName")
	clientFilter += "%|" + projectType
	
	fmt.Printf("Searching for clients with filter: %s\n", clientFilter)
	
	// Step 2: Get clients
	cf := toggl.GetClientOpts{
		Name: clientFilter,
	}
	
	clients, err := clientService.Get(&cf)
	cobra.CheckErr(err)
	
	if len(clients) == 0 {
		fmt.Println("No clients found for: " + clientFilter)
		return
	}
	
	// Step 3: Select client
	var selectedClientIndex int
	var clientOptions []huh.Option[int]
	
	for i, client := range clients {
		clientOptions = append(clientOptions, huh.NewOption(client.Name, i))
	}
	
	clientForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select Client").
				Options(clientOptions...).
				Value(&selectedClientIndex).
				Validate(func(i int) error {
					if i < 0 || i >= len(clients) {
						return fmt.Errorf("Invalid client selection")
					}
					return nil
				}),
		),
	)
	
	err = clientForm.Run()
	cobra.CheckErr(err)
	
	clientID := clients[selectedClientIndex].ID
	clientName := clients[selectedClientIndex].Name
	
	// Step 4: Get project details
	var taskIDStr string
	var ticketNo string
	var projectName string
	
	// Create project details form
	detailsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Auto Task ID").
				Validate(func(s string) error {
					// Simple validation to ensure it's a number
					if s == "" {
						return fmt.Errorf("Auto Task ID is required")
					}
					_, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return fmt.Errorf("Auto Task ID must be a number")
					}
					return nil
				}).
				Value(&taskIDStr),
		),
	)
	
	// For Service Desk type, we need a ticket number field
	if projectType == "S" {
		detailsForm = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Auto Task ID").
					Validate(func(s string) error {
						// Simple validation to ensure it's a number
						if s == "" {
							return fmt.Errorf("Auto Task ID is required")
						}
						_, err := strconv.ParseInt(s, 10, 64)
						if err != nil {
							return fmt.Errorf("Auto Task ID must be a number")
						}
						return nil
					}).
					Value(&taskIDStr),
					
				huh.NewInput().
					Title("Auto Task Ticket Number").
					Value(&ticketNo),
			),
		)
	}
	
	err = detailsForm.Run()
	cobra.CheckErr(err)
	
	// Step 5: Get project name
	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("Project name is required")
					}
					return nil
				}).
				Value(&projectName),
		),
	)
	
	err = nameForm.Run()
	cobra.CheckErr(err)
	
	// Parse task ID
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	cobra.CheckErr(err)
	
	// Format project name based on project type
	var finalProjectName string
	if projectType == "S" {
		finalProjectName = fmt.Sprintf("[%s|%d|%s] %s", clientName, taskID, ticketNo, projectName)
	} else {
		finalProjectName = fmt.Sprintf("[%s|%d] %s", clientName, taskID, projectName)
	}
	
	// Step 6: Confirm project creation
	var confirmCreate bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Project Preview").
				Description(fmt.Sprintf("Project will be created as: %s", finalProjectName)),
				
			huh.NewConfirm().
				Title("Create this project?").
				Value(&confirmCreate),
		),
	)
	
	err = confirmForm.Run()
	cobra.CheckErr(err)
	
	if !confirmCreate {
		fmt.Println("Project creation cancelled")
		return
	}
	
	// Create the project
	err = projectService.Create(finalProjectName, clientID)
	cobra.CheckErr(err)
	fmt.Println("Project created successfully")
}

func addProject(cmd *cobra.Command) {
	name, _ := cmd.Flags().GetString("name")
	clientID, _ := cmd.Flags().GetInt64("clientID")

	if clientID <= 0 {
		fmt.Println("clientID is required")
		return
	}

	if strings.TrimSpace(name) == "" {
		fmt.Println("name is required")
		return
	}

	err := projectService.Create(name, clientID)
	cobra.CheckErr(err)
}

func init() {
	projectCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "", "Name of the project")
	addCmd.Flags().Int64P("clientID", "c", 0, "ClientID that the project belongs to")
	addCmd.Flags().StringP("clientName", "f", "", "Client filter by name")
	addCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
}
