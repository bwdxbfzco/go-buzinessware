package thesslstore

import (
	"bytes"
	"encoding/json"
	conv "github.com/cstockton/go-conv"
	colly "github.com/gocolly/colly"
	"log"
	"net/http"
	"strings"
)

type Thesslstore struct {
	PartnerCode string `json:"partner_code"` //
	Token       string `json:"token"`        //
}

type SSLRequest struct {
	AuthRequest struct {
		PartnerCode string `json:"PartnerCode"` //
		AuthToken   string `json:"AuthToken"`   //
	} `json:"AuthRequest"`                                  //
	ProductCode    string `json:"ProductCode,omitempty"`    //
	ProductType    int    `json:"ProductType,omitempty"`    //
	NeedSortedList bool   `json:"NeedSortedList,omitempty"` //
	IsForceNewSKUs bool   `json:"IsForceNewSKUs,omitempty"` //
}

type SSLProductResponse struct {
	IssuanceTime string `json:"IssuanceTime"`
	PricingInfo  []struct {
		Price float64 `json:"Price"`
		Srp   float64 `json:"SRP"`
	} `json:"PricingInfo"`
	ProductCode        string `json:"ProductCode"`
	ProductDescription string `json:"ProductDescription"`
	ProductName        string `json:"ProductName"`
	ProductSlug        string `json:"ProductSlug"`
	ProductType        int    `json:"ProductType"`
	VendorName         string `json:"VendorName"`
	IsDVProduct        bool   `json:"isDVProduct"`
	IsEVProduct        bool   `json:"isEVProduct"`
	IsOVProduct        bool   `json:"isOVProduct"`
	IsWildcard         bool   `json:"isWildcard"`
}

type ScrapperResponse struct {
	Title string  `json:"title"` //
	Price float64 `json:"price"` //
}

var LIVEURL = "https://api.thesslstore.com/rest"
var DEMOURL = "https://sandbox-wbapi.thesslstore.com/rest"

func (a Thesslstore) Request(request map[string]interface{}, path string, method string, test bool) (int, []SSLProductResponse, error) {
	var postRequest SSLRequest
	var _result []SSLProductResponse

	client := &http.Client{}
	reqUrl := LIVEURL + path

	if test == true {
		reqUrl = DEMOURL + path
	}

	postRequest.AuthRequest.AuthToken = a.Token
	postRequest.AuthRequest.PartnerCode = a.PartnerCode

	_request, err := json.Marshal(postRequest)
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(_request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return 400, _result, err
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	data.Decode(&_result)

	return resp.StatusCode, _result, nil
}

func (a Thesslstore) Scrapper() ([]ScrapperResponse, error) {
	var result []ScrapperResponse

	c := colly.NewCollector()
	c.OnHTML(".ssltbllist", func(e *colly.HTMLElement) {
		e.ForEach(".tblraw", func(_ int, el *colly.HTMLElement) {
			var result1 ScrapperResponse
			log.Printf("%v - %v", el.ChildText(".rawone"), el.ChildText(".rawtwo"))
			result1.Title = el.ChildText(".rawone")
			a := strings.Replace(el.ChildText(".rawtwo"), "$", "", 1)
			a = strings.Replace(a, "/yr.", "", 1)
			result1.Price, _ = conv.Float64(a)
			result = append(result, result1)
		})
	})
	c.Visit("https://www.thesslstore.com/brands.aspx")

	return result, nil
}
