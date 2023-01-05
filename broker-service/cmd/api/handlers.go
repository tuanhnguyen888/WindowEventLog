package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string `json:"action"`
}

func Broker(ctx *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	ctx.JSON(
		http.StatusAccepted,
		gin.H{
			"message": payload,
		},
	)
}

func Handle(ctx *gin.Context) {
	var requestPayload RequestPayload

	err := ctx.BindJSON(&requestPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	switch requestPayload.Action {
	case "saveDbPg":
		SavePdDb(ctx)
	case "showDb":
		Show(ctx)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknow Action"})
	}
}

func SavePdDb(ctx *gin.Context) {
	message := struct {
		Message string
	}{}

	jsonData, _ := json.MarshalIndent(message, "", "\t")
	logServiceURL := "http://savelog-service/savePg"

	request, err := http.NewRequest("GET", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dont new request"})
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Add success",
	}

	ctx.JSON(
		http.StatusAccepted,
		gin.H{
			"message": payload,
		},
	)
}

func Show(ctx *gin.Context) {
	message := struct {
		Data any
	}{}

	jsonData, _ := json.MarshalIndent(message, "", "\t")
	logServiceURL := "http://savelog-service/find"

	request, err := http.NewRequest("GET", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dont new request"})
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = json.NewDecoder(response.Body).Decode(&message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "show data",
		Data:    message.Data,
	}

	ctx.JSON(
		http.StatusAccepted,
		gin.H{
			"message": payload,
		},
	)
}
