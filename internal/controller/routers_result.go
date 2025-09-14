package controller

import (
	"MiniProject/internal/types"

	"github.com/gin-gonic/gin"
)

func ErrorResult(err error, urlRequest int, context *gin.Context) {
	context.JSON(urlRequest, gin.H{
		"error": err.Error(),
	})
}

func SuccessResult(message string, urlRequest int, context *gin.Context) {
	context.JSON(urlRequest, gin.H{
		"result": message,
	})
}

func SuccessDeviceResult(deviceInfo []types.DeviceInfo, urlRequest int, context *gin.Context) {
	context.JSON(urlRequest, gin.H{
		"result": deviceInfo,
	})
}
