package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/georgeprice/realtime-trains-golang/model"
)

const (
	username = "username"
	password = "password"
)

var (
	getDeparturesResponse = model.Lineup{
		Location: model.LocationDetailHeader{
			Name: "getDeparturesResponse",
		},
	}
	getDeparturesDestinationResponse = model.Lineup{
		Location: model.LocationDetailHeader{
			Name: "getDeparturesDestinationResponse",
		},
	}
	getServicesDateResponse = model.Lineup{
		Location: model.LocationDetailHeader{
			Name: "getServicesDateResponse",
		},
	}
	getServicesTimeResponse = model.Lineup{
		Location: model.LocationDetailHeader{
			Name: "getServicesTimeResponse",
		},
	}
	getServiceInfoResponse = model.Service{
		ServiceUID: "getServiceInfoResponse",
	}
)

func mockServer() http.HandlerFunc {

	// searchHandler will write back responses for search endpoint parameters
	searchHandler := func(rw http.ResponseWriter, params ...string) {

		var (
			encodeError  error
			requestError error
			encoder      = json.NewEncoder(rw)
		)

		// perform action based on request type
		switch len(params) {
		case 1:
			encodeError = encoder.Encode(getDeparturesResponse)
		case 3:
			encodeError = encoder.Encode(getDeparturesDestinationResponse)
		case 4:
			encodeError = encoder.Encode(getServicesDateResponse)
		case 5:
			encodeError = encoder.Encode(getServicesTimeResponse)
		default:
			requestError = fmt.Errorf("Search request not recognised w/ params %+v", params)
		}

		// setting response header based on errors
		switch {
		case encodeError != nil:
			rw.WriteHeader(http.StatusInternalServerError)
		case requestError != nil:
			rw.WriteHeader(http.StatusBadRequest)
		}

	}

	// service handler will write back response for service endpoint parameters
	serviceHandler := func(rw http.ResponseWriter, params ...string) {

		var (
			encodeError  error
			requestError error
			encoder      = json.NewEncoder(rw)
		)

		// perform action based on request type
		switch len(params) {
		case 5:
			encodeError = encoder.Encode(getServiceInfoResponse)
		default:
			requestError = fmt.Errorf("Service request not recognised w/ params %+v", params)
		}

		// setting response header based on errors
		switch {
		case encodeError != nil:
			rw.WriteHeader(http.StatusInternalServerError)
		case requestError != nil:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}

	// entry-point handler for requests, delegates to search and service handlers
	return func(rw http.ResponseWriter, req *http.Request) {

		// check authentication
		user, pass, ok := req.BasicAuth()
		if !ok || user != username || pass != password {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		paths := strings.Split(req.URL.Path, "/")
		command := paths[1]
		switch command {
		case "search":
			searchHandler(rw, paths[2:]...)
		case "service":
			serviceHandler(rw, paths[2:]...)
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

func TestAPI(t *testing.T) {
	var (
		client User
		server *httptest.Server
	)

	// closing down the test server
	defer func() {
		switch server {
		case nil:
			t.Fatal("Got nil http test server, cannot close")
		default:
			server.Close()
		}
	}()

	// set up client and server for API testing
	t.Run("setup", func(t *testing.T) {

		// setup mock server, for accepting requests
		server = httptest.NewServer(mockServer())

		// create base URL for requests
		base, err := url.Parse(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		// create client for interacting with mock API
		client, err = New(username, password, base, &http.Client{})
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Departures", func(t *testing.T) {
		response, err := client.Departures("MAN")
		switch {
		case !reflect.DeepEqual(response, getDeparturesResponse):
			t.Fatalf("Got wrong response, got %+v, expected %+v", response, getDeparturesResponse)
		case err != nil:
			t.Fatal(err)
		}
	})

	t.Run("DeparturesToDestination", func(t *testing.T) {
		response, err := client.DeparturesToDestination("MAN", "BHM")
		switch {
		case !reflect.DeepEqual(response, getDeparturesDestinationResponse):
			t.Fatalf("Got wrong response, got %+v, expected %+v", response, getDeparturesDestinationResponse)
		case err != nil:
			t.Fatal(err)
		}
	})

	t.Run("ServicesForDate", func(t *testing.T) {
		response, err := client.ServicesForDate("MAN", time.Now())
		switch {
		case !reflect.DeepEqual(response, getServicesDateResponse):
			t.Fatalf("Got wrong response, got %+v, expected %+v", response, getServicesDateResponse)
		case err != nil:
			t.Fatal(err)
		}
	})

	t.Run("ServicesForTime", func(t *testing.T) {
		response, err := client.ServicesForTime("MAN", time.Now())
		switch {
		case !reflect.DeepEqual(response, getServicesTimeResponse):
			t.Fatalf("Got wrong response, got %+v, expected %+v", response, getServicesTimeResponse)
		case err != nil:
			t.Fatal(err)
		}
	})

	t.Run("ServiceInfo", func(t *testing.T) {
		response, err := client.ServiceInfo("W16631", time.Now())
		switch {
		case !reflect.DeepEqual(response, getServiceInfoResponse):
			t.Fatalf("Got wrong response, got %+v, expected %+v", response, getServiceInfoResponse)
		case err != nil:
			t.Fatal(err)
		}
	})

	t.Run("Errors", func(t *testing.T) {

		t.Run("Bad Base", func(t *testing.T) {

			// create base URL for requests
			badBase := &url.URL{
				Scheme: "fakeScheme!!!",
				Host:   "fakeHost!!!",
			}

			// create client for interacting with mock API
			badClient, err := New(username, password, badBase, &http.Client{})
			if err != nil {
				t.Fatal(err)
			}

			_, err = badClient.get(badBase)
			if err == nil {
				t.Fatal(err)
			}

		})

		t.Run("Cannot Connect", func(t *testing.T) {

			// setup mock server, for accepting requests
			badServer := httptest.NewServer(mockServer())

			badServer.Close()

			// create base URL for requests
			base, err := url.Parse(badServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			// create client for interacting with mock API
			client, err := New(username, password, base, &http.Client{})
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Departures("MAN")
			if err == nil {
				t.Fatal("Got nil error, expected error")
			}
		})

		t.Run("Unauthorized", func(t *testing.T) {

			// create base URL for requests
			base, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			// create client for interacting with mock API
			client, err := New("fake", "fake", base, &http.Client{})
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Departures("MAN")
			if err == nil {
				t.Fatal("Got nil error, expected error")
			}

			if err != ErrAuthenticationFailed {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrAuthenticationFailed)
			}

		})

		t.Run("Departures", func(t *testing.T) {
			_, err := client.Departures("")
			if err != ErrEmptyLocation {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrEmptyLocation)
			}
		})

		t.Run("DeparturesToDestination", func(t *testing.T) {
			_, err := client.DeparturesToDestination("MAN", "")
			if err != ErrEmptyLocation {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrEmptyLocation)
			}
		})

		t.Run("ServicesForDate", func(t *testing.T) {
			_, err := client.ServicesForDate("", time.Now())
			if err != ErrEmptyLocation {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrEmptyLocation)
			}
		})

		t.Run("ServicesForTime", func(t *testing.T) {
			_, err := client.ServicesForTime("", time.Now())
			if err != ErrEmptyLocation {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrEmptyLocation)
			}
		})

		t.Run("ServiceInfo", func(t *testing.T) {
			_, err := client.ServiceInfo("", time.Now())
			if err != ErrEmptyLocation {
				t.Fatalf("Got wrong error, got %+v, expected %+v", err, ErrEmptyLocation)
			}
		})

	})

}

type test struct {

	// inputs
	origin, destination string
	date                time.Time

	// expected outputs
	urlStr string
	err    error
}

func (t test) check(gotURL *url.URL, gotErr error) error {

	if !reflect.DeepEqual(gotErr, t.err) {
		return fmt.Errorf("Test %+v: Got wrong error, got %+v, expected %+v", t, gotErr, t.err)
	}

	var gotStr string
	if gotURL != nil {
		gotStr = gotURL.String()
	}

	if gotStr != t.urlStr {
		return fmt.Errorf("Test %+v: Got wrong URL, got %+v, expected %+v", t, gotURL.String(), t.urlStr)
	}
	return nil
}

func TestURLs(t *testing.T) {

	var (
		searchEndpoint  *url.URL
		serviceEndpoint *url.URL
	)

	const (
		base        = "https://api.rtt.io/api/v1/json/"
		searchBase  = base + "search/"
		serviceBase = base + "service/"
	)

	t.Run("setup", func(t *testing.T) {
		var err error

		// create the search endpoint URL object
		searchEndpoint, err = url.Parse(searchBase)
		if err != nil {
			t.Fatal(err)
		}

		// create the service endpoint URL object
		serviceEndpoint, err = url.Parse(serviceBase)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("getDepartures", func(t *testing.T) {
		ts := []test{
			{
				origin: "MAN",
				urlStr: searchBase + "MAN",
			},
			{
				origin: "BRM",
				urlStr: searchBase + "BRM",
			},
			{
				origin: "",
				err:    ErrEmptyLocation,
			},
		}

		for _, tc := range ts {
			gotURL, err := getDepartures(searchEndpoint, tc.origin)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getDeparturesDestination", func(t *testing.T) {
		ts := []test{
			{
				origin:      "MAN",
				destination: "BRM",
				urlStr:      searchBase + "MAN/to/BRM",
			},
			{
				origin:      "BOMO",
				destination: "BRM",
				urlStr:      searchBase + "BOMO/to/BRM",
			},
			{
				origin:      "MAN",
				destination: "MAN",
				err:         ErrOriginEqualsDestination,
			},
			{
				origin:      "",
				destination: "MAN",
				err:         ErrEmptyLocation,
			},
			{
				origin:      "MAN",
				destination: "",
				err:         ErrEmptyLocation,
			},
			{
				origin:      "",
				destination: "",
				err:         ErrEmptyLocation,
			},
		}

		for _, tc := range ts {
			gotURL, err := getDeparturesDestination(searchEndpoint, tc.origin, tc.destination)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServicesDate", func(t *testing.T) {
		ts := []test{
			{
				origin: "MAN",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: searchBase + "MAN/2020/02/03",
			},
			{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation,
			},
		}

		for _, tc := range ts {
			gotURL, err := getServicesDate(searchEndpoint, tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServicesTime", func(t *testing.T) {
		ts := []test{
			{
				origin: "MAN",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: searchBase + "MAN/2020/02/03/0405",
			},
			{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation,
			},
		}
		for _, tc := range ts {
			gotURL, err := getServicesTime(searchEndpoint, tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServiceInfo", func(t *testing.T) {
		ts := []test{
			{
				origin: "serviceName",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: serviceBase + "serviceName/2020/02/03/0405",
			},
			{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation,
			},
		}
		for _, tc := range ts {
			gotURL, err := getServiceInfo(serviceEndpoint, tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})
}
