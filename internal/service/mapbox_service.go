package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dcarbon/go-shared/gutils"
	"github.com/Dcarbon/go-shared/libs/utils"
)

const (
	baseURL = "https://api.mapbox.com/geocoding/v5/mapbox.places/"
)

// Struct to unmarshal the JSON response
type GeocodingResponse struct {
	Features []struct {
		PlaceName string `json:"place_name"`
	} `json:"features"`
}
type MapService struct {
	AccessToken string
}

func NewMapService(config *gutils.Config,
) (*MapService, error) {
	mappService := MapService{
		AccessToken: utils.StringEnv("MAPBOX_ACCESS_TOKEN", ""),
	}
	fmt.Println("log: ", mappService.AccessToken)
	return &mappService, nil
}

func (m *MapService) GetAddress(latitude, longitude float64) (string, error) {
	// Build the request URL with latitude, longitude, access token, and desired format
	url := fmt.Sprintf("%s%f,%f.json?access_token=%s", baseURL, longitude, latitude, m.AccessToken)

	// Make the HTTP request
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Parse the JSON response
	var geocodingResponse GeocodingResponse
	if err := json.NewDecoder(response.Body).Decode(&geocodingResponse); err != nil {
		return "", err
	}

	// Check if there are any results
	if len(geocodingResponse.Features) == 0 {
		return "", fmt.Errorf("no address found for the provided coordinates")
	}

	// Extract and return the formatted address
	return geocodingResponse.Features[0].PlaceName, nil
}
