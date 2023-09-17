package project

import (
	"fmt"
	"strings"

	"github.com/philipf/gt/internal/console"
	"github.com/philipf/gt/internal/togglservices"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Long: `Adds a project using a specific naming convention and client association
Example:
gt toggle project add --name "Project Name" --clientId 12345

To obtain the clientId, run the following command:
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

	projectType, err := console.Prompt("Task type [S]ervice Desk or [P]roject: ")
	cobra.CheckErr(err)

	projectType = strings.ToUpper(projectType)
	if projectType != "S" && projectType != "P" {
		fmt.Println("Invalid project type")
		return
	}

	// get client filter from flags
	clientFilter, _ := cmd.Flags().GetString("clientName")
	clientFilter += "%|" + projectType

	fmt.Printf("Searching for clients with filter: %s\n", clientFilter)
	clients, err := togglservices.GetClients(clientFilter)
	cobra.CheckErr(err)

	var clientId int64 = 0

	if len(clients) == 0 {
		fmt.Println("No clients found for: " + clientFilter)
		return
	} else if len(clients) == 1 {
		clientId = clients[0].ID
	} else {
		// More than one found, ask the user to select one
		fmt.Println("More than one client found, please select one:")
		for _, i := range clients {
			fmt.Printf("\t%d - %s\n", i.ID, i.Name)
		}

		clientId, err = console.PromptInt64("Client Id: ")
		cobra.CheckErr(err)
	}

	// get client name from clients list by id
	clientName := ""
	for _, i := range clients {
		if i.ID == clientId {
			clientName = i.Name
			break
		}
	}

	if clientName == "" {
		fmt.Println("Invalid client id")
		return
	}

	// if the type is [S]ervice Desk then get the Auto Task ticket number

	var ticketNo string = ""
	if projectType == "S" {
		// obtain the auto task id from the user
		ticketNo, err = console.Prompt("What is the Auto Task Ticket #: ")
		cobra.CheckErr(err)
	}

	// obtain the auto task id from the user
	taskId, err := console.PromptInt64("Auto Task Id: ")
	cobra.CheckErr(err)

	// obtain the project name from the user
	projectName, err := console.Prompt("Project Name: ")
	cobra.CheckErr(err)

	if projectType == "S" {
		projectName = fmt.Sprintf("[%s|%d|%s] %s", clientName, taskId, ticketNo, projectName)
	} else {
		projectName = fmt.Sprintf("[%s|%d] %s", clientName, taskId, projectName)
	}

	// get confirmation from the user
	confirm, err := console.Prompt("Confirm project name: " + projectName + " [Y/N]: ")
	cobra.CheckErr(err)

	confirm = strings.ToUpper(confirm)
	if confirm == "Y" {
		err = togglservices.CreateProject(projectName, clientId)
		cobra.CheckErr(err)
		fmt.Println("Project created successfully")
	} else {
		fmt.Println("Project not created, user cancelled")
	}
}

func addProject(cmd *cobra.Command) {
	name, _ := cmd.Flags().GetString("name")
	clientId, _ := cmd.Flags().GetInt64("clientId")

	if clientId <= 0 {
		fmt.Println("clientId is required")
		return
	}

	if strings.TrimSpace(name) == "" {
		fmt.Println("name is required")
		return
	}

	err := togglservices.CreateProject(name, clientId)
	cobra.CheckErr(err)
}

func init() {
	projectsCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "", "Name of the project")
	addCmd.Flags().Int64P("clientId", "c", 0, "ClientId that the project belongs to")
	addCmd.Flags().StringP("clientName", "f", "", "Client filter by name")
	addCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")

}
