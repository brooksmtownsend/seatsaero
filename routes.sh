#!/bin/bash

source .env

# Just a sample request to test connectivity
curl \
    --header "Partner-Authorization: $SEATS_AERO_API_KEY" \
    --header "accept: application/json" \
    https://seats.aero/partnerapi/routes | jq > routes.json