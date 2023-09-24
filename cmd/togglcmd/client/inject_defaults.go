//go:build wireinject
// +build wireinject

package clientcmd

import (
	"github.com/google/wire"
	"github.com/philipf/gt/internal/toggl"
	"github.com/philipf/gt/internal/toggl/http"
)

func provideTogglClientGateway() *http.TogglClientGateway {
	g := new(http.TogglClientGateway)
	return g
}

var set = wire.NewSet(
	provideTogglClientGateway,
	wire.Bind(new(toggl.ClientGateway), new(*http.TogglClientGateway)),
	toggl.NewClientService,
)

func initialiseClientService() toggl.ClientService {
	wire.Build(set)
	return toggl.ClientService{}
}
