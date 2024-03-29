package resellerclub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type Contact struct {
	Emailaddr  string `json:"emailaddr"`            //
	Country    string `json:"country"`              //
	Name       string `json:"name"`                 //
	Lastname   string `json:"lastname,omitempty"`   //
	Address2   string `json:"address2,omitempty"`   //
	Company    string `json:"company,omitempty"`    //
	City       string `json:"city,omitempty"`       //
	Address1   string `json:"address1,omitempty"`   //
	Zip        string `json:"zip,omitempty"`        //
	Telno      string `json:"telno,omitempty"`      //
	State      string `json:"state,omitempty"`      //
	Type       string `json:"type"`                 //
	Contactid  string `json:"contactid,omitempty"`  //
	Telnocc    string `json:"telnocc,omitempty"`    //-
	Customerid string `json:"customerid,omitempty"` //
}

type ResellerClub struct {
	ResellerClucUrl      string  `json:"resellerClucUrl,omitempty"`
	ResellerClubPassword string  `json:"resellerClubPassword,omitempty"`
	ResellerClubUser     string  `json:"resellerClubUser,omitempty"`
	Domain               string  `json:"description,omitempty"`
	Orderid              string  `json:"orderid,omitempty"`
	Actionstatus         string  `json:"actionstatus,omitempty"`
	Description          string  `json:"description,omitempty"`
	Status               string  `json:"status,omitempty"`
	Message              string  `json:"message,omitempty"`
	Actionstatusdesc     string  `json:"actionstatusdesc,omitempty"`
	Domsecret            string  `json:"domsecret,omitempty"`
	Contactid            int     `json:"contactid,omitempty"`
	TechContact          Contact `json:"techcontact,omitempty"`
	AdminContact         Contact `json:"admincontact,omitempty"`
	Registrantcontact    Contact `json:"registrantcontact,omitempty"`
	Billingcontact       Contact `json:"billingcontact,omitempty"`
	Billingcontactid     string  `json:"billingcontactid,omitempty"`
	Admincontactid       string  `json:"admincontactid,omitempty"`
	Techcontactid        string  `json:"techcontactid,omitempty"`
	Registrantcontactid  string  `json:"registrantcontactid,omitempty"`
	Ns1                  string  `json:"ns1,omitempty"`
	Ns2                  string  `json:"ns2,omitempty"`
	Ns3                  string  `json:"ns3,omitempty"`
	Ns4                  string  `json:"ns4,omitempty"`
	Customerid           string  `json:"customerid,omitempty"`
	ExpiryDate           string  `json:"endtime,omitempty"`
}

type DomainRecords struct {
	TotalRecords string `json:"recsindb"`
	PageNumber   string `json:"recsonpage"`
}

type Domains struct {
	EntityId          string `json:"entityId"`          //
	OrderId           string `json:"OrderId"`           //
	CustomerId        string `json:"customerId"`        //
	ResellerLock      string `json:"resellerLock"`      //
	OrderDate         string `json:"orderDate"`         //
	CustomerLock      string `json:"customerLock"`      //
	DomainName        string `json:"domainName"`        //
	ExpiryDate        string `json:"expiryDate"`        //
	PrivacyProtection string `json:"privacyProtection"` //
	AutoRenew         string `json:"autoRenew"`         //
	Status            string `json:"status"`            //
	TransferLock      string `json:"transferLock"`      //
}

func (u ResellerClub) apiCall(actionUrl string, params url.Values) (*http.Response, error) {
	params.Add("auth-userid", u.ResellerClubUser)
	params.Add("api-key", u.ResellerClubPassword)
	a := params.Encode()

	url := u.ResellerClucUrl + actionUrl + "?" + a

	resp, err := http.Get(url)

	if err != nil {
		log.Printf("go-buzinessware apiCall(91): %v\n", err.Error())
	}

	return resp, err
}

