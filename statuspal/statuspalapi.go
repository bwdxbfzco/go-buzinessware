package statuspal

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

const apiUrl = "https://statuspal.io/api/v1/"
const subDomain = "buzinessware-com"

type StatuspalAPI struct {
	ApiKey string `json:"apiKey"`
}

type Subscription struct {
	ServiceIds   []int  `json:"service_ids"`             //
	SMSEnabled   bool   `json:"sms_enabled,omitempty"`   //
	Type         string `json:"type"`                    //
	Email        string `json:"email"`                   //
	PhoneNumber  string `json:"phone_number,omitempty"`  //
	EmailEnabled bool   `json:"email_enabled,omitempty"` //
	Id           string `json:"id,omitempty"`            //
	Confirm      bool   `json:"confirm"`                 //
	Filter       string `json:"filter,omitempty"`        //
}

type Subscriptions struct {
	Subscriptions []Subscription         `json:"subscriptions"`  //
	Meta          map[string]interface{} `json:"meta,omitempty"` //
}

type Children struct {
	Id       float64 `json:"id"`                  //
	Name     string  `json:"name"`                //
	ParentId string  `json:"parent_id,omitempty"` //
}

type Services struct {
	Services []Children `json:"services,omitempty"`
}

type StatuspalResponse struct {
	StatusCode    int           `json:"status_code"`             //
	Message       string        `json:"message,omitempty"`       //
	Subscriptions Subscriptions `json:"Subscriptions,omitempty"` //
	Services      Services      `json:"services,omitempty"`      //
	Id            string        `json:"id,omitempty"`            //
}

type IncidentParams struct {
	Title              string                  `json:"title"`               //
	StartsAt           string                  `json:"starts_at"`           //
	Type               string                  `json:"type"`                //
	ServiceIds         []int                   `json:"service_ids"`         //
	IncidentActivities []IncidentActivityParam `json:"incident_activities"` //
}

type IncidentActivityParam struct {
	ActivityTypeId int    `json:"activity_type_id"` //
	Description    string `json:"description"`      //
	Notify         bool   `json:"notify"`           //
}

func (s StatuspalAPI) StatusPage() (StatuspalResponse, error) {
	var request []byte
	result, err := s.apiCall(request, "status", "GET")
	return result, err
}

func (s StatuspalAPI) Subscriptions() (StatuspalResponse, error) {
	var _request []byte
	result, err := s.apiCall(_request, "subscriptions", "GET")
	return result, err
}

func (s StatuspalAPI) AddSubscription(params Subscription) (StatuspalResponse, error) {
	type addSub struct {
		Subscription Subscription `json:"subscription"` //
	}
	var r addSub
	r.Subscription = params
	_request, err := json.Marshal(r)
	if err != nil {
		return StatuspalResponse{}, err
	}
	result, err := s.apiCall(_request, "subscriptions", "POST")
	return result, err
}

func (s StatuspalAPI) UpdateSubscription(params Subscription, subid string) (StatuspalResponse, error) {
	type addSub struct {
		Subscription Subscription `json:"subscription"` //
	}
	var r addSub
	r.Subscription = params
	_request, err := json.Marshal(r)
	result, err := s.apiCall(_request, "subscriptions/"+subid, "PUT")
	return result, err
}

func (s StatuspalAPI) AddIncident(params IncidentParams) (StatuspalResponse, error) {
	type incident struct {
		Incident IncidentParams `json:"incident"` //
	}
	spew.Dump(params)
	var a incident
	a.Incident = params
	_request, err := json.Marshal(a)

	if err != nil {
		return StatuspalResponse{}, err
	}
	result, err := s.apiCall(_request, "incidents", "POST")
	return result, err
}

func (s StatuspalAPI) DeleteSubscription(subid string) (StatuspalResponse, error) {
	var _request []byte
	result, err := s.apiCall(_request, "subscriptions/"+subid, "DELETE")
	return result, err
}

func (s StatuspalAPI) apiCall(request []byte, path string, method string) (StatuspalResponse, error) {
	var _result StatuspalResponse
	client := &http.Client{}

	reqUrl := apiUrl + "status_pages/" + subDomain + "/" + path

	//log.Println(reqUrl)
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Authorization", s.ApiKey)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	// check for response error
	if err != nil {
		return _result, err
	}
	defer resp.Body.Close()

	/*	var i interface{}
		b := json.NewDecoder(resp.Body)
		b.Decode(&i)
		spew.Dump(i)*/

	if resp.StatusCode == 200 {
		if path == "subscriptions" {
			if method == "GET" {
				var result Subscriptions
				data := json.NewDecoder(resp.Body)
				errjson := data.Decode(&result)
				if errjson != nil {
					_result.StatusCode = 0
					_result.Message = errjson.Error()
				}
				_result.Subscriptions = result
				_result.StatusCode = resp.StatusCode
			} else {
				var i interface{}
				b := json.NewDecoder(resp.Body)
				b.Decode(&i)
				_result.Id = i.(map[string]interface{})["subscription"].(map[string]interface{})["id"].(string)
			}
		} else if path == "status" {
			var statusResult Services
			data := json.NewDecoder(resp.Body)
			errjson := data.Decode(&statusResult)
			if errjson != nil {
				_result.StatusCode = 0
				_result.Message = errjson.Error()
			}
			_result.Services = statusResult
			_result.StatusCode = resp.StatusCode
		} else {
			_result.StatusCode = resp.StatusCode
		}
	} else if resp.StatusCode == 204 {
		_result.StatusCode = resp.StatusCode
		_result.Message = "Subscription deleted successfully"
	}

	return _result, nil
}
