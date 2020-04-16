# Realtime Trains Golang
Golang implementation of Realtime Train's API

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

### Interface

```go
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
```

### Basic usage
Create a struct to hold login credentials, use methods to get data.
```go
// create a struct to hold your credentials
user := api.New("username", "password", http.Client{ /* ... */ })

// getting departures...
lineup, err := user.GetDepartures("MAN")
lineup, err = user.GetDeparturesDestination("MAN", "BRM")
lineup, err = user.GetServicesDate("MAN", time.Now())
lineup, err = user.GetServicesTime("MAN", time.Now())

// getting service info...
service, err := user.GetServiceInfo("W16631", time.Now())

```
