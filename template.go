package gsclient

import (
	"fmt"
	"net/http"
	"path"
)

type TemplateList struct {
	List map[string]TemplateProperties `json:"templates"`
}

type Template struct {
	Properties TemplateProperties `json:"template"`
}

type TemplateProperties struct {
	Status           string   `json:"status"`
	Ostype           string   `json:"ostype"`
	LocationUuid     string   `json:"location_uuid"`
	Version          string   `json:"version"`
	LocationIata     string   `json:"location_iata"`
	ChangeTime       string   `json:"change_time"`
	Private          bool     `json:"private"`
	ObjectUuid       string   `json:"object_uuid"`
	LicenseProductNo int      `json:"license_product_no"`
	CreateTime       string   `json:"create_time"`
	UsageInMinutes   int      `json:"usage_in_minutes"`
	Capacity         int      `json:"capacity"`
	LocationName     string   `json:"location_name"`
	Distro           string   `json:"distro"`
	Description      string   `json:"description"`
	CurrentPrice     float64  `json:"current_price"`
	LocationCountry  string   `json:"location_country"`
	Name             string   `json:"name"`
	Labels           []string `json:"labels"`
}

type TemplateEventList struct {
	List []TemplateEventProperties `json:"events"`
}

type TemplateEvent struct {
	Properties TemplateEventProperties `json:"event"`
}

type TemplateEventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUuid   string `json:"request_uuid"`
	ObjectUuid    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUuid      string `json:"user_uuid"`
}

type TemplateCreateRequest struct {
	Name         string   `json:"name"`
	SnapshotUuid string   `json:"snapshot_uuid"`
	Labels       []string `json:"labels"`
}

type TemplateUpdateRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

//GetTemplate gets a template
func (c *Client) GetTemplate(id string) (Template, error) {
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodGet,
	}
	var response Template
	err := r.execute(*c, &response)

	return response, err
}

//GetTemplateList gets a list of templates
func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{
			Properties: properties,
		})
	}
	return templates, err
}

//GetTemplateByName gets a template by its name
func (c *Client) GetTemplateByName(name string) (Template, error) {
	templates, err := c.GetTemplateList()
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return c.GetTemplate(template.Properties.ObjectUuid)
		}
	}

	return Template{}, fmt.Errorf("Template %v not found", name)
}

//CreateTemplate creates a template
func (c *Client) CreateTemplate(body TemplateCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var resonse CreateResponse
	err := r.execute(*c, &resonse)
	return resonse, err
}

//UpdateTemplate updates a template
func (c *Client) UpdateTemplate(id string, body TemplateUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteTemplate deletes a template
func (c *Client) DeleteTemplate(id string) error {
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetTemplateEventList gets a list of a template's events
func (c *Client) GetTemplateEventList(id string) ([]TemplateEvent, error) {
	r := Request{
		uri:    path.Join(apiTemplateBase, id, "events"),
		method: http.MethodGet,
	}
	var response TemplateEventList
	var templateEvents []TemplateEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templateEvents = append(templateEvents, TemplateEvent{Properties:properties})
	}
	return templateEvents, err
}