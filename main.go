package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Airport struct {
	CountryIsoName string  `json:"country_iso_name"`
	CountryName    string  `json:"country_name"`
	Elevation      float64 `json:"elevation"`
	IataCode       string  `json:"iata_code"`
	IcaoCode       string  `json:"icao_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Municipality   string  `json:"municipality"`
	Name           string  `json:"name"`
}

type FlightRoute struct {
	Callsign    string  `json:"callsign"`
	Origin      Airport `json:"origin"`
	Destination Airport `json:"destination"`
	Midpoint    Airport `json:"midpoint"`
}

type ApiResponse struct {
	Response struct {
		FlightRoute FlightRoute `json:"flightroute"`
	} `json:"response"`
}

func main() {
	for {
		var input string
		fmt.Printf("1: To get info on a plane\n")
		fmt.Printf("exit\n")
		_, err := fmt.Scanln(&input)
		if err != nil {
			return
		}
		switch input {
		case "1":
			fmt.Print("Enter a call sign: ")
			_, err := fmt.Scanln(&input)
			if err != nil {
				return
			}
			planeInfo, err := fetchPlaneInfo(input)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Ground Information for callsign %s:\n", input)
			fmt.Printf("Origin Airport: %s, %s\n", planeInfo.Origin.Name, planeInfo.Origin.CountryName)
			fmt.Printf("Destination Airport: %s, %s\n", planeInfo.Destination.Name, planeInfo.Destination.CountryName)
			if planeInfo.Midpoint.Name != "" {
				fmt.Printf("Midpoint Airport: %s, %s\n", planeInfo.Midpoint.Name, planeInfo.Midpoint.CountryName)
			}
		case "exit":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid input")
		}
	}
}
func fetchPlaneInfo(callsign string) (*FlightRoute, error) {
	url := fmt.Sprintf("https://api.adsbdb.com/v0/callsign/%s", callsign)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse.Response.FlightRoute, nil
}

//func fetchAirportInfo(airportCode string) (*FlightRoute, error) {

//
