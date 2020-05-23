[![Go Report Card](https://goreportcard.com/badge/github.com/georgeprice/realtime-trains-golang)](https://goreportcard.com/report/github.com/georgeprice/realtime-trains-golang)

# Realtime Trains Golang
Golang wrapper of Realtime Train's API

## Model
Each object returned from the API is available as a struct in the __model__ package.

### Sample
```go
var (
    service model.Service
    lineup model.Lineup
)

// create a http request for the api endpoint
request, err := http.NewRequest(http.MethodGet, endpoint, nil)
if err != nil {
    return
}

// add authentication, send the request
req.SetBasicAuth("username", "password")
response, err := http.Do(req)
if err != nil {
    return
}

// create a decoder to unpack the JSON response body bytes into a struct
decoder := json.NewDecoder(response.Body)

// extracting the lineup 
err = decoder.Decode(&lineup)

// or, extracting the service
err = decoder.Decode(&service)

// ...
```

## API

The __API__ package provides an easier way to retrieve data from the Realtime Trains API from your own project.

[Realtime Trains API Docs](https://www.realtimetrains.co.uk/about/developer/pull/docs/)

### Features

```go
// User contains data for a RTT API account, wrapping requests
type User struct {
	Username        string
	Password        string
	SearchEndpoint  *url.URL
	ServiceEndpoint *url.URL
	Client          *http.Client
}

// Departures returns all of the departures from a starting station
func (c User) Departures(origin string) (lineup model.Lineup, err error) {
	// ... 
}

// DeparturesToDestination returns all of the departures from one station to another
func (c User) DeparturesToDestination(origin, destination string) (lineup model.Lineup, err error) {
	// ...
}

// ServicesForDate returns all of the services on a given day
func (c User) ServicesForDate(origin string, date time.Time) (lineup model.Lineup, err error) {
	// ...
}

// ServicesForTime returns all the services ot a given time
func (c User) ServicesForTime(origin string, date time.Time) (lineup model.Lineup, err error) {
	// ...
}

// ServiceInfo returns information about a specific service id
func (c User) ServiceInfo(id string, date time.Time) (service model.Service, err error) {
	// ...
}
```

### Basic usage
Create a struct to hold login credentials, use methods to get data.
```go

// create base URL where RTT API is hosted
apiBase, err := url.Parse(/* ... */)
if err != nil {
	// ...
}

// create a struct to hold your credentials
user := api.New("username", "password", apiBase, http.Client{ /* ... */ })

// getting departures...
lineup, err := user.Departures("MAN")
lineup, err = user.DeparturesToDestination("MAN", "BRM")
lineup, err = user.ServicesForDate("MAN", time.Now())
lineup, err = user.ServicesForTime("MAN", time.Now())

// getting service info...
service, err := user.ServiceInfo("W16631", time.Now())

```
