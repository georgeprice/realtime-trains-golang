package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/georgeprice/realtime-trains-golang/model"
)

var (
	searchEndpoint, _  = url.Parse("https://api.rtt.io/api/v1/json/search/")
	serviceEndpoint, _ = url.Parse("https://api.rtt.io/api/v1/json/service/")
)

type user struct {
	username, password string
	client             *http.Client
}

// New creates a new user login for RTT
func New(username, password string, client *http.Client) (API, error) {
	return &user{
		username: username, password: password,
		client: client,
	}, nil
}

// API handles interacting with a RTT REST service
type API interface {
	// /json/search/<station>
	GetDepartures(origin string) (model.Lineup, error)

	// /json/search/<station>/to/<toStation>
	GetDeparturesDestination(origin, destination string) (model.Lineup, error)

	// /json/search/<station>/<year>/<month>/<day>
	GetServicesDate(origin string, date time.Time) (model.Lineup, error)

	// /json/search/<station>/<year>/<month>/<day>/<time>
	GetServicesTime(origin string, date time.Time) (model.Lineup, error)

	// /json/service/<serviceUid>/<year>/<month>/<day>
	GetServiceInfo(service string, date time.Time) (model.Service, error)
}

func (c user) get(u *url.URL) (*http.Response, error) {

	// setup the basic GET request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// add authentication, send the request
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)

	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrAuthenticationFailed{}
	default:
		return resp, err
	}
}

// GetDepartures returns all of the departures from a starting station
func (c user) GetDepartures(origin string) (lineup model.Lineup, err error) {

	// send the get request for the custom resource endpoint
	url, err := searchEndpoint.Parse(origin)
	if err != nil {
		return lineup, err
	}

	// get response and parse out into service
	resp, err := c.get(url)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&lineup)
	}

	return lineup, err
}

// creates the url to access a lineup resource from an origin
func getDepartures(origin string) (*url.URL, error) {
	if origin == "" {
		return nil, ErrEmptyLocation{}
	}
	return searchEndpoint.Parse(origin)
}

// GetDeparturesDestination returns all of the departures from one station to another
func (c user) GetDeparturesDestination(origin, destination string) (lineup model.Lineup, err error) {

	url, err := getDeparturesDestination(origin, destination)
	if err != nil {
		return lineup, err
	}

	// get response and parse out into service
	resp, err := c.get(url)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&lineup)
	}
	return lineup, err
}

// creates the url to access a lineup resource from an origin to a destination
func getDeparturesDestination(origin, destination string) (*url.URL, error) {

	switch {
	case origin == "":
		return nil, ErrEmptyLocation{}
	case destination == "":
		return nil, ErrEmptyLocation{}
	case origin == destination:
		return nil, ErrOriginEqualsDestination{location: origin}
	}

	paths := []string{origin, "to", destination}
	ext := strings.Join(paths, "/")
	return searchEndpoint.Parse(ext)
}

// GetServicesDate returns all of the services on a given day
func (c user) GetServicesDate(origin string, date time.Time) (lineup model.Lineup, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServicesDate(origin, date)
	if err != nil {
		return lineup, err
	}

	// get response and parse out into service
	resp, err := c.get(url)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&lineup)
	}
	return lineup, err
}

// creates a url to access the service resource from an origin station, on a given date
func getServicesDate(origin string, date time.Time) (*url.URL, error) {

	if origin == "" {
		return nil, ErrEmptyLocation{}
	}

	paths := []string{
		origin,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
	}
	ext := strings.Join(paths, "/")
	return searchEndpoint.Parse(ext)
}

// GetServicesTime returns all the services ot a given time
func (c user) GetServicesTime(origin string, date time.Time) (lineup model.Lineup, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServicesTime(origin, date)
	if err != nil {
		return lineup, err
	}

	// get response and parse out into service
	resp, err := c.get(url)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&lineup)
	}
	return lineup, err
}

// creates a url to access the service resource from an origin station, at a given time
func getServicesTime(origin string, date time.Time) (*url.URL, error) {

	if origin == "" {
		return nil, ErrEmptyLocation{}
	}

	paths := []string{
		origin,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
		fmt.Sprintf("%02d%02d", date.Hour(), date.Minute()),
	}
	ext := strings.Join(paths, "/")
	return searchEndpoint.Parse(ext)
}

func (c user) GetServiceInfo(id string, date time.Time) (service model.Service, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServiceInfo(id, date)
	if err != nil {
		return service, err
	}

	// get response and parse out into service
	resp, err := c.get(url)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&service)
	}
	return service, err
}

func getServiceInfo(service string, date time.Time) (*url.URL, error) {
	if service == "" {
		return nil, ErrEmptyLocation{}
	}

	paths := []string{
		service,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
		fmt.Sprintf("%02d%02d", date.Hour(), date.Minute()),
	}
	ext := strings.Join(paths, "/")
	return serviceEndpoint.Parse(ext)
}
