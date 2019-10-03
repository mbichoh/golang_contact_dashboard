package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mbichoh/contactDash/pkg/models"
)

func main() {
	// endpoint
	var sendMessageURL string = "https://api.amisend.com/v1/sms/send"

	
	// data

	messageData := map[string]string{
		"phoneNumbers": ,
		"message":      ,
		"senderId":     "", // leave blank if you do not have a custom sender Id
	}

	params, _ := json.Marshal(messageData)

	request, err := http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(params))

	if err != nil {
		panic(err.Error())
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("x-api-user", models.Username)
	request.Header.Set("x-api-key", models.ApiKey)
	request.Header.Set("Content-Length", strconv.Itoa(len(params)))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	defer response.Body.Close()

	fmt.Println(string(body))
}
