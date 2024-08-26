package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const API_KEY_ENV = "SEATS_AERO_API_KEY"
const API_BASE_URL = "https://seats.aero/partnerapi"

// Looking for flights from USA/WAS to Seoul and Japan
const SEARCH_PARAMS = "origin_airport=USA%2CDCA%2CBWI&destination_airport=SEL%2CJPN&take=1000&order_by=lowest_mileage"

// Agreed upon point thresholds
const PREMIUM_POINT_THRESHOLD = 37500
const BUSINESS_POINT_THRESHOLD = 75000

func main() {
	if os.Getenv(API_KEY_ENV) == "" {
		log.Fatal("{\"context\": \"init\", \"error\": \"Missing API key\"}")
	}

	trips := search("2025-04-01", "2025-05-31")
	// TODO: Enable this when we're a year out
	// fallTrips := search("2025-09-01", "2025-10-31")
	// trips = append(trips, fallTrips...)

	// Fetch all the trips for additional info, concurrently
	var wg sync.WaitGroup
	results := make([]TripBooking, len(trips))
	for i, trip := range trips {
		wg.Add(1)
		go func(i int, trip string) {
			defer wg.Done()
			result := getTrip(trip)
			var rr RouteResponse
			err := json.Unmarshal(result, &rr)
			checkError("unmarshal_route_response", err)
			results[i] = usefulData(rr.Data, rr.BookingLinks)
		}(i, trip)
	}
	wg.Wait()

	tripJson, err := json.Marshal(results)
	checkError("marshal_results", err)
	fmt.Println(string(tripJson))
}

func search(startDate string, endDate string) []string {
	cabins := []string{"premium", "business", "first"}
	availabilities := []Availability{}

	var wg sync.WaitGroup
	for _, cabin := range cabins {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("%s/search?%s&cabin=%s&start_date=%s&end_date=%s", API_BASE_URL, SEARCH_PARAMS, cabin, startDate, endDate)
			body, err := query(url)
			checkError("cached_search", err)

			var response SearchResponse
			err = json.Unmarshal(body, &response)
			checkError("unmarshal_cached_search", err)
			availabilities = append(availabilities, response.Data...)
		}()
	}

	wg.Wait()

	// TODO: search fall too
	// TODO: look for more results if we're under the threshold in the first 1000
	// TODO: establish tiers of trips, preferring low direct mileage, then low, then just results
	// TODO: all kinds of filtering

	var tripIds []string

	for _, availability := range availabilities {
		if meetsPremiumCriteria(availability) || meetsBusinessCriteria(availability) || meetsFirstCriteria(availability) {
			tripIds = append(tripIds, availability.ID)
		}

		// Just a circuit breaker to avoid continuing the search if we're past the threshold already
		// if mileageCost > BUSINESS_POINT_THRESHOLD && directMileageCost > BUSINESS_POINT_THRESHOLD {
		// 	break
		// }
	}

	// TODO: search for return options as well

	return tripIds
}

// Retrieve an individual trip
func getTrip(id string) []byte {
	url := fmt.Sprintf("%s/trips/%s", API_BASE_URL, id)
	response, err := query(url)
	checkError("get_trip", err)

	return response
}

// Helper function to make a query to the API
func query(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]byte, 0), err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Partner-Authorization", os.Getenv(API_KEY_ENV))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return make([]byte, 0), err
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return make([]byte, 0), err
	}

	if res.StatusCode != 200 {
		return make([]byte, 0), fmt.Errorf("error querying API (%d): %s", res.StatusCode, resp)
	}

	return resp, nil
}

// Helper function to check if the availability meets "worth it" criteria
func meetsPremiumCriteria(a Availability) bool {
	directMileageCost := a.WDirectMileageCost
	mileageCost, err := strconv.Atoi(a.WMileageCost)
	checkError("convert_premium_criteria", err)
	return (mileageCost > 0 && mileageCost <= PREMIUM_POINT_THRESHOLD && a.WRemainingSeats >= 2) || (directMileageCost > 0 && directMileageCost <= PREMIUM_POINT_THRESHOLD && a.WDirectRemainingSeats >= 2)
}

// Helper function to check if the availability meets "worth it" criteria
func meetsBusinessCriteria(a Availability) bool {
	directMileageCost := a.JDirectMileageCost
	mileageCost, err := strconv.Atoi(a.JMileageCost)
	checkError("convert_business_criteria", err)
	return (mileageCost > 0 && mileageCost <= BUSINESS_POINT_THRESHOLD && a.JRemainingSeats >= 2) || (directMileageCost > 0 && directMileageCost <= BUSINESS_POINT_THRESHOLD && a.JDirectRemainingSeats >= 2)
}

// Helper function to check if the availability meets "worth it" criteria
func meetsFirstCriteria(a Availability) bool {
	directMileageCost := a.FDirectMileageCost
	mileageCost, err := strconv.Atoi(a.FMileageCost)
	checkError("first_business_criteria", err)
	return (mileageCost > 0 && mileageCost <= BUSINESS_POINT_THRESHOLD && a.FRemainingSeats >= 2) || (directMileageCost > 0 && directMileageCost <= BUSINESS_POINT_THRESHOLD && a.FDirectRemainingSeats >= 2)
}

// Helper function to check for the error and output to JSON in case of jq use
func checkError(ctx string, err error) {
	if err != nil {
		log.Fatalf("{\"context\": \"%s\"\"error\": \"%s\"}", ctx, err)
	}
}

// Short transform to pull a subset of the data included in the response
func usefulData(trips []Trip, bookings []BookingLink) TripBooking {
	minimalTrips := make([]MinimalTrip, 0)
	for _, trip := range trips {
		if trip.Cabin == "economy" {
			continue
		}
		minimalTrips = append(minimalTrips, MinimalTrip{
			ID:             trip.ID,
			RemainingSeats: trip.RemainingSeats,
			Cabin:          trip.Cabin,
			DepartsAt:      trip.DepartsAt,
			ArrivesAt:      trip.ArrivesAt,
			Stops:          trip.Stops,
			MileageCost:    trip.MileageCost,
			TotalTaxes:     trip.TotalTaxes,
			Source:         trip.Source,
		})
	}

	return TripBooking{Trips: minimalTrips, Bookings: bookings}
}
