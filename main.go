package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicolas-nannoni/fingy-server/events"
	"github.com/nicolas-nannoni/fingy-server/services"
	"log"
	"time"
)

func main() {

	router := gin.Default()
	router.GET("/device/:deviceId/socket", socket)
	registerServices()

	go reg.run()
	go router.Run()

	c := time.NewTimer(time.Second * 10)
	<-c.C

	if err := reg.Send("uuid.NewV1()", &events.Event{Path: "/hello/123"}); err != nil {
		log.Print(err)
	}

	select {}
}

func socket(c *gin.Context) {
	socketHandler(c.Param("deviceId"), c.Writer, c.Request)
}

func registerServices() {
	services.Registry.RegisterService(&services.Service{Id: "alfred", Host: "localhost:8081"})
}
