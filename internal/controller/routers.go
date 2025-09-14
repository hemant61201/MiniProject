package controller

import (
	"MiniProject/internal/status"
	"MiniProject/internal/storage/sqlite"
	"MiniProject/internal/types"
	"MiniProject/internal/utiles"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterDevice(router *gin.Engine, storage **sqlite.Sqlite) {

	router.POST(utiles.DEVICE_URL, func(context *gin.Context) {

		var device types.Device

		if err := context.ShouldBindJSON(&device); err != nil {
			ErrorResult(err, http.StatusBadRequest, context)
		}

		status.CheckStatus(&device)

		_, err := (*storage).RegisterDevice(&device)

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		SuccessResult(utiles.DEVICE_REGISTERED, http.StatusOK, context)
	})
}

func GetDeviceListResult(router *gin.Engine, storage **sqlite.Sqlite) {

	router.GET(utiles.DEVICE_URL, func(context *gin.Context) {

		result, err := (*storage).GetDeviceList()

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		SuccessDeviceResult(result, http.StatusOK, context)
	})
}

func GetDevice(router *gin.Engine, storage **sqlite.Sqlite) {

	router.GET(utiles.DEVICE_PARAM_URL, func(context *gin.Context) {

		id, err := strconv.ParseInt(context.Param(utiles.DEVICE_ID), utiles.BASE_VALUE, utiles.BIT_SIZE)

		if err != nil {
			ErrorResult(err, http.StatusBadRequest, context)
		}

		result, err := (*storage).GetDevice(id)

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		SuccessDeviceResult(result, http.StatusOK, context)
	})
}

func UpdateDevice(router *gin.Engine, storage **sqlite.Sqlite) {

	router.PUT(utiles.DEVICE_PARAM_URL, func(context *gin.Context) {

		id, err := strconv.ParseInt(context.Param(utiles.DEVICE_ID), utiles.BASE_VALUE, utiles.BIT_SIZE)

		if err != nil {
			ErrorResult(err, http.StatusBadRequest, context)
		}

		idPresent, err := (*storage).CheckDevice(id)

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		if idPresent {
			UpdateDeviceValue(id, &storage, context)
		} else {
			SuccessResult(utiles.DEVICE_NOT_FOUND, http.StatusOK, context)
		}
	})
}

func UpdateDeviceValue(id int64, storage ***sqlite.Sqlite, context *gin.Context) {

	var deviceInput types.UpdateDeviceInput

	if err := context.ShouldBindJSON(&deviceInput); err != nil {
		ErrorResult(err, http.StatusBadRequest, context)
	}

	result, err := (**storage).UpdateDevice(id, &deviceInput)

	if err != nil {
		ErrorResult(err, http.StatusInternalServerError, context)
	}

	if result == 0 {
		SuccessResult(utiles.DEVICE_NOT_UPDATED, http.StatusOK, context)
	} else {
		SuccessResult(utiles.DEVICE_UPDATED, http.StatusOK, context)
	}
}

func DeleteDevice(router *gin.Engine, storage **sqlite.Sqlite) {

	router.DELETE(utiles.DEVICE_PARAM_URL, func(context *gin.Context) {

		id, err := strconv.ParseInt(context.Param(utiles.DEVICE_ID), utiles.BASE_VALUE, utiles.BIT_SIZE)

		if err != nil {
			ErrorResult(err, http.StatusBadRequest, context)
		}

		result, err := (*storage).DeleteDevice(id)

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		if result == 0 {
			SuccessResult(utiles.DEVICE_NOT_FOUND, http.StatusOK, context)
		} else {
			SuccessResult(utiles.DEVICE_DELETED, http.StatusOK, context)
		}
	})
}

func GetMonitoringResult(router *gin.Engine, storage **sqlite.Sqlite) {

	router.GET(utiles.DEVICE_MONITORING_PARAM_URL, func(context *gin.Context) {

		id, err := strconv.ParseInt(context.Param(utiles.DEVICE_ID), utiles.BASE_VALUE, utiles.BIT_SIZE)

		if err != nil {
			ErrorResult(err, http.StatusBadRequest, context)
		}

		result, err := (*storage).GetDevice(id)

		if err != nil {
			ErrorResult(err, http.StatusInternalServerError, context)
		}

		for i := range result {
			result[i].Metadata.CPU = fmt.Sprintf("%d%%", rand.Intn(100))
			result[i].Metadata.Memory = fmt.Sprintf("%dMB", rand.Intn(16000))
			result[i].Metadata.Disk = fmt.Sprintf("%dGB", rand.Intn(512))
		}

		SuccessDeviceResult(result, http.StatusOK, context)
	})
}
