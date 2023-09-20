//go:build wireinject
// +build wireinject

package client

import (
	"github.com/google/wire"
	"github.com/philipf/gt/internal/toggl"
	"github.com/philipf/gt/internal/toggl/http"
)

func provideToggleClientGateway() *http.TogglClientGateway {
	g := new(http.TogglClientGateway)
	return g
}

var set = wire.NewSet(
	provideToggleClientGateway,
	wire.Bind(new(toggl.ClientGateway), new(*http.TogglClientGateway)),
	toggl.NewClientService,
)

func initializeClientService() toggl.ClientService {
	wire.Build(set)
	return toggl.ClientService{}
}
