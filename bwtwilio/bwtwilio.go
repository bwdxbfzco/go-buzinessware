package bwtwilio

import (
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var LOCALE = "en"

type BWTwilio struct {
	Sid          string `json:"sid"`          //
	Token        string `json:"token"`        //
	To           string `json:"to"`           //
	Channel      string `json:"channel"`      //
	ServiceId    string `json:"serviceId"`    //
	VerifySid    string `json:"verifySid"`    //
	VerifyStatus string `json:"verifyStatus"` //
	VerifyCode   string `json:"verifyCode"`   //
}

func (t BWTwilio) CreateVerification() (interface{}, error) {
	var s interface{}
	accountSid := t.Sid
	authToken := t.Token

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateVerificationParams{}
	params.Channel = &t.Channel
	params.To = &t.To
	params.Locale = &LOCALE
	resp, err := client.VerifyV2.CreateVerification(t.ServiceId, params)
	if err != nil {
		return s, err
	} else {
		s = *resp.Sid
	}
	return s, nil
}

func (t BWTwilio) FetchVerification() (interface{}, error) {
	var s interface{}
	accountSid := t.Sid
	authToken := t.Token

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	resp, err := client.VerifyV2.FetchVerification(t.ServiceId, t.VerifySid)
	if err != nil {
		return s, err
	} else {
		s = *resp
	}

	return s, nil
}

func (t BWTwilio) VerificationCheck() (interface{}, error) {
	var s interface{}
	accountSid := t.Sid
	authToken := t.Token

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateVerificationCheckParams{}
	params.Code = &t.VerifyCode
	params.To = &t.To
	resp, err := client.VerifyV2.CreateVerificationCheck(t.ServiceId, params)
	if err != nil {
		return s, err
	} else {
		s = *resp.Status
	}

	return s, nil
}

func (t BWTwilio) VerificationUpdate() (interface{}, error) {
	var s interface{}
	accountSid := t.Sid
	authToken := t.Token

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.UpdateVerificationParams{}
	resp, err := client.VerifyV2.UpdateVerification(t.ServiceId, t.VerifySid, params)
	if err != nil {
		return s, err
	} else {
		s = *resp
	}
	return s, nil
}
