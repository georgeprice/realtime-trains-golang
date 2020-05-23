package model

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"reflect"
	"testing"
)

const (
	lineupFile  = "lineup.json"
	serviceFile = "service.json"
)

var (
	expectedLineup = Lineup{
		Location: LocationDetailHeader{
			Name: "Bournemouth", CRS: "BMH", TIPLOC: "BOMO",
		},
		Services: []LocationContainer{
			{
				LocationDetail: LocationDetail{
					RealTimeActivated:   true,
					TIPLOC:              "BOMO",
					CRS:                 "BMH",
					Description:         "Bournemouth",
					WTTBookedArrival:    "011630",
					WTTBookedDeparture:  "011830",
					GBTTBookedArrival:   "0117",
					GBTTBookedDeparture: "0118",
					Origin: []Pair{
						{
							TIPLOC:      "WATRLMN",
							Description: "London Waterloo",
							WorkingTime: "230500",
							PublicTime:  "2305",
						},
					},
					Destination: []Pair{
						{
							TIPLOC:      "POOLE",
							Description: "Poole",
							WorkingTime: "013000",
							PublicTime:  "0130",
						},
					},
					IsCall:                  true,
					IsCallPublic:            true,
					RealTimeArrival:         "0114",
					RealTimeArrivalActual:   false,
					RealTimeDeparture:       "0118",
					RealTimeDepartureActual: false,
					Platform:                "3",
					PlatformConfirmed:       false,
					PlatformChanged:         false,
					DisplayAs:               "CALL",
				},
				ServiceUID:      "W90091",
				RunDate:         "2013-06-11",
				TrainIdentity:   "1B77",
				RunningIdentity: "1B77",
				ATOCCode:        "SW",
				ATOCName:        "South West Trains",
				ServiceType:     "train",
				IsPassenger:     true,
			},
		},
	}
	expectedService = Service{
		ServiceUID:           "Q13773",
		RunDate:              "2020-02-12",
		ServiceType:          TrainService,
		IsPassenger:          true,
		TrainIdentity:        "1Z73",
		PowerType:            "EMU",
		TrainClass:           "B",
		ATOCCode:             "SW",
		ATOCName:             "South Western Railway",
		PerformanceMonitored: true,
		Origin: []Pair{
			{
				TIPLOC:      "ELGH",
				Description: "Eastleigh",
				WorkingTime: "233700",
				PublicTime:  "2337",
			},
		},
		Destination: []Pair{
			{
				TIPLOC:      "POOLE",
				Description: "Poole",
				WorkingTime: "004700",
				PublicTime:  "0047",
			},
		},
		Locations: []LocationDetail{

			// first from the JSON
			{
				RealTimeActivated:   true,
				TIPLOC:              "ELGH",
				CRS:                 "ESL",
				Description:         "Eastleigh",
				GBTTBookedDeparture: "2337",
				Origin: []Pair{
					{
						TIPLOC:      "ELGH",
						Description: "Eastleigh",
						WorkingTime: "233700",
						PublicTime:  "2337",
					},
				},
				Destination: []Pair{
					{
						TIPLOC:      "POOLE",
						Description: "Poole",
						WorkingTime: "004700",
						PublicTime:  "0047",
					},
				},
				IsCall:                true,
				IsCallPublic:          true,
				RealTimeDeparture:     "2337",
				RealTimeArrivalActual: false,
				Platform:              "2N",
				PlatformConfirmed:     false,
				PlatformChanged:       false,
				DisplayAs:             "ORIGIN",
			},

			{},
			{},
			{},

			// fifth from the HSON
			{
				RealTimeActivated:   true,
				TIPLOC:              "LYNDHRD",
				CRS:                 "ANF",
				Description:         "Ashurst New Forest",
				GBTTBookedArrival:   "2359",
				GBTTBookedDeparture: "2359",
				Origin: []Pair{
					{
						TIPLOC:      "ELGH",
						Description: "Eastleigh",
						WorkingTime: "233700",
						PublicTime:  "2337",
					},
				},
				Destination: []Pair{
					{
						TIPLOC:      "POOLE",
						Description: "Poole",
						WorkingTime: "004700",
						PublicTime:  "0047",
					},
				},
				IsCall:                  true,
				IsCallPublic:            true,
				RealTimeArrival:         "2359",
				RealTimeArrivalActual:   false,
				RealTimeDeparture:       "2359",
				RealTimeDepartureActual: false,
				DisplayAs:               "CALL",
			},

			{},
			{},
			{},
			{},

			// tenth from the JSON
			{
				RealTimeActivated:          true,
				TIPLOC:                     "CHRISTC",
				CRS:                        "CHR",
				Description:                "Christchurch",
				GBTTBookedArrival:          "0026",
				GBTTBookedArrivalNextDay:   true,
				GBTTBookedDeparture:        "0027",
				GBTTBookedDepartureNextDay: true,
				Origin: []Pair{
					{
						TIPLOC:      "ELGH",
						Description: "Eastleigh",
						WorkingTime: "233700",
						PublicTime:  "2337",
					},
				},
				Destination: []Pair{
					{
						TIPLOC:      "POOLE",
						Description: "Poole",
						WorkingTime: "004700",
						PublicTime:  "0047",
					},
				},
				IsCall:                   true,
				IsCallPublic:             true,
				RealTimeArrival:          "0026",
				RealTimeArrivalActual:    false,
				RealTimeArrivalNextDay:   true,
				RealTimeDeparture:        "0027",
				RealTimeDepartureActual:  false,
				RealTimeDepartureNextDay: true,
				DisplayAs:                "CALL",
			},
		},
		RealtimeActivated: true,
		RunningIdentity:   "1Z73",
		PlannedCancel:     true,
	}
)

