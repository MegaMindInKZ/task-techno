package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendInternalServerErrorMessage(res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)

	output := struct {
		Message string `json:"message"`
	}{
		Message: "Internal Server Error",
	}

	jsonResp, err := json.Marshal(output)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func SendErrorMessage(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)

	output := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	jsonResp, err := json.Marshal(output)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func SendMessage(res http.ResponseWriter, status int, obj interface{}) {
	res.WriteHeader(status)
	output := struct {
		Data interface{} `json:"data"`
	}{
		Data: obj,
	}
	jsonResp, err := json.Marshal(output)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}
