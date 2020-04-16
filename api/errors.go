package api

import (
	"fmt"
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
