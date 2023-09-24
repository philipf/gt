//go:build wireinject
// +build wireinject

package report

import (
	"github.com/google/wire"
	"github.com/philipf/gt/internal/toggl"
	"github.com/philipf/gt/internal/toggl/http"
)

func provideTogglProjectGateway() *http.TogglProjectGateway {
	g := new(http.TogglProjectGateway)
	return g
}

func provideTogglClientGateway() *http.TogglClientGateway {
	g := new(http.TogglClientGateway)
	return g
}

func provideTogglTimeEntryGateway() *http.TogglTimeEntriesGateway {
	g := new(http.TogglTimeEntriesGateway)
	return g
}

var set = wire.NewSet(
	provideTogglProjectGateway,
	provideTogglClientGateway,
	provideTogglTimeEntryGateway,
	wire.Bind(new(toggl.ProjectGateway), new(*http.TogglProjectGateway)),
	wire.Bind(new(toggl.ClientGateway), new(*http.TogglClientGateway)),
	wire.Bind(new(toggl.TimeEntryGateway), new(*http.TogglTimeEntriesGateway)),
	toggl.NewProjectService,
	toggl.NewClientService,
	toggl.NewTimeService,
	toggl.NewReportService,
)

func initialiseReportService() toggl.ReportService {
	wire.Build(set)
	return toggl.ReportService{}
}
