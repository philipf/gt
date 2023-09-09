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
	kanbanSettings := make(map[string]string)

	kanbanSettings["basePath"] = "G:/My Drive/SecondBrain"
	kanbanSettings["inboxPath"] = "_GTD/Inbox"
	kanbanSettings["boardPath"] = "_GTD/_Board.md"
	kanbanSettings["inTemplate"] = defaultInTemplate
	kanbanSettings["In"] = "In"

	viper.SetDefault("kanban", kanbanSettings)
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
