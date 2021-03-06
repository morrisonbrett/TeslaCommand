//
// Brett Morrison, Februrary 2016
//
// A program to alert if Tesla is not plugged in at a specified GeoFence
//
package main

import (
	"TeslaCommand/lib/HaversinFormula"
	"TeslaCommand/lib/NotifyLib"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	tesla "github.com/jsgoecke/tesla"
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
var twilioSID string
var twilioToken string
var senderPhoneNumber string
var recipientPhoneNumber string
var radius int
var checkInterval int
var alertThreshold int

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
	flag.StringVar(&twilioSID, "twilioSID", "", "Twilio SID")
	flag.StringVar(&twilioToken, "twilioToken", "", "Twilio Token")
	flag.StringVar(&senderPhoneNumber, "senderPhoneNumber", "", "Sender Phone Number")
	flag.StringVar(&recipientPhoneNumber, "recipientPhoneNumber", "", "Recipient Phone Number")
	flag.IntVar(&radius, "radius", 0, "Radius in meters from center geoFence - Typically use 200")
	flag.IntVar(&checkInterval, "checkInterval", 300, "Time in seconds between checks for geoFence")
	flag.IntVar(&alertThreshold, "alertThreshold", 50, "Percentage charged threshold to send alert. If charge level is above threshold, alert won't be sent")
}

func main() {
	// Setup Logging
	logFileName := fmt.Sprintf("TeslaCommand-%v.log", time.Now().Unix())
	logf, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	defer logf.Close()

	log.SetOutput(logf)
	logger := log.New(io.MultiWriter(logf, os.Stdout), "TeslaCommand: ", log.Lshortfile|log.LstdFlags)

	logger.Printf("Num Args %d\n", len(os.Args))
	flag.Parse()
	if len(os.Args) != 19 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := tesla.NewClient(&tesla.Auth{ClientID: clientid, ClientSecret: clientsecret, Email: teslaLoginEmail, Password: teslaLoginPassword})
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}
	//logger.Println("token: " + li.Token)

	vehicles, err := client.Vehicles()
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	vehicle := vehicles[vehicleIndex]

	vehicleState, err := vehicle.VehicleState()
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	// Need to set this flag for every time vehicle exits and enters GeoFence (so we don't send repeated alerts)
	ingeofenceandstopped := false
	waitmessage := fmt.Sprintf("Waiting to check vehicle %v location for %v seconds...\n", vehicleState.VehicleName, checkInterval)

	// Loop every N seconds.
	logger.Printf(waitmessage)
	for _ = range time.Tick(time.Duration(checkInterval) * time.Second) {
		logger.Printf("Checking vehicle %v location after waiting %v seconds.\n", vehicleState.VehicleName, checkInterval)

		driveState, err := vehicle.DriveState()
		if err != nil {
			logger.Println(err)
			logger.Printf(waitmessage)
			continue
		}

		distance := HaversinFormula.Distance(geoFenceLatitude, geoFenceLongitude, driveState.Latitude, driveState.Longitude)

		logger.Printf("Distance: %v\n\n", distance)

		// If the distance is outside the radius, that means vehicle is outside the GeoFence.  Ok to get out
		if distance > float64(radius) {
			ingeofenceandstopped = false
			logger.Printf(waitmessage)
			continue
		}

		// This is to prevent the below logic, if it's already executed, no need to keep doing it
		if ingeofenceandstopped == true {
			logger.Printf(waitmessage)
			continue
		}

		// Check if the vehicle is stopped.
		if (driveState.ShiftState == nil || driveState.ShiftState == "P") && driveState.Speed == 0 {
			ingeofenceandstopped = true

			// In the GeoFence and stopped.  Get the vehicle charge state.
			chargeState, err := vehicle.ChargeState()
			if err != nil {
				logger.Println(err)
				logger.Printf(waitmessage)
				continue
			}

			// Check if the vehicle is plugged in.
			logger.Printf("Vehicle %v is within %v meters of GeoFence with a battery level of %v percent and charging state of %v.", vehicleState.VehicleName, int(distance), chargeState.BatteryLevel, chargeState.ChargingState)
			if chargeState.ChargingState == "Disconnected" {
				// Disconnected, stopped, and within the radius - send alert
				if chargeState.BatteryLevel >= alertThreshold {
					logger.Printf("Battery level is above %v alert threshold. Not sending alert.", chargeState.BatteryLevel)
					logger.Printf(waitmessage)
					continue
				}

				subject := "Tesla Command - " + vehicleState.VehicleName
				body := fmt.Sprintf("Vehicle %v is within %v meters of GeoFence with a battery level of %v percent.  Please plug in.", vehicleState.VehicleName, int(distance), chargeState.BatteryLevel)
				logger.Println(body)

				// Send mail
				err = NotifyLib.SendMail(logger, mailServer, mailServerPort, mailServerLogin, mailServerPassword, fromAddress, toAddress, subject, body)
				if err != nil {
					logger.Println(err)
				}

				// Send text
				err = NotifyLib.SendText(logger, twilioSID, twilioToken, senderPhoneNumber, recipientPhoneNumber, body)
				if err != nil {
					logger.Println(err)
				}
			}
		}
		logger.Printf(waitmessage)
	}
}
