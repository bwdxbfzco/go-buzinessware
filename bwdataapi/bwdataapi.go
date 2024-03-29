package bwdataapi

import (
	"bytes"
	"errors"
	"net/http"

	validator "github.com/go-playground/validator/v10"
)

var reqUrl = "https://api.buzinessware.com"

type BWDataApi struct {
	Url string `json:"url"`
}

func (a BWDataApi) PostRequest(request []byte, path string, method string, username string, password string) (*http.Response, error) {
	var resp *http.Response

	//Validation
	validate := validator.New()
	err := validate.Var(path, "required")
	if err != nil {
		return resp, errors.New("Path not provided")
	}

	err = validate.Var(method, "required")
	if err != nil {
		return resp, errors.New("Method not provided")
	}

	client := &http.Client{}

	reqUrl = a.Url + path

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err = client.Do(req)

	// check for response error
	if err != nil {
		return resp, err
	}

	return resp, nil
}
