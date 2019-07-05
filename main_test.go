package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	container *Container
	router *gin.Engine
)

func initContext() *cli.Context {
	set := flag.NewFlagSet("test", 0)
	set.String("mqtt-host", "tcp://localhost:1883", "doc")
	set.String("mqtt-user", "", "doc")
	set.String("mqtt-pass", "", "doc")
	set.String("username", "user1", "doc")
	set.String("password", "pass1", "doc")
	set.String("port", "8080", "doc")
	ctx := cli.NewContext(nil, set, nil)
	return ctx
}

func init() {
	ctx := initContext()
	container = InitializeContainer(ctx)
	router = SetupGin(container)
	
	// Add Mock MQTT Client to Container
	c := newMockClient()
	c.Connect()
	container.MqttClient = c
}

func TestStaticRoute(t *testing.T) {
	
	t.Run("/Ping Static Router", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})
	
	t.Run("/ Static Router", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "i am ready !", w.Body.String())
	})
	
	
	
}

func TestPublishSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	u := PublishDTO{Topic: "topic/test", Message: "t,1,1"}
	req := makePublishRequest("user1", "pass1", u)
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["status"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, "ok", value)
}

func TestPublishFailureWithWrongCredentials(t *testing.T)  {
	w := httptest.NewRecorder()
	u := PublishDTO{Topic: "topic/test", Message: "t,1,1"}
	req := makePublishRequest("wronguser1", "wrongpass1", u)
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPublishFailureWithWrongPayload(t *testing.T)  {
	w := httptest.NewRecorder()
	u := gin.H{}
	req := makePublishRequest("user1", "pass1", u)
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func makePublishRequest(username string, password string, payload interface{}) *http.Request {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(payload)
	req, _ := http.NewRequest("POST", "/publish", b)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(username, password)
	return req
}