package togglservices

import (
	"log"

	"github.com/spf13/viper"
)

const (
	BASE_URI = "https://api.track.toggl.com/api/v9"
)

func getWorkspaceID() string {
	r := viper.GetString("toggl.workspace")

	if r == "" {
		log.Fatal("No workspace id found in config file")
	}

	return r
}
