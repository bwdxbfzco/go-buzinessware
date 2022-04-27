package businessemail

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

var reqUrl = "https://popbox.apps.ae:4443/api/"

type BusinessEmail struct {
	Url string `json:"url"`
}

func (c BusinessEmail) Request(request url.Values, action string) (map[string]interface{}, error) {
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
