package sendinblue

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Sendinblue struct {
	StatusCode int    `json:"statusCode,omitempty"` //
	MessageId  string `json:"messageId,omitempty"`  //
	AuthKey    string `json:"authKey,omitempty"`    //
	Id         int    `json:"id,omitempty"`         //
	Message    string `json:"message,omitempty"`    //
	Code       string `json:"code,omitempty"`       //
}

type SibContact struct {
	Email         string `json:"email"`         //
	Firstname     string `json:"firstname"`     //
	Lastname      string `json:"lastname"`      //
	ListId        []int  `json:"listId"`        //
	UpdateEnabled bool   `json:"updateEnabled"` //

}

var reqUrl = "https://api.sendinblue.com/v3/smtp/email"
var apiUrl = "https://api.sendinblue.com/v3/"

func Sendmail(request []byte, apiKey string) (Sendinblue, error) {
	var t Sendinblue
	method := "POST"
	client := &http.Client{}

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Api-Key", apiKey)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	// check for response error
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode

	if resp.StatusCode == 201 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&t)
		if errjson != nil {
			return t, err
		}
	}
	// close response body
	resp.Body.Close()

	return t, nil
}

func Contacts(params SibContact, key string) (int, error) {

	return 201, nil
}

func SibRequest(request []byte, apiKey string, path string, method string) (Sendinblue, error) {
	var t Sendinblue

	client := &http.Client{}

	if path != "" {
		apiUrl = path
	}

	log.Println(apiUrl)
	log.Println(string(request))
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(request))
	req.Header.Add("Api-Key", apiKey)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	// check for response error
	if err != nil {
		return t, err
	}

	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode

	log.Println(resp.StatusCode)
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&t)
		if errjson != nil {
			return t, errjson
		}
	} else if resp.StatusCode == 204 {
		t.Message = "Request successful"
		return t, nil
	} else {
		data := json.NewDecoder(resp.Body)
		data.Decode(&t)
		return t, errors.New(t.Message)
	}

	return t, nil
}
