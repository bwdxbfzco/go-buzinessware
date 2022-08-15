package bps

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Bps struct {
	Username string `json:"Username"` //
	Password string `json:"password"` //
}

const Url = "https://portal.bpsme.com/api"

func (f Bps) PostRequest(request []byte, path string, method string) (int, []interface{}, error) {
	var _result []interface{}

	client := &http.Client{}
	reqUrl := Url + path

	token, err := f.auth()
	if err != nil {
		return 0, _result, err
	}

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("token", token)
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

func (a Bps) auth() (string, error) {
	var _result interface{}
	client := &http.Client{}

	reqUrl := Url + "/PublicApi/GetAuth?username=" + a.Username + "&password=" + a.Password

	req, err := http.NewRequest("GET", reqUrl, nil)
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
		return _result.(map[string]interface{})["EncryptedToken"].(string), nil
	} else {
		return "error", errors.New("Error accessing token.")
	}

	return "error", nil
}
