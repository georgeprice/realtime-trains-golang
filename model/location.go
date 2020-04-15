package model

// Lineup defines the response from a search for services for a given location
type Lineup struct {
	Location LocationDetailHeader `json:"location,omitempty"`
	Filter   []LocationDetail     `json:"filter,omitempty"`
	Services []LocationContainer  `json:"services,omitempty"`
}

// LocationDetailHeader describes the shorthand location used in the query
type LocationDetailHeader struct {
	Name   string `json:"name,omitempty"`
	CRS    string `json:"crs,omitempty"`
	TIPLOC string `json:"tiploc,omitempty"`
}

// LocationContainer contains a description of a service which is running for a lineup
type LocationContainer struct {
	LocationDetail   `json:"locationDetail,omitempty"`
	ServiceUID       string `json:"serviceUid,omitempty"`
	RunDate          string `json:"runDate,omitempty"`
	TrainIdentity    string `json:"trainIdentity,omitempty"`
	RunningIdentity  string `json:"runningIdentity,omitempty"`
	ATOCCode         string `json:"atocCode,omitempty"`
	ATOCName         string `json:"atocName,omitempty"`
	ServiceType      string `json:"serviceType,omitempty"`
	IsPassenger      bool   `json:"isPassenger,omitempty"`
	PlannedCancel    bool   `json:"plannedCancel,omitempty"`
	Origin           []Pair `json:"origin,omitempty"`
	Destination      []Pair `json:"destination,omitempty"`
	CountdownMinutes int    `json:"countdownMinutes,omitempty"`
}
