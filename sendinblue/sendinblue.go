package sendinblue

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Sendinblue struct {
	StatusCode int    `json:"statusCode"` //
	MessageId  string `json:"messageId"`  //
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
