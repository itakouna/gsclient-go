package gsclient

import (
	"net/http"
	"path"
)

//IPList is JSON struct of a list of IPs
type IPList struct {
	List map[string]IPProperties `json:"ips"`
}

//IP is JSON struct if a single IP
type IP struct {
	Properties IPProperties `json:"ip"`
}

//IPProperties is JSON struct of an IP's properties
type IPProperties struct {
	Name            string      `json:"name"`
	LocationCountry string      `json:"location_country"`
	LocationUUID    string      `json:"location_uuid"`
	ObjectUUID      string      `json:"object_uuid"`
	ReverseDNS      string      `json:"reverse_dns"`
	Family          int         `json:"family"`
	Status          string      `json:"status"`
	CreateTime      string      `json:"create_time"`
	Failover        bool        `json:"failover"`
	ChangeTime      string      `json:"change_time"`
	LocationIata    string      `json:"location_iata"`
	LocationName    string      `json:"location_name"`
	Prefix          string      `json:"prefix"`
	IP              string      `json:"ip"`
	DeleteBlock     string      `json:"delete_block"`
	UsagesInMinutes float64     `json:"usage_in_minutes"`
	CurrentPrice    float64     `json:"current_price"`
	Labels          []string    `json:"labels"`
	Relations       IPRelations `json:"relations"`
}

//IPRelations is JSON struct of a list of an IP's relations
type IPRelations struct {
	Loadbalancers []IPLoadbalancer                  `json:"loadbalancers"`
	Servers       []IPServer                        `json:"servers"`
	PublicIPs     []ServerIPRelationProperties      `json:"public_ips"`
	Storages      []ServerStorageRelationProperties `json:"storages"`
}

//IPLoadbalancer is JSON struct of the relation between an IP and a Load Balancer
type IPLoadbalancer struct {
	CreateTime       string `json:"create_time"`
	LoadbalancerName string `json:"loadbalancer_name"`
	LoadbalancerUUID string `json:"loadbalancer_uuid"`
}

//IPServer is JSON struct of the relation between an IP and a Server
type IPServer struct {
	CreateTime string `json:"create_time"`
	ServerName string `json:"server_name"`
	ServerUUID string `json:"server_uuid"`
}

//IPCreateResponse is JSON struct of a response for creating an IP
type IPCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
	Prefix      string `json:"prefix"`
	IP          string `json:"ip"`
}

//IPCreateRequest is JSON struct of a request for creating an IP
type IPCreateRequest struct {
	Name         string   `json:"name,omitempty"`
	Family       int      `json:"family"`
	LocationUUID string   `json:"location_uuid"`
	Failover     bool     `json:"failover,omitempty"`
	ReverseDNS   string   `json:"reverse_dns,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

//IPUpdateRequest is JSON struct of a request for updating an IP
type IPUpdateRequest struct {
	Name       string   `json:"name,omitempty"`
	Failover   bool     `json:"failover"`
	ReverseDNS string   `json:"reverse_dns,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

//IPEventList is JSON struct of a list of an IP's events
type IPEventList struct {
	List []IPEventProperties `json:"events"`
}

//IPEvent is JSON struct of a single IP
type IPEvent struct {
	Properties IPEventProperties `json:"event"`
}

//IPEventProperties is JSON struct of an IP's properties
type IPEventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUUID   string `json:"request_uuid"`
	ObjectUUID    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUUID      string `json:"user_uuid"`
}

//GetIP get a specific IP based on given id
func (c *Client) GetIP(id string) (IP, error) {
	r := Request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodGet,
	}

	var response IP
	err := r.execute(*c, &response)

	return response, err
}

//GetIPList gets a list of available IPs
func (c *Client) GetIPList() ([]IP, error) {
	r := Request{
		uri:    apiIPBase,
		method: http.MethodGet,
	}

	var response IPList
	var IPs []IP
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}

	return IPs, err
}

//CreateIP creates an IP
func (c *Client) CreateIP(body IPCreateRequest) (IPCreateResponse, error) {
	r := Request{
		uri:    apiIPBase,
		method: http.MethodPost,
		body:   body,
	}

	var response IPCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return IPCreateResponse{}, err
	}

	err = c.WaitForRequestCompletion(response.RequestUUID)

	return response, err
}

//DeleteIP deletes a specific IP based on given id
func (c *Client) DeleteIP(id string) error {
	r := Request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodDelete,
	}

	return r.execute(*c, nil)
}

//UpdateIP updates a specific IP based on given id
func (c *Client) UpdateIP(id string, body IPUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodPatch,
		body:   body,
	}

	return r.execute(*c, nil)
}

//GetIPEventList gets a list of an IP's events
func (c *Client) GetIPEventList(id string) ([]IPEvent, error) {
	r := Request{
		uri:    path.Join(apiIPBase, id, "events"),
		method: http.MethodGet,
	}
	var response IPEventList
	var IPEvents []IPEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		IPEvents = append(IPEvents, IPEvent{Properties: properties})
	}
	return IPEvents, err
}

//GetIPVersion gets IP's version, returns 0 if an error was encountered
func (c *Client) GetIPVersion(id string) int {
	ip, err := c.GetIP(id)
	if err != nil {
		return 0
	}
	return ip.Properties.Family
}
