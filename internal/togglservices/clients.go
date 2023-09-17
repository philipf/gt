package togglservices

import (
	"github.com/philipf/gt/internal/togglservices/gateways"
)

func GetClients(filter string) (TogglClients, error) {
	clientGateway := gateways.NewToggleClientGateway()
	return clientGateway.GetClients(filter)
}
