package mobishastra

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type MobiShastra struct {
	Number  string `json:"number"`  //
	Message string `json:"message"` //
}

var User = "20079597"
var Password = "Dubai@123"
var Sender = "AD-BW TEAM"

func (c MobiShastra) SendSMS() interface{} {
	var s interface{}
	var request url.Values
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	client := http.Client{
		Timeout:   100 * time.Second,
		Transport: tr,
	}
	params := url.Values{}
	params.Set("user", User)
	params.Set("pwd", Password)
	params.Set("sendid", Sender)
	params.Set("priority", "High")
	params.Set("CountryCode", "ALL")
	params.Set("ShowError", "C")
	params.Set("mobileno", c.Number)
	params.Set("msgtext", c.Message)

	reqUrl := "https://mshastra.com/sendurl.aspx?" + params.Encode()

	resp, err := client.PostForm(reqUrl, request)
	defer resp.Body.Close()

	if err != nil {
		var r interface{} = err.Error()
		return r
	}
	data := json.NewDecoder(resp.Body)
	data.Decode(&s)
	return s
}
