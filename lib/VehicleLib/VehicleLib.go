// Package VehicleLib ...
//
// Brett Morrison, Februrary 2016
//
// A library to support Tesla vehicle commands
//
// API Documented here: http://docs.timdorr.apiary.io/#
//
package VehicleLib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://owner-api.teslamotors.com/"

// LoginInfo ...
type LoginInfo struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	Expires   string `json:"expires_in"`
	Created   string `json:"created_at"`
}

// VehicleInfo ...
type VehicleInfo struct {
	DisplayName string `json:"display_name"`
	ID          int    `json:"id"`
	OptionCodes string `json:"option_codes"`
	UserID      int    `json:"user_id"`
	VehicleID   int    `json:"vehicle_id"`
	Vin         string `json:"vin"`
	State       string `json:"state"`
}

// VehicleInfoResponse ...
type VehicleInfoResponse struct {
	Vehicles []VehicleInfo `json:"response"`
	Count    int           `json:"count"`
}

// VehicleLocation ...
type VehicleLocation struct {
	ShiftState string  `json:"shift_state"`
	Speed      int     `json:"speed"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Heading    int     `json:"heading"`
	GPSTime    int     `json:"gps_as_of"`
}

// VehicleLocationResponse ...
type VehicleLocationResponse struct {
	VehicleLocation VehicleLocation `json:"response"`
}

// VehicleChargeState ...
type VehicleChargeState struct {
	ChargingState      string `json:"charging_state"`
	ChargeToMaxRange   bool   `json:"charge_to_max_range"`
	ChargePartDoorOpen bool   `json:"charge_port_door_open"`
	BatteryLevel       int    `json:"battery_level"`
}

// VehicleChargeStateResponse ...
type VehicleChargeStateResponse struct {
	VehicleChargeState VehicleChargeState `json:"response"`
}

// ListVehicles ...
func ListVehicles(token string, vir *VehicleInfoResponse) error {
	resource := "api/1/vehicles"

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println(urlStr)

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("listVehicles request error: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("listVehicles response code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&vir); err == io.EOF {
		return fmt.Errorf("listVehicles decode error: %s", err)
	}

	fmt.Printf("Count: %d\n", vir.Count)

	for i := 0; i < vir.Count; i++ {
		fmt.Printf("Name: %v\n", vir.Vehicles[i].DisplayName)
		fmt.Printf("ID: %v\n", vir.Vehicles[i].ID)
		fmt.Printf("VIN: %v\n", vir.Vehicles[i].Vin)
		fmt.Printf("OptionCodes: %v\n", vir.Vehicles[i].OptionCodes)
		fmt.Printf("State: %v\n", vir.Vehicles[i].State)
		fmt.Println()
	}

	return nil
}

// GetChargeState ...
func GetChargeState(token string, id int, vcsr *VehicleChargeStateResponse) error {
	resource := fmt.Sprintf("api/1/vehicles/%d/data_request/charge_state", id)

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println(urlStr)

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("GetChargeState request error: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GetChargeState response code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&vcsr); err == io.EOF {
		return fmt.Errorf("GetChargeState decode error: %s", err)
	}

	fmt.Printf("ChargingState: %v\n", vcsr.VehicleChargeState.ChargingState)
	fmt.Printf("BatteryLevel: %v\n", vcsr.VehicleChargeState.BatteryLevel)
	fmt.Printf("ChargeToMaxRange: %v\n", vcsr.VehicleChargeState.ChargeToMaxRange)
	fmt.Printf("ChargePartDoorOpen: %v\n", vcsr.VehicleChargeState.ChargePartDoorOpen)
	fmt.Println()

	return nil
}

// GetLocation ...
func GetLocation(token string, id int, vlr *VehicleLocationResponse) error {
	resource := fmt.Sprintf("api/1/vehicles/%d/data_request/drive_state", id)

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println(urlStr)

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("getLocation request error: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("getLocation response code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&vlr); err == io.EOF {
		return fmt.Errorf("getLocation decode error: %s", err)
	}

	fmt.Printf("ShiftState: %v\n", vlr.VehicleLocation.ShiftState)
	fmt.Printf("Speed: %v\n", vlr.VehicleLocation.Speed)
	fmt.Printf("Latitude: %v\n", vlr.VehicleLocation.Latitude)
	fmt.Printf("Longitude: %v\n", vlr.VehicleLocation.Longitude)
	fmt.Printf("Heading: %v\n", vlr.VehicleLocation.Heading)
	fmt.Printf("GPS Time: %v\n", vlr.VehicleLocation.GPSTime)
	fmt.Println()

	return nil
}

// TeslaLogin ...
func TeslaLogin(clientid string, clientsecret string, email string, password string, li *LoginInfo) error {
	resource := "/oauth/token"

	data := url.Values{}
	data.Add("grant_type", "password")
	data.Add("client_id", clientid)
	data.Add("client_secret", clientsecret)
	data.Add("email", email)
	data.Add("password", password)

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println(urlStr)
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("URL: %v\n", urlStr)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("teslaLogin request error: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("teslaLogin response code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&li); err == io.EOF {
		return fmt.Errorf("teslaLogin decode error: %s", err)
	}

	return nil
}
