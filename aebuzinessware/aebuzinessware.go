package aebuzinessware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const liveUrl = "http://ae-epp.apps.ae/index.php?"
const testUrl = "http://ae-epp-test.apps.ae/index.php?"

type AEDomain struct {
	Testmode string `json:"testmode"` //
	Username string `json:"username"` //
	Password string `json:"password"` //
}

type Response struct {
	Status          string          `json:"status"`                    //
	DomainResponse  DomainResponse  `json:"domainResponse,omitempty"`  //
	ContactResponse ContactResponse `json:"contactResponse,omitempty"` //
	ContactId       string          `json:"contactId,omitempty"`       //
}

type DomainResponse struct {
	Status         string        `json:"status"`          //
	ExpirationDate time.Time     `json:"expirationDate"`  //
	AuthInfo       string        `json:"auth_info"`       //
	Domainname     string        `json:"domainname"`      //
	CreatedOn      time.Time     `json:"createdOn"`       //
	LastUpdateOn   time.Time     `json:"lastUpdateOn"`    //
	Registrant     string        `json:"registrant"`      //
	ContactInfo    []interface{} `json:"contact_info"`    //
	NameserverInfo []interface{} `json:"nameserver_info"` //
}

type ContactResponse struct {
	ContactID          string    `json:"ContactId"`
	ContactEmail       string    `json:"ContactEmail"`
	ContactClientID    string    `json:"ContactClientId"`
	ContactCreateDate  time.Time `json:"ContactCreateDate"`
	ContactUpdateDate  time.Time `json:"ContactUpdateDate"`
	ContactStatus      []string  `json:"ContactStatus"`
	ContactVoice       string    `json:"ContactVoice"`
	ContactFax         string    `json:"ContactFax"`
	ContactName        string    `json:"ContactName"`
	ContactStreet      string    `json:"ContactStreet"`
	ContactCity        string    `json:"ContactCity"`
	ContactZipcode     string    `json:"ContactZipcode"`
	ContactProvince    string    `json:"ContactProvince"`
	ContactCountrycode string    `json:"ContactCountrycode"`
	ContactCompanyname string    `json:"ContactCompanyname"`
}

type Nameserver struct {
	Nameserver string
}

type Contact struct {
	FirstName string `json:"First Name"`        //
	LastName  string `json:"Last Name"`         //
	Company   string `json:"Organisation Name"` //
	Address1  string `json:"Address 1"`         //
	Address2  string `json:"Address 2"`         //
	City      string `json:"City"`              //
	State     string `json:"State"`             //
	Postcode  string `json:"Postcode"`          //
	Country   string `json:"Country"`           //
	Email     string `json:"Email"`             //
	Phone     string `json:"Phone"`             //

}

func (c AEDomain) DomainDetails(domainName string) (Response, error) {
	params := url.Values{}
	params.Add("u", c.Username)
	params.Add("p", c.Password)
	params.Add("action", "info")
	params.Add("domain", domainName)

	result, err := c.apiCall(params, "info")
	return result, err
}

func (c AEDomain) ContactDetails(contactId string) (Response, error) {
	params := url.Values{}
	params.Add("u", c.Username)
	params.Add("p", c.Password)
	params.Add("action", "contactinfo")
	params.Add("ContactId", contactId)

	result, err := c.apiCall(params, "contactinfo")
	return result, err
}

func (c AEDomain) CheckContactByEmail(email string) (Response, error) {
	params := url.Values{}
	params.Set("action", "CheckConatctByEmail")
	params.Set("email", email)
	result, err := c.apiCall(params, "checkcontactbyemail")

	return result, err
}

func (c AEDomain) CreateContact(contact Contact) (Response, error) {
	_contact, _ := json.Marshal(contact)
	params := url.Values{}
	params.Set("action", "CreateContact")
	params.Set("data", string(_contact))
	result, err := c.apiCall(params, "createcontact")
	return result, err
}

func (c AEDomain) CreateChildNameServers() {

}

func (c AEDomain) UpdateChildNameServers() {

}

func (c AEDomain) DeleteChildNameServers() {

}

func (c AEDomain) UpdateNameServers(nameservers map[string]string, domainName string) (Response, error) {
	//var result Response
	params := url.Values{}
	var _oldns []string
	var _newns []string
	var oldn string
	var newn string
	//Get Domain Info
	_domainInfo, err := c.DomainDetails(domainName)
	if err != nil {
		return _domainInfo, err
	}

	for _, x := range _domainInfo.DomainResponse.NameserverInfo {
		oldn = strings.TrimSpace(x.(interface{}).(string))
		_oldns = append(_oldns, oldn)
	}
	_oldnsJSON, _ := json.Marshal(_oldns)

	if nameservers["nameserver1"] != "" {
		newn = nameservers["nameserver1"]
		_newns = append(_newns, newn)
	}
	if nameservers["nameserver2"] != "" {
		newn = nameservers["nameserver2"]
		_newns = append(_newns, newn)
	}
	if nameservers["nameserver3"] != "" {
		newn = nameservers["nameserver3"]
		_newns = append(_newns, newn)
	}
	if nameservers["nameserver4"] != "" {
		newn = nameservers["nameserver4"]
		_newns = append(_newns, newn)
	}

	_newnsJSON, _ := json.Marshal(_newns)

	params.Set("action", "savenameservers")
	params.Set("remnameservers", string(_oldnsJSON))
	params.Set("nameserver", string(_newnsJSON))
	params.Set("domain", domainName)

	result, err := c.apiCall(params, "savenameservers")
	return result, err
}

