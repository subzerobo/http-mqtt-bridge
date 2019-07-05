# Http-MQTT-Bridge  
  
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/fb288a8d7dbd4ab9b8ca5d02657fd972)](https://app.codacy.com/app/subzerobo/http-mqtt-bridge?utm_source=github.com&utm_medium=referral&utm_content=subzerobo/http-mqtt-bridge&utm_campaign=Badge_Grade_Dashboard)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)   [![Build Status](https://travis-ci.org/subzerobo/http-mqtt-bridge.svg?branch=master)](https://travis-ci.org/subzerobo/http-mqtt-bridge)  [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/subzerobo/http-mqtt-bridge/issues)  
  
  
The HTTP to MQTT bridge should feel that gap of IFTTT Actions for your Custom IoT Hardwares.  
The idea is to receive signals using HTTP requests and transfer them to your MQTT broker. The HTTP to MQTT bridge is written using Golang with Gin for HTTP server and Paho MQTT client.  
  
# Usage  
This app could be hosted on Heroku  
  
### Running on Local Machine  
  
 * Clone the Project `git clone git@github.com:subzerobo/http-mqtt-bridge.git` * `cd http-mqtt-bridge`  
 * `go get ./...`  
 * `go build -o hmb`  
 * `hmb s --mqtt-host=tcp://localhost:1883`  
  
### CLI Help  
 * `hmb s --help`  
  
### Configuration Environment Variables  
  
 * MQTT Host : `MQTT_HOST`   
 * MQTT User : `MQTT_USER`   
 * MQTT Pass : `MQTT_PASS`      
 * Basic Auth Username : `AUTH_USERNAME`   
 * Basic AuthPassword : `AUTH_PASSWORD`
  
### Test  
  
`curl -XPOST -H 'Content-Type: application/json' -d '{"topic": "ali","message": "This is a Test Message"}' --user alikaviani:F#@{fW+/ localhost:8090/publish`  
  
This will publish message: `this is a test message` to the topic : `ali`  
  
### Deploying to Heroku  
  
You can add required environment variable through .env File accordingly or after creating the heroku app   
Navigate to your app then go to Settings -> Config Vars -> Add Environment variables  
  
```  
$ git clone git@github.com:subzerobo/http-mqtt-bridge.git  
$ cd  http-mqtt-bridge  
$ heroku create  
$ git push heroku master  
$ heroku open  
```  
  OR  
    
 [![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/subzerobo/http-mqtt-bridge)