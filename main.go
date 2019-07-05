package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	MqttServerURL string
	MqttUsername  string
	MqttPassword  string
	BasicAuthUser string
	BasicAuthPass string
	GinPort       string
}

type Container struct {
	Config     Configuration
	MqttClient mqtt.Client
}

type PublishDTO struct {
	Topic   string `binding:"required" json:"topic"`
	Message string `binding:"required" json:"message"`
}

func main() {
	
	app := cli.NewApp()
	app.Name = "HTTP-MQTT Bridge"
	app.Usage = "Makes a bridge from HTTP to MQTT"
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}
	
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "starts the server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "mqtt-host, mh",
					Value:  "tcp://localhost:1883",
					Usage:  "MQTT Host Address",
					EnvVar: "MQTT_HOST",
				},
				cli.StringFlag{
					Name:   "mqtt-user,mu",
					Value:  "",
					Usage:  "MQTT Host Username",
					EnvVar: "MQTT_USER",
				},
				cli.StringFlag{
					Name:   "mqtt-pass,mp",
					Value:  "",
					Usage:  "MQTT Host Password",
					EnvVar: "MQTT_PASS",
				},
				cli.StringFlag{
					Name:   "username,u",
					Value:  "alikaviani",
					Usage:  "Basic Authentication Username",
					EnvVar: "AUTH_USERNAME",
				},
				cli.StringFlag{
					Name:   "password,p",
					Value:  "F#@{fW+/",
					Usage:  "Basic Authentication Password",
					EnvVar: "AUTH_PASSWORD",
				},
				cli.StringFlag{
					Name:   "port",
					Value:  "8090",
					Usage:  "GIN Router Port",
					EnvVar: "PORT",
				},
			},
			Action: func(c *cli.Context) error {
				// Fills the Configuration
				container := InitializeContainer(c)
				// Setup Gin Server + MQTT
				Setup(container)
				return nil
			},
		},
	}
	
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	
}

func InitializeContainer(c *cli.Context) *Container {
	container := new(Container)
	container.Config = Configuration{
		MqttServerURL: c.String("mqtt-host"),
		MqttUsername:  c.String("mqtt-user"),
		MqttPassword:  c.String("mqtt-pass"),
		BasicAuthUser: c.String("username"),
		BasicAuthPass: c.String("password"),
		GinPort:       c.String("port"),
	}
	return container
}

func Setup(container *Container) {
	// Connect to MQTT Client
	SetupMQTT(container)
	
	// Setup GIN
	r := SetupGin(container)
	
	// Run Gin Server
	r.Run(":" + container.Config.GinPort)
}

func SetupMQTT(container *Container) {
	opts := mqtt.NewClientOptions().AddBroker(container.Config.MqttServerURL)
	if container.Config.MqttUsername != "" {
		opts.Username = container.Config.MqttUsername
		opts.Password = container.Config.MqttPassword
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT Error: ", token.Error())
	}
	container.MqttClient = client
}

func SetupGin(container *Container) *gin.Engine {
	r := gin.Default()
	
	// Ping Test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	
	// Ping Test
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "i am ready !")
	})
	
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		container.Config.BasicAuthUser: container.Config.BasicAuthPass, // user:foo password:bar
	}))
	
	// Post To Topic
	authorized.POST("publish", func(c *gin.Context) {
		var publishDTO PublishDTO
		// Validate Payloadd
		if err := c.ShouldBindJSON(&publishDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
			return
		}
		// Go For MQTT Publish
		client := container.MqttClient
		if token := client.Publish(publishDTO.Topic, 0, false, publishDTO.Message); token.Error() != nil {
			// Return Error
			log.Println("Error:", token.Error())
			c.JSON(http.StatusFailedDependency, gin.H{"status": "error", "error": token.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	return r
	
}
