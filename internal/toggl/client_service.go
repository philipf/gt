package toggl

type ClientService struct {
	ClientGateway ClientGateway
}

func (c *ClientService) Get(filter string) (TogglClients, error) {
	gw := c.ClientGateway
	return gw.Get(filter)
}

func NewClientService(clientGateway ClientGateway) ClientService {
	s := ClientService{}
	s.ClientGateway = clientGateway
	return s
}
