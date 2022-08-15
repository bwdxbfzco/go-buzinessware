package twiliosms

import (
	"fmt"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMS struct {
	Sid        string `json:"sid"`        //
	Token      string `json:"token"`      //
	FromNumber string `json:"fromNumber"` //
	ToNumber   string `json:"toNumber"`   //
	Message    string `json:"message"`    //
}

func (f TwilioSMS) SendSMS() interface{} {
	var s interface{}
	accountSid := f.Sid
	authToken := f.Token

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(f.ToNumber)
	params.SetFrom(f.FromNumber)
	params.SetBody(f.Message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	} else {
		fmt.Println("Message Sid: " + *resp.Sid)
	}
	return s
}
