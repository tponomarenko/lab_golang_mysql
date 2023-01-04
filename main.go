package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuthentication(settings *Settings) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authentication")

		if authHeader != fmt.Sprintf("Token %s", settings.AuthToken) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}

}

func main() {
	settings, err := NewSettings()
	if err != nil {
		return
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	if settings.AuthToken != "" {
		engine.Use(TokenAuthentication(settings))
	}

	endpoint, err := NewEndpoint(settings)
	if err != nil {
		return
	}

	engine.GET("/records/", endpoint.GetRecords)
	engine.POST("/records/", endpoint.AddRecord)
	engine.GET("/records/:recordId", endpoint.GetRecordById)
	engine.DELETE("/records/:recordId", endpoint.DeleteRecord)
	engine.PUT("/records/:recordId", endpoint.UpdateRecord)

	_ = engine.Run(fmt.Sprintf("0.0.0.0:%s", settings.ServicePort))
}
