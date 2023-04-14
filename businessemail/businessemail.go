package businessemail

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	validator "github.com/go-playground/validator/v10"
)

var reqUrl = "http://popbox.apps.ae:9002/"

type BusinessEmail struct {
	Url string `json:"url"`
}

func (c BusinessEmail) RequestOld(request url.Values, action string) (map[string]interface{}, error) {
	_response := make(map[string]interface{})
	client := http.Client{
		Timeout: 100 * time.Second,
	}
	reqUrl := reqUrl + action
	u, _ := url.ParseRequestURI(reqUrl)
	urlStr := u.String()

	resp, err := client.PostForm(urlStr, request)

	// check for response error
	if err != nil {
		return _response, err
	}

	var _r interface{}
	defer resp.Body.Close()
	data := json.NewDecoder(resp.Body)
	data.Decode(&_r)

	_response["statuscode"] = resp.StatusCode
	_response["status"] = _r

	return _response, nil
}

func (c BusinessEmail) Request(request []byte, path string, method string, username string, password string) (interface{}, error) {
	var result interface{}
	var url string

	//Validation
	validate := validator.New()
	err := validate.Var(path, "required")
	if err != nil {
		return result, errors.New("Path not provided")
	}

	err = validate.Var(method, "required")
	if err != nil {
		return result, errors.New("Method not provided")
	}

	client := &http.Client{Timeout: 100 * time.Second}

	url = reqUrl + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	// check for response error
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		data.Decode(&result)
		return result, nil
	} else if resp.StatusCode == 404 {
		data := json.NewDecoder(resp.Body)
		data.Decode(&result)
		_err := "error"
		if result.(map[string]interface{}) != nil {
			_err = result.(map[string]interface{})["description"].(string)
		}
		return result, errors.New(_err)
	} else if resp.StatusCode == 400 {
		data := json.NewDecoder(resp.Body)
		data.Decode(&result)
	} else {
		var _errResult string
		data := json.NewDecoder(resp.Body)
		data.Decode(&_errResult)
		return result, errors.New(_errResult)
	}

	return result, nil
}
