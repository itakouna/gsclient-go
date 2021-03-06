package gsclient

import (
	"net/http"
	"path"
)

//SshkeyList JSON struct of a list of SSH-keys
type SshkeyList struct {
	List map[string]SshkeyProperties `json:"sshkeys"`
}

//Sshkey JSON struct of a single SSH-key
type Sshkey struct {
	Properties SshkeyProperties `json:"sshkey"`
}

//SshkeyProperties JSON struct of properties of a single SSH-key
type SshkeyProperties struct {
	Name       string   `json:"name"`
	ObjectUUID string   `json:"object_uuid"`
	Status     string   `json:"status"`
	CreateTime string   `json:"create_time"`
	ChangeTime string   `json:"change_time"`
	Sshkey     string   `json:"sshkey"`
	Labels     []string `json:"labels"`
	UserUUID   string   `json:"user_uuid"`
}

//SshkeyCreateRequest JSON struct of a request for creating a SSH-key
type SshkeyCreateRequest struct {
	Name   string   `json:"name"`
	Sshkey string   `json:"sshkey"`
	Labels []string `json:"labels,omitempty"`
}

//SshkeyUpdateRequest JSON struct of a request for updating a SSH-key
type SshkeyUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Sshkey string   `json:"sshkey,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//SshkeyEventList JSON struct of a list of a SSH-key's events
type SshkeyEventList struct {
	List []SshkeyEventProperties `json:"events"`
}

//SshkeyEvent JSON struct of an event of a SSH-key
type SshkeyEvent struct {
	Properties SshkeyEventProperties `json:"event"`
}

//SshkeyEventProperties JSON struct of properties of an event of a SSH-key
type SshkeyEventProperties struct {
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

//GetSshkey gets a ssh key
func (c *Client) GetSshkey(id string) (Sshkey, error) {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodGet,
	}
	var response Sshkey
	err := r.execute(*c, &response)
	return response, err
}

//GetSshkeyList gets a list of ssh keys
func (c *Client) GetSshkeyList() ([]Sshkey, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: http.MethodGet,
	}

	var response SshkeyList
	var sshKeys []Sshkey
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		sshKeys = append(sshKeys, Sshkey{Properties: properties})
	}
	return sshKeys, err
}

//CreateSshkey creates a ssh key
func (c *Client) CreateSshkey(body SshkeyCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: "POST",
		body:   body,
	}
	var response CreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return CreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//DeleteSshkey deletes a ssh key
func (c *Client) DeleteSshkey(id string) error {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateSshkey updates a ssh key
func (c *Client) UpdateSshkey(id string, body SshkeyUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetSshkeyEventList gets a ssh key's events
func (c *Client) GetSshkeyEventList(id string) ([]SshkeyEvent, error) {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id, "events"),
		method: http.MethodGet,
	}
	var response SshkeyEventList
	var sshEvents []SshkeyEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		sshEvents = append(sshEvents, SshkeyEvent{Properties: properties})
	}
	return sshEvents, err
}