func (c AEDomain) UpdateDomain(contactId string, domainName string, contactType string) (Response, error) {
	contacts := make(map[string]string)
	params := url.Values{}
	if strings.ToLower(contactType) == "registrant" {
		contacts["type"] = "CONTACT_TYPE_REGISTRANT"
	}
	if strings.ToLower(contactType) == "admin" {
		contacts["type"] = "CONTACT_TYPE_ADMIN"
	}
	if strings.ToLower(contactType) == "technical" || strings.ToLower(contactType) == "tech" {
		contacts["type"] = "CONTACT_TYPE_TECH"
	}
	if strings.ToLower(contactType) == "billing" {
		contacts["type"] = "CONTACT_TYPE_BILLING"
	}
	//Get Domain Info
	_domInfo, err := c.DomainDetails(domainName)
	if err != nil {
		return Response{}, err
	}

	for _, x := range _domInfo.DomainResponse.ContactInfo {
		_split := strings.Split(strings.TrimSpace(x.(interface{}).(string)), ":")
		if _split[0] == "tech" && (contactType == "technical" || contactType == "tech") {
			contacts["remContactId"] = strings.TrimSpace(_split[1])
		} else if _split[0] == "admin" && contactType == "admin" {
			contacts["remContactId"] = strings.TrimSpace(_split[1])
		} else if _split[0] == "billing" && contactType == "billing" {
			contacts["remContactId"] = strings.TrimSpace(_split[1])
		}
	}

	if contactType == "registrant" {
		contacts["remContactId"] = _domInfo.DomainResponse.Registrant
	}

	contacts["addContactId"] = contactId

	_contact, _ := json.Marshal(contacts)

	params.Set("action", "domaincontactupdate")
	params.Set("domain", domainName)
	params.Set("data", string(_contact))

	log.Printf("%s", params)
	result, err := c.apiCall(params, "updatedomain")
	return result, err
}

func (c AEDomain) apiCall(request url.Values, action string) (Response, error) {
	var _response Response
	var _contact ContactResponse
	var _domain DomainResponse

	client := http.Client{
		Timeout: 100 * time.Second,
	}
	reqUrl := liveUrl
	if c.Testmode == "on" {
		reqUrl = testUrl
	}
	u, _ := url.ParseRequestURI(reqUrl)
	urlStr := u.String()

	resp, err := client.PostForm(urlStr, request)

	// check for response error
	if err != nil {
		return _response, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		if action == "contactinfo" {
			errjson := data.Decode(&_contact)
			if errjson != nil {
				return _response, errjson
			}
			if _contact.ContactID != "" {
				_response.Status = "success"
				_response.ContactResponse = _contact
			} else {
				_response.Status = "error"
				return _response, errors.New("no such contact")
			}
		} else if action == "info" {
			errjson := data.Decode(&_domain)
			if errjson != nil {
				return _response, errjson
			}
			if _domain.Status != "" {
				_response.Status = "success"
				_response.DomainResponse = _domain
			} else {
				_response.Status = "error"
				return _response, errors.New("no such domain")
			}
		} else if action == "savenameservers" {
			var r interface{}
			errjson := data.Decode(&r)
			if errjson != nil {
				return _response, errjson
			}
			_response.Status = "success"
			log.Printf("%s", r)
		} else if action == "checkcontactbyemail" {
			var r interface{}
			errjson := data.Decode(&r)
			if errjson != nil {
				return _response, errjson
			}
			if len(r.(interface{}).(string)) > 0 {
				_response.Status = "success"
				_response.ContactId = r.(interface{}).(string)
			} else {
				_response.Status = "error"
				return _response, errors.New("No such contact.")
			}
		} else if action == "createcontact" {
			var r interface{}
			errjson := data.Decode(&r)
			if errjson != nil {
				return _response, errjson
			}
			if len(r.(interface{}).(map[string]interface{})["ContactId"].(string)) > 0 {
				_response.Status = "success"
				_response.ContactId = r.(interface{}).(map[string]interface{})["ContactId"].(string)
			} else {
				_response.Status = "error"
				return _response, errors.New("error create contact")
			}
		} else if action == "updatedomain" {
			_response.Status = "success"
		}
	}
	return _response, nil
}
