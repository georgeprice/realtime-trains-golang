package model

import (
	"reflect"
)

// ServiceType describes what type of train is running
type ServiceType string

// BusService is a bus?
// ShipService is a ship?
// TrainService is a train? (Maybe they're all trains, idk)
const (
	BusService   ServiceType = "bus"
	ShipService  ServiceType = "ship"
	TrainService ServiceType = "train"
)

// Service defines the response from a request for a particular train running between two points
type Service struct {
	ServiceUID           string `json:"serviceUid,omitempty"`
	RunDate              string `json:"runDate,omitempty"`
	ServiceType          `json:"serviceType,omitempty"`
	IsPassenger          bool             `json:"isPassenger,omitempty"`
	TrainIdentity        string           `json:"trainIdentity,omitempty"`
	PowerType            string           `json:"powerType,omitempty"`
	TrainClass           string           `json:"trainClass,omitempty"`
	Sleeper              string           `json:"sleeper,omitempty"`
	ATOCCode             string           `json:"atocCode,omitempty"`
	ATOCName             string           `json:"atocName,omitempty"`
	PerformanceMonitored bool             `json:"performanceMonitored,omitempty"`
	Origin               []Pair           `json:"origin,omitempty"`
	Destination          []Pair           `json:"destination,omitempty"`
	Locations            []LocationDetail `json:"locations,omitempty"`
	RealtimeActivated    bool             `json:"realtimeActivated,omitempty"`
	RunningIdentity      string           `json:"runningIdentity,omitempty"`
	PlannedCancel        bool             `json:"plannedCancel,omitempty"`
}

// FresherThan returns whether another service struct contains more up-to-date data
//	I just assume that different data is newer data lol
func (s Service) FresherThan(thou Service) bool {
	return !reflect.DeepEqual(s, thou)
}

// Pair describes a start or end of a train's journey (don't ask)
type Pair struct {
	TIPLOC      string `json:"tiploc,omitempty"`
	Description string `json:"description,omitempty"`
	WorkingTime string `json:"workingTime,omitempty"`
	PublicTime  string `json:"publicTime,omitempty"`
}

// LocationDetail describes a station which is passed through between the origin and destination of a service
type LocationDetail struct {
	RealTimeActivated bool   `json:"realtimeActivated,omitempty"`
	TIPLOC            string `json:"tiploc,omitempty"`
	CRS               string `json:"crs,omitempty"`
	Description       string `json:"description,omitempty"`

	WTTBookedArrival           string `json:"wttBookedArrival,omitempty"`
	WTTBookedDeparture         string `json:"wttBookedDeparture,omitempty"`
	WTTBookedPass              string `json:"wttBookedPass,omitempty"`
	GBTTBookedArrival          string `json:"gbttBookedArrival,omitempty"`
	GBTTBookedArrivalNextDay   bool   `json:"gbttBookedArrivalNextDay,omitempty"`
	GBTTBookedDeparture        string `json:"gbttBookedDeparture,omitempty"`
	GBTTBookedDepartureNextDay bool   `json:"gbttBookedDepartureNextDay,omitempty"`

	Origin       []Pair `json:"origin,omitempty"`
	Destination  []Pair `json:"destination,omitempty"`
	IsCall       bool   `json:"isCall,omitempty"`
	IsCallPublic bool   `json:"isPublicCall,omitempty"`

	RealTimeArrival         string `json:"realtimeArrival,omitempty"`
	RealTimeArrivalActual   bool   `json:"realtimeArrivalActual,omitempty"`
	RealTimeArrivalNoReport bool   `json:"realtimeArrivalNoReport,omitempty"`
	RealTimeArrivalNextDay  bool   `json:"realtimeArrivalNextDay,omitempty"`

	RealTimeGBTTArrivalLateness        int `json:"realtimeGbttArrivalLateness,omitempty"`
	RealTimeWTTArrivalLateness         int `json:"realtimeWttDepartureLateness,omitempty"`
	RealTimeWTTArrivalLatenessDetailed int `json:"realtimeWttArrivalLateness,omitempty"`

	RealTimeDeparture         string `json:"realtimeDeparture,omitempty"`
	RealTimeDepartureActual   bool   `json:"realtimeDepartureActual,omitempty"`
	RealTimeDepartureNoReport bool   `json:"realtimeDepartureNoReport,omitempty"`
	RealTimeDepartureNextDay  bool   `json:"realtimeDepartureNextDay,omitempty"`

	RealTimePass         string `json:"realtimePass,omitempty"`
	RealTimePassActual   bool   `json:"realtimePassActual,omitempty"`
	RealTimePassNoReport bool   `json:"realtimePassNoReport,omitempty"`

	Platform              string `json:"platform,omitempty"`
	PlatformConfirmed     bool   `json:"platformConfirmed,omitempty"`
	PlatformChanged       bool   `json:"platformChanged,omitempty"`
	Line                  string `json:"line,omitempty"`
	LineConfirmed         bool   `json:"lineConfirmed,omitempty"`
	Path                  string `json:"path,omitempty"`
	PathConfirmed         bool   `json:"pathConfirmed,omitempty"`
	CancelReasonCode      string `json:"cancelReasonCode,omitempty"`
	CancelReasonShortText string `json:"cancelReasonShortText,omitempty"`
	CancelReasonLongText  string `json:"cancelReasonLongText,omitempty"`
	DisplayAs             string `json:"displayAs,omitempty"`
	ServiceLocation       string `json:"serviceLocation,omitempty"`
}
