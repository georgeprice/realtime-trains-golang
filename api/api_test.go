package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestAPI(t *testing.T) {
	var (
		login  credentials
		client API
	)

	t.Run("load-credentials-file", func(t *testing.T) {

		// get working directory
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		// create path to test credentials JSON file
		p := path.Join(wd, "test-credentials.json")

		// open the file up
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		// unpack the file contents into our login struct
		err = json.NewDecoder(f).Decode(&login)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("create-api-login", func(t *testing.T) {
		var err error
		client, err = New(login.Username, login.Password, &http.Client{Timeout: time.Second})
		switch err {
		case nil:
			t.Logf("Using test credentials {%+v}", login)
		default:
			t.Fatal(err)
		}
	})

	t.Run("GetDepartures", func(t *testing.T) {
		lineup, err := client.GetDepartures("MAN")
		switch err {
		case nil:
			t.Logf("Got lineup %+v", lineup)
		default:
			t.Fatal(err)
		}
	})

	t.Run("GetDeparturesDestination", func(t *testing.T) {
		lineup, err := client.GetDeparturesDestination("MAN", "BHM")
		switch err {
		case nil:
			t.Logf("Got lineup %+v", lineup)
		default:
			t.Fatal(err)
		}
	})

	t.Run("GetServicesDate", func(t *testing.T) {
		lineup, err := client.GetServicesDate("MAN", time.Now())
		switch err {
		case nil:
			t.Logf("Got lineup %+v", lineup)
		default:
			t.Fatal(err)
		}
	})

	t.Run("GetServicesTime", func(t *testing.T) {
		lineup, err := client.GetServicesTime("MAN", time.Now())
		switch err {
		case nil:
			t.Logf("Got lineup %+v", lineup)
		default:
			t.Fatal(err)
		}
	})

	t.Run("GetServiceInfo", func(t *testing.T) {
		service, err := client.GetServicesTime("W16631", time.Now())
		switch err {
		case nil:
			t.Logf("Got service %+v", service)
		default:
			t.Fatal(err)
		}
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
	t.Run("getDepartures", func(t *testing.T) {
		ts := []test{
			test{
				origin: "MAN",
				urlStr: "https://api.rtt.io/api/v1/json/search/MAN",
			},
			test{
				origin: "BRM",
				urlStr: "https://api.rtt.io/api/v1/json/search/BRM",
			},
			test{
				origin: "",
				err:    ErrEmptyLocation{},
			},
		}

		for _, tc := range ts {
			gotURL, err := getDepartures(tc.origin)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getDeparturesDestination", func(t *testing.T) {
		ts := []test{
			test{
				origin:      "MAN",
				destination: "BRM",
				urlStr:      "https://api.rtt.io/api/v1/json/search/MAN/to/BRM",
			},
			test{
				origin:      "BOMO",
				destination: "BRM",
				urlStr:      "https://api.rtt.io/api/v1/json/search/BOMO/to/BRM",
			},
			test{
				origin:      "MAN",
				destination: "MAN",
				err:         ErrOriginEqualsDestination{location: "MAN"},
			},
			test{
				origin:      "",
				destination: "MAN",
				err:         ErrEmptyLocation{},
			},
			test{
				origin:      "MAN",
				destination: "",
				err:         ErrEmptyLocation{},
			},
			test{
				origin:      "",
				destination: "",
				err:         ErrEmptyLocation{},
			},
		}

		for _, tc := range ts {
			gotURL, err := getDeparturesDestination(tc.origin, tc.destination)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServicesDate", func(t *testing.T) {
		ts := []test{
			test{
				origin: "MAN",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: "https://api.rtt.io/api/v1/json/search/MAN/2020/02/03",
			},
			test{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation{},
			},
		}

		for _, tc := range ts {
			gotURL, err := getServicesDate(tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServicesTime", func(t *testing.T) {
		ts := []test{
			test{
				origin: "MAN",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: "https://api.rtt.io/api/v1/json/search/MAN/2020/02/03/0405",
			},
			test{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation{},
			},
		}
		for _, tc := range ts {
			gotURL, err := getServicesTime(tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})

	t.Run("getServiceInfo", func(t *testing.T) {
		ts := []test{
			test{
				origin: "serviceName",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				urlStr: "https://api.rtt.io/api/v1/json/service/serviceName/2020/02/03/0405",
			},
			test{
				origin: "",
				date:   time.Date(2020, 2, 3, 4, 5, 6, 0, &time.Location{}),
				err:    ErrEmptyLocation{},
			},
		}
		for _, tc := range ts {
			gotURL, err := getServiceInfo(tc.origin, tc.date)
			err = tc.check(gotURL, err)
			if err != nil {
				t.Error(err.Error())
			}
		}

	})
}
