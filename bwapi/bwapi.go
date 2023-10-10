package bwapi

import (
	"bytes"
	"errors"
	"net/http"

	validator "github.com/go-playground/validator/v10"
)

var reqUrl = "http://bwapi.buzinessware.com:9001"

type BWApi struct {
	Url string `json:"url"`
}

func (a BWApi) PostRequest(request []byte, path string, method string, username string, password string) (*http.Response, error) {
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

	err = validate.Var(username, "required")
	if err != nil {
		return resp, errors.New("Username not provided")
	}

	err = validate.Var(password, "required")
	if err != nil {
		return resp, errors.New("Password not provided")
	}

	client := &http.Client{}

	reqUrl = a.Url + path

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err = client.Do(req)

	// check for response error
	if err != nil {
		return resp, err
	}

	return resp, nil
}
