#!/bin/bash

source .env

FROM=2025-09-01
TO=2025-10-30
CABIN=business

# Just a sample request to test connectivity
curl \
    --header "Partner-Authorization: $SEATS_AERO_API_KEY" \
    --header "accept: application/json" \
    "https://seats.aero/partnerapi/search?cabin=${CABIN}&start_date=${FROM}&end_date=${TO}&origin_airport=USA%2CDCA%2CBWI&destination_airport=SEL%2CJPN&take=1000&order_by=lowest_mileage"