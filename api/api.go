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

// ErrEmptyLocation is returned when an empty location string is given for an endpoint
type ErrEmptyLocation struct {
}

func (e ErrEmptyLocation) Error() string {
	return "Empty location given"
}

// ErrOriginEqualsDestination is returned when a matching origin and destination are provided for an endpoint
type ErrOriginEqualsDestination struct {
	location string
}

func (e ErrOriginEqualsDestination) Error() string {
	return fmt.Sprintf("Origin location is equal destination (%s)", e.location)
}

// ErrAuthenticationFailed is returned when API credentials aren't accepted
type ErrAuthenticationFailed struct {
}

func (e ErrAuthenticationFailed) Error() string {
	return "API Authentication failed"
}

// User contains data for a RTT API account, wrapping requests
type User struct {
	Username        string
	Password        string
	SearchEndpoint  *url.URL
	ServiceEndpoint *url.URL
	Client          *http.Client
}

// New creates a new user login for RTT
func New(username, password string, baseURL *url.URL, client *http.Client) (User, error) {

	// create the search endpoint from the base URL
	searchEndpoint, err := baseURL.Parse("/search/")
	if err != nil {
		return User{}, err
	}

	// create the service endpoint from the base URL
	serviceEndpoint, err := baseURL.Parse("/service/")
	return User{
		Username:        username,
		Password:        password,
		SearchEndpoint:  searchEndpoint,
		ServiceEndpoint: serviceEndpoint,
		Client:          client,
	}, err
}

func (c User) get(u *url.URL) (*http.Response, error) {

	// setup the basic GET request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Add authentication
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	// send the request to the API
	resp, err := c.Client.Do(req)
	if err != nil {
		return resp, err
	}

	// check the response status code, return custom error
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrAuthenticationFailed{}
	default:
		return resp, err
	}
}

// Departures returns all of the departures from a starting station
func (c User) Departures(origin string) (lineup model.Lineup, err error) {

	// get the URL for this request
	url, err := getDepartures(c.SearchEndpoint, origin)
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
func getDepartures(endpoint *url.URL, origin string) (*url.URL, error) {
	if origin == "" {
		return nil, ErrEmptyLocation{}
	}
	return endpoint.Parse(origin)
}

// DeparturesToDestination returns all of the departures from one station to another
func (c User) DeparturesToDestination(origin, destination string) (lineup model.Lineup, err error) {

	// get the URL for this request
	url, err := getDeparturesDestination(c.SearchEndpoint, origin, destination)
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
func getDeparturesDestination(endpoint *url.URL, origin, destination string) (*url.URL, error) {

	// checking for dodgy input data
	switch {
	case origin == "":
		return nil, ErrEmptyLocation{}
	case destination == "":
		return nil, ErrEmptyLocation{}
	case origin == destination:
		return nil, ErrOriginEqualsDestination{location: origin}
	}

	// append path data to endpoint
	paths := []string{origin, "to", destination}
	ext := strings.Join(paths, "/")
	return endpoint.Parse(ext)
}

// ServicesForDate returns all of the services on a given day
func (c User) ServicesForDate(origin string, date time.Time) (lineup model.Lineup, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServicesDate(c.SearchEndpoint, origin, date)
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
func getServicesDate(endpoint *url.URL, origin string, date time.Time) (*url.URL, error) {

	// checking for dodgy input data
	if origin == "" {
		return nil, ErrEmptyLocation{}
	}

	// append path data to endpoint
	paths := []string{
		origin,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
	}
	ext := strings.Join(paths, "/")
	return endpoint.Parse(ext)
}

// ServicesForTime returns all the services ot a given time
func (c User) ServicesForTime(origin string, date time.Time) (lineup model.Lineup, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServicesTime(c.SearchEndpoint, origin, date)
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
func getServicesTime(endpoint *url.URL, origin string, date time.Time) (*url.URL, error) {

	// checking for dodgy input data
	if origin == "" {
		return nil, ErrEmptyLocation{}
	}

	// append path data to endpoint
	paths := []string{
		origin,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
		fmt.Sprintf("%02d%02d", date.Hour(), date.Minute()),
	}
	ext := strings.Join(paths, "/")
	return endpoint.Parse(ext)
}

// ServiceInfo returns information about a specific service id
func (c User) ServiceInfo(id string, date time.Time) (service model.Service, err error) {

	// send the get request for the custom resource endpoint
	url, err := getServiceInfo(c.ServiceEndpoint, id, date)
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

func getServiceInfo(endpoint *url.URL, service string, date time.Time) (*url.URL, error) {

	// checking for dodgy input data
	if service == "" {
		return nil, ErrEmptyLocation{}
	}

	// append path data to endpoint
	paths := []string{
		service,
		strconv.Itoa(date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
		fmt.Sprintf("%02d%02d", date.Hour(), date.Minute()),
	}
	ext := strings.Join(paths, "/")
	return endpoint.Parse(ext)
}
