package settings

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultInTemplate = "in.md"
	gtHomePath        = ".gt"

	CONFIG_FILE_NAME = "config"
	CONFIG_FILE_TYPE = "yaml"
)

func Init() error {
	// Check if there is a config file in the home directory for gt, if not, create one otherwise use it

	// If there is no config file in the home directory, create one
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("cannot find user's home directory [%v]", err)
	}

	gtHomePath := filepath.Join(homeDir, gtHomePath)

	// Check if the directory exists, if not, create it
	_, err = os.Stat(gtHomePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = initialiseConfigFile(gtHomePath)
			if err != nil {
				return fmt.Errorf("cannot initialise config file in user's home directory: [%v]", err)
			}
		} else {
			return fmt.Errorf("cannot stat .gt directory in user's home directory:  [%v]", err)
		}
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath(gtHomePath)

	viper.ReadInConfig()

	return nil
}

func initialiseConfigFile(gtHomePath string) error {
	err := os.Mkdir(gtHomePath, 0755)
	if err != nil {
		return fmt.Errorf("cannot .gt directory in user's home directory:  [%v]", err)
	}

	viper.AddConfigPath(gtHomePath)
	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType(CONFIG_FILE_TYPE)

	setViperDefaults()

	fmt.Println("Creating a default config file in: ", gtHomePath)

	err = viper.SafeWriteConfig()
	if err != nil {
		return fmt.Errorf("cannot write config file in user's home directory:  [%v]", err)
	}

	// write default in.md file
	inTemplate := getDefaultInTemplate()

	os.WriteFile(filepath.Join(gtHomePath, defaultInTemplate), []byte(inTemplate), 0644)

	return nil
}

func getDefaultInTemplate() string {
	inTemplate := `---
type: kanban
created: {{ .CreatedAt.Format "2006-01-02 15:04" }}
updated: {{ .UpdatedAt.Format "2006-01-02 15:04" }}
status: {{.Status}}
channel: {{.Channel}}
externalId: {{.ExternalID}}
---

# {{.Title}}
{{.Description}}

## Notes
- [ ] Step 1
`
	return inTemplate
}

func setViperDefaults() {
	// Personal settings
	personalSettings := make(map[string]interface{})
	personalSettings["name"] = "<name>"
	personalSettings["surname"] = "<surname>"
	personalSettings["company"] = "<company>"
	personalSettings["email"] = "<email optional>"

	viper.SetDefault("personal", personalSettings)

	// Kanban settings
	kanbanSettings := make(map[string]interface{})

	kanbanSettings["basePath"] = "G:/My Drive/SecondBrain"
	kanbanSettings["inboxPath"] = "_GTD"
	kanbanSettings["boardPath"] = "_GTD/_Board.md"
	kanbanSettings["inTemplate"] = defaultInTemplate
	kanbanSettings["In"] = "In"

	viper.SetDefault("kanban", kanbanSettings)

	// AI settings
	openAiSettings := make(map[string]interface{})
	openAiSettings["openAiKey"] = "<your key here from https://platform.openai.com/account/api-keys>"
	//openAiSettings["openAiModel"] = "gpt-3.5-turbo"
	openAiSettings["openAiModel"] = "gpt-4"

	viper.SetDefault("ai", openAiSettings)

	// GTD settings
	gtdSettings := make(map[string]interface{})
	gtdSettings["actionPrompt.user"] = `I am {{.name}} {{.surname}} and work for {{.company}}, you are an efficient task master using the Getting Things Done method, your job is to extract the next action for me and to provide a summary from supplied emails or online conversations. 
It is now {{.currentDate}}, use: 
{{.input}}`
	gtdSettings["agctionPrompt.temperature"] = 0.7

	viper.SetDefault("gtd", gtdSettings)

}

func GetInTemplatePath() string {
	inTemplate := viper.GetString("kanban.inTemplate")
	homeDir, _ := os.UserHomeDir()

	return filepath.Join(homeDir, gtHomePath, inTemplate)
}

func GetKanbanBasePath() string {
	return viper.GetString("kanban.basePath")
}

func GetKanbanGtdPath() string {
	basePath := GetKanbanBasePath()
	return filepath.Join(basePath, viper.GetString("kanban.gtdPath"))
}

func GetKanbanInboxPath() string {
	gtdPath := GetKanbanGtdPath()
	return filepath.Join(gtdPath, viper.GetString("kanban.inboxPath"))
}

func GetKanbanBoardPath() string {
	gtdPath := GetKanbanGtdPath()
	return filepath.Join(gtdPath, viper.GetString("kanban.boardPath"))
}

func GetKanbanInColumn() string {
	return viper.GetString("kanban.In")
}
