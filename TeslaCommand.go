//
// Brett Morrison, Februrary 2016
//
// A program to alert if Tesla is not plugged in at a specified GeoFence
//
package main

import (
	"TeslaCommand/lib/HaversinFormula"
	"TeslaCommand/lib/VehicleLib"
	"flag"
	"fmt"
	"os"
	"time"
)

// Magic clientid and clientsecret available here: http://pastebin.com/fX6ejAHd
const clientid = "e4a9949fcfa04068f59abb5a658f2bac0a3428e4652315490b659d5ab3f35a9e"
const clientsecret = "c75f14bbadc8bee3a7594412c31416f8300256d7668ea7e6e7f06727bfb9d220"

var teslaLoginEmail string
var teslaLoginPassword string
var vehicleIndex int
var geoFenceLatitude float64
var geoFenceLongitude float64
var mailServer string
var mailServerPort int
var mailServerLogin string
var mailServerPassword string
var fromAddress string
var toAddress string
var radius int
var checkInterval int

func init() {
	flag.StringVar(&teslaLoginEmail, "teslaLoginEmail", "", "Email for teslamotors.com account")
	flag.StringVar(&teslaLoginPassword, "teslaLoginPassword", "", "Password for teslamotors.com account")
	flag.IntVar(&vehicleIndex, "vehicleIndex", 0, "Index of vehicles in your account - If just 1 vehicle, use 0")
	flag.Float64Var(&geoFenceLatitude, "geoFenceLatitude", 0.0, "Latitude of GeoFence Center")
	flag.Float64Var(&geoFenceLongitude, "geoFenceLongitude", 0.0, "Longitude of GeoFence Center")
	flag.StringVar(&mailServer, "mailServer", "", "SMTP Server hostname")
	flag.IntVar(&mailServerPort, "mailServerPort", 25, "SMTP Server port number")
	flag.StringVar(&mailServerLogin, "mailServerLogin", "", "SMTP Server login username")
	flag.StringVar(&mailServerPassword, "mailServerPassword", "", "SMTP Server password")
	flag.StringVar(&fromAddress, "fromAddress", "", "Alert send from email")
	flag.StringVar(&toAddress, "toAddress", "", "Alert send to email")
	flag.IntVar(&radius, "radius", 0, "Radius in meters from center geoFence - Typically use 200")
	flag.IntVar(&checkInterval, "checkInterval", 5, "Time in minutes between checks for geoFence")
}

func main() {
	fmt.Printf("Num Args %d\n", len(os.Args))
	flag.Parse()
	if len(os.Args) != 14 {
		flag.Usage()
		os.Exit(1)
	}

	var li VehicleLib.LoginInfo
	err := VehicleLib.TeslaLogin(clientid, clientsecret, teslaLoginEmail, teslaLoginPassword, &li)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("token: " + li.Token)

	var vir VehicleLib.VehicleInfoResponse
	err = VehicleLib.ListVehicles(li.Token, &vir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Need to set this flag for every time vehicle exits and enters GeoFence (so we don't send repeated alerts)
	ingeofence := false

	// Loop every N minutes.
	fmt.Printf("Waiting to check vehicle %v location for %v minutes...\n", vir.Vehicles[vehicleIndex].DisplayName, checkInterval)
	for _ = range time.Tick(time.Duration(checkInterval) * time.Minute) {
		fmt.Printf("Checking vehicle %v location after waiting %v minutes.\n", vir.Vehicles[vehicleIndex].DisplayName, checkInterval)

		var vlr VehicleLib.VehicleLocationResponse
		err = VehicleLib.GetLocation(li.Token, vir.Vehicles[vehicleIndex].ID, &vlr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		distance := HaversinFormula.Distance(geoFenceLatitude, geoFenceLongitude, vlr.VehicleLocation.Latitude, vlr.VehicleLocation.Longitude)

		fmt.Printf("Distance: %v\n\n", distance)

		// If the distance is outside the radius, that means vehicle is outside the GeoFence.  Ok to get out
		if distance > float64(radius) {
			ingeofence = false
		}

		// This is to prevent the below logic, if it's already executed, no need to keep doing it
		if ingeofence == true {
			continue
		}

		// In the GeoFence.  Get the vehicle charge state.
		var vcsr VehicleLib.VehicleChargeStateResponse
		err = VehicleLib.GetChargeState(li.Token, vir.Vehicles[vehicleIndex].ID, &vcsr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if the vehicle is stopped.
		if vlr.VehicleLocation.Speed == 0 {
			//Check if the vehicle is plugged in.
			if vcsr.VehicleChargeState.ChargingState == "Disconnected" {
				// Disconnected, stopped, and within the radius - send alert
				ingeofence = true
				subject := "Tesla Command - " + vir.Vehicles[vehicleIndex].DisplayName
				body := fmt.Sprintf("Vehicle is within %v meters of GeoFence with a battery level of %v percent.  Please plug in.", int(distance), vcsr.VehicleChargeState.BatteryLevel)
				fmt.Println(body)
				err = VehicleLib.SendMail(mailServer, mailServerPort, mailServerLogin, mailServerPassword, fromAddress, toAddress, subject, body)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		fmt.Printf("Waiting to check vehicle %v location for %v minutes...\n", vir.Vehicles[vehicleIndex].DisplayName, checkInterval)
	}
}