func (u ResellerClub) ResellerClubApi(_params map[string]string) ResellerClub {
	var _resellerclubresponse ResellerClub
	var a string

	params := url.Values{}

	if _params["action"] == "getorder" {
		params.Add("domain-name", _params["domain-name"])
		a = params.Encode()
	}

	if _params["action"] == "modifyns" {
		var ns string
		m := make(map[string]string)
		json.Unmarshal([]byte(_params["nameserver"]), &m)
		for _, v := range m {
			ns = ns + "&ns=" + v
		}
		params.Add("order-id", _params["order-id"])
		a = params.Encode()
		a = a + ns
	}

	if _params["action"] == "createchildnameserver" {
		params.Add("order-id", _params["order-id"])
		params.Add("cns", _params["cns"])
		params.Add("ip", _params["ip"])
		a = params.Encode()
	}

	if _params["action"] == "modifychildnameserver" {
		params.Add("order-id", _params["order-id"])
		params.Add("cns", _params["cns"])
		params.Add("old-ip", _params["old-ip"])
		params.Add("new-ip", _params["new-ip"])
		a = params.Encode()
	}

	if _params["action"] == "deletechildnameserver" {
		params.Add("order-id", _params["order-id"])
		params.Add("cns", _params["cns"])
		params.Add("ip", _params["ip"])
		a = params.Encode()
	}

	if _params["action"] == "theftprotection" {
		params.Add("order-id", _params["order-id"])
		a = params.Encode()
	}

	if _params["action"] == "domaininfo" {
		params.Add("domain-name", _params["domain"])
		params.Add("options", _params["options"])
		a = params.Encode()
	}

	if _params["action"] == "modifydomaincontact" {
		params.Add("order-id", _params["order-id"])
		params.Add("designated-agent", "true")
		params.Add("admin-contact-id", _params["admin-contact-id"])
		params.Add("tech-contact-id", _params["tech-contact-id"])
		params.Add("billing-contact-id", _params["billing-contact-id"])
		params.Add("reg-contact-id", _params["reg-contact-id"])
		a = params.Encode()
	}

	if _params["action"] == "createcontact" {
		params.Add("name", _params["name"])
		params.Add("company", _params["company"])
		params.Add("email", _params["email"])
		params.Add("address-line-1", _params["address-line1"])
		params.Add("city", _params["city"])
		params.Add("country", _params["country"])
		params.Add("zipcode", _params["zipcode"])
		params.Add("phone", _params["phone"])
		params.Add("phone-cc", _params["phonecc"])
		params.Add("type", _params["type"])
		params.Add("address-line-2", _params["address-line-2"])
		params.Add("state", _params["state"])
		params.Add("customer-id", _params["customer-id"])

		a = params.Encode()
	}

	if _params["action"] == "suspendorder" {
		params.Add("order-id", _params["order-id"])
		params.Add("reason", _params["reason"])
		a = params.Encode()
	}

	if _params["action"] == "resumeorder" {
		params.Add("order-id", _params["order-id"])
		a = params.Encode()
	}

	if _params["action"] == "reneworder" {
		params.Add("order-id", _params["order-id"])
		params.Add("months", _params["months"])
		params.Add("invoice-option", "NoInvoice")
		a = params.Encode()
	}

	resp, err := u.apiCall(_params["actionUrl"], params)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("go-buzinessware ResellerClubApi(191): %v\n", err.Error())
	}

	var r map[string]interface{}
	json.Unmarshal(body, &r)

	if resp.StatusCode == 200 {
		if _params["action"] != "getorder" && _params["action"] != "createcontact" {
			err = json.Unmarshal(body, &_resellerclubresponse)
			if err != nil {
				log.Printf("go-buzinessware ResellerClubApi(200): %v\n", err.Error())
			}
		} else if _params["action"] == "getorder" {
			_resellerclubresponse.Orderid = fmt.Sprintf("%s", body)
		} else if _params["action"] == "createcontact" {
			_resellerclubresponse.Contactid, _ = strconv.Atoi(fmt.Sprintf("%s", body))
		}
	} else {
		err = json.Unmarshal(body, &_resellerclubresponse)
		if err != nil {
			log.Printf("go-buzinessware ResellerClubApi(211): %v\n", err.Error())
		}
	}

	return _resellerclubresponse
}

func (u ResellerClub) ResellerClubApiPost(actionUrl string, params url.Values) ResellerClub {
	var result ResellerClub
	params.Add("auth-userid", u.ResellerClubUser)
	params.Add("api-key", u.ResellerClubPassword)
	a := params.Encode()

	reqUrl := u.ResellerClucUrl + actionUrl + "?" + a

	method := "POST"
	var request []byte

	client := &http.Client{}
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("%+v", err)
		}
	}

	return result
}

func (u ResellerClub) DomainDetails(domain string, orderType string) ResellerClub {
	var result ResellerClub

	_params := make(map[string]string)
	_params["action"] = "domaininfo"
	_params["actionUrl"] = "domains/details-by-name.json"
	_params["domain"] = domain
	_params["options"] = orderType

	result = u.ResellerClubApi(_params)

	if result.Domsecret != "" {
		result.Status = "Success"
	}

	return result
}

