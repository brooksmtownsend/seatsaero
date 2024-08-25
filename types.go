package main

import "time"

// seats.aero types

type SearchResponse struct {
	Data    []Availability `json:"data"`
	Count   int            `json:"count"`
	HasMore bool           `json:"hasMore"`
	Cursor  int            `json:"cursor"`
}

type Availability struct {
	ID                    string    `json:"ID"`
	RouteID               string    `json:"RouteID"`
	Route                 Route     `json:"Route"`
	Date                  string    `json:"Date"`
	ParsedDate            time.Time `json:"ParsedDate"`
	YAvailable            bool      `json:"YAvailable"` // Economy
	WAvailable            bool      `json:"WAvailable"` // Premium Economy
	JAvailable            bool      `json:"JAvailable"` // Business
	FAvailable            bool      `json:"FAvailable"` // First
	YMileageCost          string    `json:"YMileageCost"`
	WMileageCost          string    `json:"WMileageCost"`
	JMileageCost          string    `json:"JMileageCost"`
	FMileageCost          string    `json:"FMileageCost"`
	YDirectMileageCost    int       `json:"YDirectMileageCost"`
	WDirectMileageCost    int       `json:"WDirectMileageCost"`
	JDirectMileageCost    int       `json:"JDirectMileageCost"`
	FDirectMileageCost    int       `json:"FDirectMileageCost"`
	YRemainingSeats       int       `json:"YRemainingSeats"`
	WRemainingSeats       int       `json:"WRemainingSeats"`
	JRemainingSeats       int       `json:"JRemainingSeats"`
	FRemainingSeats       int       `json:"FRemainingSeats"`
	YDirectRemainingSeats int       `json:"YDirectRemainingSeats"`
	WDirectRemainingSeats int       `json:"WDirectRemainingSeats"`
	JDirectRemainingSeats int       `json:"JDirectRemainingSeats"`
	FDirectRemainingSeats int       `json:"FDirectRemainingSeats"`
	YAirlines             string    `json:"YAirlines"`
	WAirlines             string    `json:"WAirlines"`
	JAirlines             string    `json:"JAirlines"`
	FAirlines             string    `json:"FAirlines"`
	YDirectAirlines       string    `json:"YDirectAirlines"`
	WDirectAirlines       string    `json:"WDirectAirlines"`
	JDirectAirlines       string    `json:"JDirectAirlines"`
	FDirectAirlines       string    `json:"FDirectAirlines"`
	YDirect               bool      `json:"YDirect"`
	WDirect               bool      `json:"WDirect"`
	JDirect               bool      `json:"JDirect"`
	FDirect               bool      `json:"FDirect"`
	Source                string    `json:"Source"`
	CreatedAt             time.Time `json:"CreatedAt"`
	UpdatedAt             time.Time `json:"UpdatedAt"`
	AvailabilityTrips     *string   `json:"AvailabilityTrips"`
}

type Route struct {
	ID                 string `json:"ID"`
	OriginAirport      string `json:"OriginAirport"`
	OriginRegion       string `json:"OriginRegion"`
	DestinationAirport string `json:"DestinationAirport"`
	DestinationRegion  string `json:"DestinationRegion"`
	NumDaysOut         int    `json:"NumDaysOut"`
	Distance           int    `json:"Distance"`
	Source             string `json:"Source"`
}

type RouteResponse struct {
	Data                   []Trip        `json:"data"`
	OriginCoordinates      Coordinates   `json:"origin_coordinates"`
	DestinationCoordinates Coordinates   `json:"destination_coordinates"`
	BookingLinks           []BookingLink `json:"booking_links"`
	RevalidationID         string        `json:"revalidation_id"`
}

type Trip struct {
	ID                   string    `json:"ID"`
	RouteID              string    `json:"RouteID"`
	AvailabilityID       string    `json:"AvailabilityID"`
	AvailabilitySegments []Segment `json:"AvailabilitySegments"`
	TotalDuration        int       `json:"TotalDuration"`
	Stops                int       `json:"Stops"`
	Carriers             string    `json:"Carriers"`
	RemainingSeats       int       `json:"RemainingSeats"`
	MileageCost          int       `json:"MileageCost"`
	TotalTaxes           int       `json:"TotalTaxes"`
	TaxesCurrency        string    `json:"TaxesCurrency"`
	TaxesCurrencySymbol  string    `json:"TaxesCurrencySymbol"`
	AllianceCost         int       `json:"AllianceCost"`
	TotalSegmentDistance int       `json:"TotalSegmentDistance"`
	FlightNumbers        string    `json:"FlightNumbers"`
	DepartsAt            time.Time `json:"DepartsAt"`
	Cabin                string    `json:"Cabin"`
	ArrivesAt            time.Time `json:"ArrivesAt"`
	CreatedAt            time.Time `json:"CreatedAt"`
	UpdatedAt            time.Time `json:"UpdatedAt"`
	Source               string    `json:"Source"`
}

type Segment struct {
	ID                 string    `json:"ID"`
	RouteID            string    `json:"RouteID"`
	AvailabilityID     string    `json:"AvailabilityID"`
	AvailabilityTripID string    `json:"AvailabilityTripID"`
	FlightNumber       string    `json:"FlightNumber"`
	Distance           int       `json:"Distance"`
	FareClass          string    `json:"FareClass"`
	AircraftName       string    `json:"AircraftName"`
	AircraftCode       string    `json:"AircraftCode"`
	OriginAirport      string    `json:"OriginAirport"`
	DestinationAirport string    `json:"DestinationAirport"`
	DepartsAt          time.Time `json:"DepartsAt"`
	ArrivesAt          time.Time `json:"ArrivesAt"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
	Source             string    `json:"Source"`
	Order              int       `json:"Order"`
}

type Coordinates struct {
	Lat float64 `json:"Lat"`
	Lon float64 `json:"Lon"`
}

type BookingLink struct {
	Label   string `json:"label"`
	Link    string `json:"link"`
	Primary bool   `json:"primary"`
}

// Internal use

type TripBooking struct {
	Trips    []Trip        `json:"trips"`
	Bookings []BookingLink `json:"bookings"`
}
