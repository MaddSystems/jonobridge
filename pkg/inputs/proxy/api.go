package main

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SendCommandRequest struct {
	IMEI string `json:"imei" binding:"required"`
	Data string `json:"data" binding:"required"` // Hex encoded data
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	
	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/trackerlist", func(c *gin.Context) {
			connMutex.Lock()
			trackers := make([]TrackerAssign, 0)
			for imei, connInfo := range imeiConnections {
				tracker := TrackerAssign{
					Imei:       imei,
					Protocol:   connInfo.Protocol,
					RemoteAddr: connInfo.RemoteAddr,
				}
				trackers = append(trackers, tracker)
			}
			connMutex.Unlock()
			c.JSON(http.StatusOK, trackers)
		})

		v1.POST("/sendcommand", func(c *gin.Context) {
			var req SendCommandRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			// Check if IMEI exists
			connMutex.Lock()
			connInfo, exists := imeiConnections[req.IMEI]
			connMutex.Unlock()

			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Tracker not found"})
				return
			}

			// Convert hex string to bytes
			data, err := hex.DecodeString(req.Data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hex data"})
				return
			}

			// Send data to connection
			err = SendDataToConnection(connInfo.RemoteAddr, data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Command sent successfully"})
		})
	}

	return r
}
