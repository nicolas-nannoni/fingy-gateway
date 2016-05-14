package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicolas-nannoni/fingy-gateway/events"
)

func SetupServiceSideGateway() {

	r := gin.Default()
	r.GET("/", index)
	r.POST("/service/:serviceId/sendEvent/device/:deviceId/*path", sendEventToDevice)

	r.Run(":8090")
}

func index(c *gin.Context) {
	c.Status(200)
}

func sendEventToDevice(c *gin.Context) {

	serviceId := c.Param("serviceId")
	deviceId := c.Param("deviceId")
	path := c.Param("path")

	evt := events.Event{
		ServiceId: serviceId,
		Path:      path,
	}

	Registry.SendToDevice(serviceId, deviceId, &evt)
}
