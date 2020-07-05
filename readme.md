# TeslaCommand

### Overview
A Golang program to connect to a Tesla vehicle, determine if it's within a GeoFence, and once it enters, if it's not plugged in, send an email alert.

Given a teslamotors.com account, interval, coordinates, and radius, it connects to the Tesla RESTful API, and determines the vehicles location and charging state.

The command line args are specified via minus sign, argname, equal sign, value.

You must first go on [Google Maps][1] or [LatLong][2] and get the Longitude and Latitude of the center point of your charging destination.

### Installation
Install [Go][3] and [Git][4] if you don't have them.  Go can be tricky with paths.  Type `$ go env` and then `$ cd` into directory in `$GOROOT/src`.  Clone this repository from within the `src` directory.

`$ git clone https://github.com/morrisonbrett/TeslaCommand.git`

`$ cd TeslaCommand`

Below is an example run.  The long/lat is for the Tesla Hawthorne, CA Supercharger (replace with your own values):

`$ go run TeslaCommand.go -checkInterval=300 -fromAddress="user@gmail.com" -geoFenceLatitude=33.921063 -geoFenceLongitude=-118.33015434 -mailServer="smtp.gmail.com" -mailServerLogin="user@gmail.com" -mailServerPassword="the-gmail-password" -mailServerPort=587 -radius=200 -teslaLoginEmail="user@gmail.com" -teslaLoginPassword="the-teslamotors-password" -toAddress="user@gmail.com" -vehicleIndex=0`

Please see the "Issues" link for a list of "TO DO" items.  It's a work in progress... :-)

[1]: https://support.google.com/maps/answer/18539?hl=en
[2]: http://www.latlong.net/
[3]: https://golang.org/
[4]: http://git-scm.com/download/

################Modifications and fixes
Vehicles will "sleep" after a certain period. When asleep, vehicle data is unavailable so a wake command has been added. 
A hard coded wait time of 60 seconds is used to allow the vehicle to wake. Note that repeated waking of the vehicle drains the battery.

The original used a fee-for-service provider (twilio) to send texts. This code has been commented out, but input parameters remain, are required and
are unused. A second toaddress has been added to allow free email-to-text transmissions.

Changes were made to email transmission. Most SMTP servers require the message to contain TO:, FROM: and SUBJECT: since FROM is often checked for 
validity. Also, port 25 seems to be the only one that works (Gmail, Charter).

The original was designed to loop if the vehicle is found outside the fence. Here we assume that the norm is for the vehicle to be inside the fence when checking occurs
so the primary check is for charging door open. If true, the program quits. If false, it loops. Suggested loop interval 3600 seconds. Program is
designed to be started with a scheduler (e.g.cron or Window Task Manager) at a time when vehicle **should** be charging.

Note that the example command above is **out of date, even for the original**. Here is a sample command for the current version:
go run TeslaCommand.go  -alertThreshold=100 -checkInterval=3600 -fromAddress="user@gmail.com" -geoFenceLatitude=33.921063 -geoFenceLongitude=-118.33015434 -mailServer="mobile.charter.net" -mailServerLogin="user@charter.net" -mailServerPassword="the-password" -mailServerPort=25 -radius=200 -recipientPhoneNumber="7775551212" -senderPhoneNumber="7775551212" -teslaLoginEmail"user@gmail.com" -teslaLoginPassword="the-teslamotors-password" -toAddress1="user@gmail.com" -toAddress2="user@gmail.com" -twilioSID="3334445555" -twilioToken="2223334444" -vehicleIndex=0

