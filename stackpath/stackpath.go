package stackpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const liveUrl = "https://gateway.stackpath.com"

type Stackpath struct {
	StackId  string `json:"stack_id"`  //
	SPToken  string `json:"sp_token"`  //
	ClientId string `json:"client_id"` //
	Secret   string `json:"secret"`    //
}

type StackpathResponse struct {
	Code    int    `json:"code,omitempty"`    //
	Message string `json:"message,omitempty"` //
	Site    []struct {
		Id       string      `json:"id,omitempty"`       //
		StackId  string      `json:"stackId,omitempty"`  //
		Label    string      `json:"label,omitempty"`    //
		Status   string      `json:"status,omitempty"`   //
		Features interface{} `json:"features,omitempty"` //
	} `json:"site,omitempty"`
}

func (a Stackpath) auth() (string, error) {
	type authParam struct {
		GrantType    string `json:"grant_type"`    //
		ClientId     string `json:"client_id"`     //
		ClientSecret string `json:"client_secret"` //
	}

	var _params authParam

	_params.GrantType = "client_credentials"
	_params.ClientId = a.ClientId
	_params.ClientSecret = a.Secret

	request, _ := json.Marshal(_params)

	var _result interface{}
	client := &http.Client{}

	reqUrl := liveUrl + "/identity/v1/oauth2/token"

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&_result)
		if errjson != nil {
			log.Println(errjson)
		}
		return _result.(map[string]interface{})["access_token"].(string), nil
	} else {
		return "error", errors.New("Error accessing token.")
	}

	return "error", nil
}

func (a Stackpath) Request(request []byte, path string, method string) (int, interface{}, error) {
	var _result interface{}

	client := &http.Client{}
	reqUrl := liveUrl + path

	//Get token
	spToken, errToken := a.auth()

	if errToken != nil {
		return 400, _result, errToken
	}

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Authorization", "Bearer "+spToken)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return 400, _result, err
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	data.Decode(&_result)

	return resp.StatusCode, _result, nil
}
