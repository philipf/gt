//go:build wireinject
// +build wireinject

package projectcmd

import (
	"github.com/google/wire"
	"github.com/philipf/gt/internal/toggl"
	"github.com/philipf/gt/internal/toggl/http"
)

func provideTogglProjectGateway() *http.TogglProjectGateway {
	g := new(http.TogglProjectGateway)
	return g
}

var set = wire.NewSet(
	provideTogglProjectGateway,
	wire.Bind(new(toggl.ProjectGateway), new(*http.TogglProjectGateway)),
	toggl.NewProjectService,
)

func initialiseProjectService() toggl.ProjectService {
	wire.Build(set)
	return toggl.ProjectService{}
}
