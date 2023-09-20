package toggl

type ClientService struct {
	ClientGateway ClientGateway
}

func (c *ClientService) GetClients(filter string) (TogglClients, error) {
	gw := c.ClientGateway
	return gw.GetClients(filter)
}

func NewClientService(clientGateway ClientGateway) ClientService {
	s := ClientService{}
	s.ClientGateway = clientGateway
	return s
}
