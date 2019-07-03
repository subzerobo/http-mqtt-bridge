# Http-MQTT-Bridge

The HTTP to MQTT bridge should feel that gap of IFTTT Actions for your Custom IoT Hardwares.
The idea is to receive signals using HTTP requests and transfer them to your MQTT broker. The HTTP to MQTT bridge is written using Golang with Gin for HTTP server and Paho MQTT client.

#Usage
This app could be hosted on Heroku

###Running on Local Machine

* Clone the Project `git clone git@github.com:subzerobo/http-mqtt-bridge.git`
* `cd http-mqtt-bridge`
* `go get ./...`
* `go build -o hmb`
* `hmb s --mqtt-host=tcp://localhost:1883`

###CLI Help
* `hmb s --help`

###Configuration Environment Variables

* MQTT Host : `MQTT_HOST`
* MQTT User : `MQTT_USER`
* MQTT Pass : `MQTT_PASS`

* Basic Auth Username : `AUTH_USERNAME`
* Basic Auth Password : `AUTH_PASSWORD`

###Test

`curl -XPOST -H 'Content-Type: application/json' -d '{"topic": "ali","message": "This is a Test Message"}' --user alikaviani:F#@{fW+/ localhost:8090/publish`

This will publish message: `this is a test message` to the topic : `ali`



 