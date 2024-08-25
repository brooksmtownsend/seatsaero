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

const API_BASE_URL = "https://seats.aero/partnerapi"

// Should probably be less, this * 4 is the whole trip (2 ways, 2 people)
const POINT_THRESHOLD = 75000

func main() {
	if len(os.Args) < 2 {
		log.Fatal("{\"error\": \"Expected one of 'gettrip' or 'search'\"}")
	}

	args := os.Args[1]

	var response string
	switch args {
	case "gettrip":
		response = string(getTrip(os.Args[2]))
	case "search":
		trips := cachedSearch()
		var wg sync.WaitGroup
		results := make([]TripBooking, len(trips))
		for i, trip := range trips {
			wg.Add(1)
			go func(i int, trip string) {
				defer wg.Done()
				result := getTrip(trip)
				var rr RouteResponse
				err := json.Unmarshal(result, &rr)
				checkError(err)
				results[i] = TripBooking{Trips: rr.Data, Bookings: rr.BookingLinks}
			}(i, trip)
		}
		wg.Wait()
		out, err := json.Marshal(results)
		checkError(err)
		response = string(out)
	default:
		response = "{\"error\": \"Invalid command, expected one of 'gettrip' or 'search'\"}"
	}

	fmt.Println(response)
}

func getTrip(id string) []byte {
	url := fmt.Sprintf("%s/trips/%s", API_BASE_URL, id)
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Partner-Authorization", os.Getenv("SEATS_AERO_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	checkError(err)

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	checkError(err)

	return resp
}

func cachedSearch() []string {
	queryParams := "origin_airport=USA%2CDCA%2CBWI&destination_airport=SEL%2CJPN&cabin=business&start_date=2025-04-01&end_date=2025-05-31&take=1000&order_by=lowest_mileage"
	url := fmt.Sprintf("%s/search?%s", API_BASE_URL, queryParams)
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Partner-Authorization", os.Getenv("SEATS_AERO_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	checkError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	checkError(err)

	var response SearchResponse
	err = json.Unmarshal(body, &response)
	checkError(err)

	// TODO: search fall too
	// TODO: look for more results if we're under the threshold in the first 1000
	// TODO: establish tiers of trips, preferring low direct mileage, then low, then just results
	// TODO: all kinds of filtering

	var tripIds []string

	for _, availability := range response.Data {
		mileageCost, err := strconv.Atoi(availability.JMileageCost)
		checkError(err)

		directMileageCost := availability.JDirectMileageCost
		if (mileageCost > 0 && mileageCost <= POINT_THRESHOLD && availability.JRemainingSeats >= 2) || (directMileageCost > 0 && directMileageCost <= POINT_THRESHOLD && availability.JDirectRemainingSeats >= 2) {
			tripIds = append(tripIds, availability.ID)
		}

		// Just a circuit breaker to avoid too many requests
		if mileageCost > POINT_THRESHOLD && directMileageCost > POINT_THRESHOLD {
			break
		}
	}

	return tripIds
}

// Helper function to check for the error and output to JSON in case of jq use
func checkError(err error) {
	if err != nil {
		log.Fatalf("{\"error\": \"%s\"}", err)
	}
}
