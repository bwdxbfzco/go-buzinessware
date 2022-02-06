package odoo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	conv "github.com/cstockton/go-conv"
	validator "github.com/go-playground/validator/v10"
)

const odooUrl = "https://buzinessware-staging-3962576.dev.odoo.com/api"
const odooAPI = "KFVY8S8Y5XV575GD7UP4XU85EINFEKEM"

type InvoiceItems struct {
	ProductID        int     `json:"product_id" validate:"required"`         //
	Name             string  `json:"name" validate:"required"`               //
	PriceUnit        float64 `json:"price_unit" validate:"required"`         //
	PriceSubtotal    float64 `json:"price_subtotal" validate:"required"`     //
	PriceTotal       float64 `json:"price_total" validate:"required"`        //
	TermName         string  `json:"term_name" validate:"required"`          //
	DateStartingg    string  `json:"date_startingg" validate:"required"`     //
	DateEndingg      string  `json:"date_endingg"`                           //
	TaxIds           []int   `json:"tax_ids"`                                //
	OrderTypee       string  `json:"order_typee" validate:"required"`        //
	ProductGroupName int     `json:"product_group_name" validate:"required"` //
	BusinessUnit     string  `json:"business_unit" validate:"required"`      //
}

type Invoice struct {
	MoveType                  string         `json:"move_type"`                                        //
	Name                      string         `json:"name" validate:"required"`                         //
	InvoiceDate               string         `json:"invoice_date" validate:"required"`                 //
	InvoiceDateDue            string         `json:"invoice_date_due"`                                 //
	InvoiceLineIds            []InvoiceItems `json:"invoice_line_ids" validate:"required"`             //
	PaymentReference          string         `json:"payment_reference" validate:"required"`            //
	ClientID                  int            `json:"client_id" validate:"required"`                    //
	PartnerID                 int            `json:"partner_id" validate:"required"`                   //
	InvoicePartnerDisplayName string         `json:"invoice_partner_display_name" validate:"required"` //
	CurrencyID                int            `json:"currency_id" validate:"required"`                  //
	JournalID                 int            `json:"journal_id,omitempty"`                             //
	AmountUntaxed             float64        `json:"amount_untaxed,omitempty"`                         //
	AmountTax                 float64        `json:"amount_tax,omitempty"`                             //
	AmountTotal               float64        `json:"amount_total,omitempty"`                           //
	AmountResidual            float64        `json:"amount_residual,omitempty"`                        //
	AmountUntaxedSigned       float64        `json:"amount_untaxed_signed,omitempty"`                  //
	AmountTaxSigned           float64        `json:"amount_tax_signed,omitempty"`                      //
	AmountTotalSigned         float64        `json:"amount_total_signed,omitempty"`                    //
	AmountResidualSigned      float64        `json:"amount_residual_signed,omitempty"`                 //
	ExtractState              string         `json:"extract_state,omitempty"`                          //
	State                     string         `json:"state,omitempty"`                                  //
}

type ModelOdoo struct {
	Result string `json.result`
}

type Odoocreatecontact struct {
	Name         string `json:"name,omitempty"`
	Emailaddress string `json:"emailaddress,omitempty"`
}

type Odooresponse struct {
	Success      bool   `json:"success,omitempty"`
	Message      string `json:"message,omitempty"`
	ResponseCode int    `json:"responseCode,omitempty"`
	CreateID     int    `json:"create_id,omitempty"`
	Data         []struct {
		ID   int    `json:"id,omitempty"`   //
		Name string `json:"name,omitempty"` //
	} `json:"data,omitempty"`
}

type Contact struct {
	Name            string `json:"name" validate:"required"`                      //
	DisplayName     string `json:"display_name"`                                  //
	ClientId        int    `json:"client_id" validate:"required,numeric"`         //
	Phone           string `json:"phone"`                                         //
	EmailNormalized string `json:"email_normalized"`                              //
	Email           string `json:"email" validate:"required"`                     //
	CustomerType    string `json:"x_studio_customer_segment" validate:"required"` //
}

func Search(object string, field string, search string, filter string) (Odooresponse, error) {
	values := []byte(`{}`)
	params := url.Values{}
	var a string

	if search != "" && field != "" {
		if filter == "" {
			filter = "ilike"
		}
		params.Add("domain", `[["`+field+`", "`+filter+`", "`+search+`"]]`)
		a = params.Encode()
	}

	reqUrl := odooUrl + "/" + object + "/search?" + a

	log.Println(reqUrl)
	method := "GET"

	return odooApiCall(values, reqUrl, method)
}

func odooApiCall(request []byte, reqUrl string, method string) (Odooresponse, error) {
	var t Odooresponse
	client := &http.Client{}

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	if method == "GET" {
		req, err = http.NewRequest(method, reqUrl, nil)
		if err != nil {
			return t, err
		}
	}
	req.Header.Add("api-key", odooAPI)

	if method == "POST" {
		req.Header.Add("Content-type", "text/plain")
	}

	resp, err := client.Do(req)

	// check for response error
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&t)
		if errjson != nil {
			return t, errjson
		}
	}

	// close response body
	resp.Body.Close()

	return t, nil
}

func CreateCustomer(_customerRequest Contact) (Odooresponse, error) {
	var t Odooresponse
	validate := validator.New()
	err := validate.Struct(&_customerRequest)

	if err != nil {
		return t, err
	}

	values, _ := json.Marshal(_customerRequest)
	reqUrl := odooUrl + "/res.partner/create"

	method := "POST"
	t, err = odooApiCall(values, reqUrl, method)

	return t, err
}

func CreateInvoice(_invoiceRequest Invoice) (Odooresponse, error) {
	if _invoiceRequest.MoveType == "" {
		_invoiceRequest.MoveType = "out_invoice"
	}
	if _invoiceRequest.JournalID == 0 {
		_invoiceRequest.JournalID = 2
	}
	if _invoiceRequest.ExtractState == "" {
		_invoiceRequest.ExtractState = "no_extract_requested"
	}

	var t Odooresponse
	validate := validator.New()
	err := validate.Struct(&_invoiceRequest)

	if err != nil {
		return t, err
	}

	values, _ := json.Marshal(_invoiceRequest)
	reqUrl := odooUrl + "/account.move/create"

	method := "POST"
	t, err = odooApiCall(values, reqUrl, method)

	/*	//update the created invoice from draft to posted
		valuesPut := []byte(`{
			"state": "posted"
		}`)
		//invoiceid, _ := conv.String(t.CreateID)
		//reqUrl = odooUrl + "/account.move/" + invoiceid
		reqUrl = odooUrl + "/account.move/create"

		log.Println(reqUrl)
		method = "POST"
		t, err = odooApiCall(valuesPut, reqUrl, method)*/

	return t, err
}

func DeleteRecord(id int, object string) (Odooresponse, error) {
	_id, _ := conv.String(id)
	values := []byte(`{}`)
	reqUrl := odooUrl + "/" + object + "/" + _id

	method := "DELETE"
	t, err := odooApiCall(values, reqUrl, method)

	return t, err
}

func UpdateRecord(id int, values []byte, object string) (Odooresponse, error) {
	_id, _ := conv.String(id)
	reqUrl := odooUrl + "/" + object + "/" + _id

	method := "PUT"
	t, err := odooApiCall(values, reqUrl, method)
	return t, err
}

func CreateRecord(values []byte, object string) (Odooresponse, error) {
	method := "POST"
	reqUrl := odooUrl + "/" + object + "/create"

	t, err := odooApiCall(values, reqUrl, method)

	return t, err
}
