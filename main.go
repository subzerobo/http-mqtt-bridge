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
}

type Container struct {
	Config Configuration
	MqttClient mqtt.Client
}

type PublishDTO struct {
	Topic string `json="topic" binding:"required"`
	Message string `json="message" binding:"required"`
}

func main()  {
	
	app := cli.NewApp()
	app.Name = "HTTP-MQTT Bridge"
	app.Usage = "Makes a bridge from HTTP to MQTT"
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}
	
	app.Commands = []cli.Command {
		{
			Name: "start",
			Aliases: []string{"s"},
			Usage: "starts the server",
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "mqtt, m",
					Value : "tcp://localhost:1883",
					Usage : "MQTT Server Address",
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
		MqttServerURL: c.String("mqtt"),
	}
	return container
}

func Setup(container *Container)  {
	// Connect to MQTT Client
	SetupMQTT(container)
	
	// Setup GIN
	r := SetupGin(container)
	
	// Run Gin Server
	r.Run(":8090")
}

func SetupMQTT(container *Container) {
	opts := mqtt.NewClientOptions().AddBroker(container.Config.MqttServerURL)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	container.MqttClient = client
}

func SetupGin(container *Container) *gin.Engine  {
	r := gin.Default()
	
	// Ping Test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"alikaviani":  "F#@{fW+/", // user:foo password:bar
	}))
	
	// Post To Topic
	authorized.POST("publish", func(c *gin.Context) {
		var publishDTO PublishDTO
		if err := c.BindJSON(&publishDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		}
		// Go For MQTT Publish
		client := container.MqttClient
		if token := client.Publish(publishDTO.Topic,0,false, publishDTO.Message); token.Wait() && token.Error() != nil {
			log.Println("Error:", token.Error())
			c.JSON(http.StatusFailedDependency, gin.H{"status": "error", "error": token.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	return r
	
}