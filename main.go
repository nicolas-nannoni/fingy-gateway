package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/service/:serviceId/device/:deviceId/socket", socket)
	registerServices()

	go router.Run()
	go SetupServiceSideGateway()

	select {}
}

func socket(c *gin.Context) {
	socketHandler(c.Param("serviceId"), c.Param("deviceId"), c.Writer, c.Request)
}

func registerServices() {
	Registry.RegisterService(&Service{Id: "alfred", Host: "localhost:8092"})
}
