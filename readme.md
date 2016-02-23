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

`$ go run TeslaCommand.go -checkInterval=5 -fromAddress="user@gmail.com" -geoFenceLatitude=33.921063 -geoFenceLongitude=-118.33015434 -mailServer="smtp.gmail.com" -mailServerLogin="user@gmail.com" -mailServerPassword="the-gmail-password" -mailServerPort=587 -radius=200 -teslaLoginEmail="user@gmail.com" -teslaLoginPassword="the-teslamotors-password" -toAddress="user@gmail.com" -vehicleIndex=0`

Please see the "Issues" link for a list of "TO DO" items.  It's a work in progress... :-)

[1]: https://support.google.com/maps/answer/18539?hl=en
[2]: http://www.latlong.net/
[3]: https://golang.org/
[4]: http://git-scm.com/download/
