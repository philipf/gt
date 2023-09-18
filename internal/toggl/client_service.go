package toggl

type ClientService interface {
	GetClients(filter string) (TogglClients, error)
}

type ClientServiceImplementation struct {
	ClientGateway ClientGateway
}

func (c *ClientServiceImplementation) GetClients(filter string) (TogglClients, error) {
	return c.ClientGateway.GetClients(filter)
}