func (u ResellerClub) GetAllDomains(noofrecords string, pageno string) []Domains {
	params := url.Values{}
	params.Add("no-of-records", noofrecords)
	params.Add("page-no", pageno)
	params.Add("show-child-orders", "true")
	params.Add("order-by", "creationtime desc")

	resp, err := u.apiCall("domains/search.json", params)
	if err != nil {
		log.Printf("go-buzinessware GetAllDomains(247): %v\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("go-buzinessware GetAllDomains(254): %v\n", err.Error())
	}

	var x2 []Domains

	var r map[string]interface{}
	json.Unmarshal(body, &r)

	for _, b := range r {
		var x1 Domains
		if reflect.TypeOf(b).Kind().String() == "map" && reflect.TypeOf(b).Kind().String() != "string" {
			if b.(map[string]interface{})["entity.entityid"] != nil {
				x1.EntityId = b.(map[string]interface{})["entity.entityid"].(string)
			}
			if b.(map[string]interface{})["orders.orderid"] != nil {
				x1.OrderId = b.(map[string]interface{})["orders.orderid"].(string)
			}
			if b.(map[string]interface{})["entity.customerid"] != nil {
				x1.CustomerId = b.(map[string]interface{})["entity.customerid"].(string)
			}
			if b.(map[string]interface{})["orders.resellerlock"] != nil {
				x1.ResellerLock = b.(map[string]interface{})["orders.resellerlock"].(string)
			}
			if b.(map[string]interface{})["orders.customerlock"] != nil {
				x1.CustomerLock = b.(map[string]interface{})["orders.customerlock"].(string)
			}
			if b.(map[string]interface{})["entity.description"] != nil {
				x1.DomainName = b.(map[string]interface{})["entity.description"].(string)
			}
			if b.(map[string]interface{})["orders.endtime"] != nil {
				x1.ExpiryDate = b.(map[string]interface{})["orders.endtime"].(string)
			}
			if b.(map[string]interface{})["orders.privacyprotection"] != nil {
				x1.PrivacyProtection = b.(map[string]interface{})["orders.privacyprotection"].(string)
			}
			if b.(map[string]interface{})["orders.autorenew"] != nil {
				x1.AutoRenew = b.(map[string]interface{})["orders.autorenew"].(string)
			}
			if b.(map[string]interface{})["entity.currentstatus"] != nil {
				x1.Status = b.(map[string]interface{})["entity.currentstatus"].(string)
			}
			if b.(map[string]interface{})["orders.transferlock"] != nil {
				x1.TransferLock = b.(map[string]interface{})["orders.transferlock"].(string)
			}
			if b.(map[string]interface{})["orders.creationdt"] != nil {
				x1.OrderDate = b.(map[string]interface{})["orders.creationdt"].(string)
			}
		}
		if x1.EntityId != "" {
			x2 = append(x2, x1)
		}
	}
	return x2
}

func (u ResellerClub) GetTotalDomainCount() DomainRecords {
	params := url.Values{}
	params.Add("no-of-records", "10")
	params.Add("page-no", "1")
	params.Add("show-child-orders", "true")
	params.Add("order-by", "creationtime desc")

	resp, err := u.apiCall("domains/search.json", params)
	if err != nil {
		log.Printf("go-buzinessware GetTotalDomainCount(317): %v\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("go-buzinessware GetTotalDomainCount(326): %v\n", err.Error())
	}

	var record DomainRecords
	json.Unmarshal(body, &record)

	return record
}

func (u ResellerClub) SuspendService(orderId int, reason string) (string, error) {
	_params := url.Values{}
	_params.Add("order-id", strconv.Itoa(orderId))
	_params.Add("reason", reason)

	result := u.ResellerClubApiPost("orders/suspend.json", _params)

	if result.Status == "ERROR" {
		return result.Status, errors.New(result.Message)
	}

	return result.Status, nil
}

func (u ResellerClub) ResumeService(orderId int) (string, error) {
	_params := url.Values{}
	_params.Add("order-id", strconv.Itoa(orderId))

	result := u.ResellerClubApiPost("orders/unsuspend.json", _params)

	if result.Status == "ERROR" {
		return result.Status, errors.New(result.Message)
	}

	return result.Status, nil
}

func (u ResellerClub) RenewService(orderId int, duration int, service string) (string, error) {
	var result ResellerClub

	if service == "gsuite" {

		_params := url.Values{}
		_params.Add("order-id", strconv.Itoa(orderId))
		_params.Add("month", strconv.Itoa(duration))

		result = u.ResellerClubApiPost("gapps/gbl/renew.json", _params)

		if result.Status == "ERROR" {
			return result.Status, errors.New(result.Message)
		}
		return result.Status, nil
	}
	return "", errors.New("no service details provided")
}