func TestLineup(t *testing.T) {

	var jsonReader io.ReadCloser
	t.Run("load-expected-data", func(t *testing.T) {

		var testDir string

		t.Run("directory", func(t *testing.T) {
			pwd, err := os.Getwd()
			if err != nil {
				t.Errorf("Could not load test data, got error %s", err.Error())
			}
			testDir = path.Join(pwd, "expected", lineupFile)
		})

		t.Run("json-file", func(t *testing.T) {
			var err error
			jsonReader, err = os.Open(testDir)
			if err != nil {
				t.Errorf("Could not open line up file, got error %s", err.Error())
			}
		})
	})

	// file cleanup
	defer func() {
		if jsonReader != nil {
			return
		}
		jsonReader.Close()
	}()

	var gotLineup Lineup
	t.Run("decoding", func(t *testing.T) {
		err := json.NewDecoder(jsonReader).Decode(&gotLineup)
		if err != nil {
			t.Errorf("Could not decode lineup JSON from file reader, got error ")
		}
	})

	t.Run("comparison", func(t *testing.T) {
		if !reflect.DeepEqual(gotLineup, expectedLineup) {
			t.Errorf("Lineup struct mismatch\nGot %+v\nExpected %+v", gotLineup, expectedLineup)
		}
	})
}

func TestService(t *testing.T) {

	var jsonReader io.ReadCloser
	t.Run("load-expected-data", func(t *testing.T) {

		var testDir string

		t.Run("directory", func(t *testing.T) {
			pwd, err := os.Getwd()
			if err != nil {
				t.Errorf("Could not load test data, got error %s", err.Error())
			}
			testDir = path.Join(pwd, "expected", serviceFile)
		})

		t.Run("json-file", func(t *testing.T) {
			var err error
			jsonReader, err = os.Open(testDir)
			if err != nil {
				t.Errorf("Could not open line up file, got error %s", err.Error())
			}
		})
	})

	// file cleanup
	defer func() {
		if jsonReader != nil {
			return
		}
		jsonReader.Close()
	}()

	var gotService Service
	t.Run("decoding", func(t *testing.T) {
		err := json.NewDecoder(jsonReader).Decode(&gotService)
		if err != nil {
			t.Errorf("Could not decode service JSON from file reader, got error ")
		}
	})

	t.Run("comparison", func(t *testing.T) {

		t.Run("location-details", func(t *testing.T) {

			locationIndexes := []int{0, 4, 9}

			for _, i := range locationIndexes {

				gotLocation := gotService.Locations[i]
				expectedLocation := expectedService.Locations[i]

				if !reflect.DeepEqual(gotLocation, expectedLocation) {
					t.Errorf("Service location detail %d struct mismatch\nGot %+v\nExpected %+v",
						i, gotLocation, expectedLocation)
				}

			}
		})

		t.Run("other-fields", func(t *testing.T) {
			// super dirty workaround to get around checking locations :'(
			gotService.Locations = []LocationDetail{}
			expectedService.Locations = []LocationDetail{}

			if !reflect.DeepEqual(gotService, expectedService) {
				t.Errorf("Service struct mismatch\nGot %+v\nExpected %+v",
					gotService, expectedService)
			}
		})
	})
}
