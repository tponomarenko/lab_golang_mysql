package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	settings, err := NewSettings()
	if err != nil {
		return
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	endpoint, err := NewEndpoint(settings)
	if err != nil {
		return
	}

	engine.GET("/records/", endpoint.GetRecords)
	engine.POST("/records/", endpoint.AddRecord)
	engine.DELETE("/records/:recordId", endpoint.DeleteRecord)
	engine.PUT("/records/:recordId", endpoint.UpdateRecord)

	_ = engine.Run(fmt.Sprintf("0.0.0.0:%s", settings.ServicePort))
}
