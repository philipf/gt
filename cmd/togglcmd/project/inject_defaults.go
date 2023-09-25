//go:build wireinject
// +build wireinject

package project

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

func initialiseClientService() toggl.ClientService {
	wire.Build(clientSet)
	return toggl.ClientService{}
}

func provideTogglClientGateway() *http.TogglClientGateway {
	g := new(http.TogglClientGateway)
	return g
}

var clientSet = wire.NewSet(
	provideTogglClientGateway,
	wire.Bind(new(toggl.ClientGateway), new(*http.TogglClientGateway)),
	toggl.NewClientService,
)
