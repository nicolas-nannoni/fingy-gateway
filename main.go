package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicolas-nannoni/fingy-gateway/events"
	"log"
	"time"
)

func main() {

	router := gin.Default()
	router.GET("/service/:serviceId/device/:deviceId/socket", socket)
	registerServices()

	go router.Run()
	go SetupServiceSideGateway()

	c := time.NewTimer(time.Second * 10)
	<-c.C

	if err := Registry.SendToDevice("alfred", "uuid.NewV1()", &events.Event{Path: "/hello/123"}); err != nil {
		log.Print(err)
	}

	select {}
}

func socket(c *gin.Context) {
	socketHandler(c.Param("serviceId"), c.Param("deviceId"), c.Writer, c.Request)
}

func registerServices() {
	Registry.RegisterService(&Service{Id: "alfred", Host: "localhost:8081"})
}
