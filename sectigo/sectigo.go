package sectigo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const liveUrl = "https://secure.sectigo.com/products/"

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
var bitSize = 2048

type SectigoSSL struct {
	Testmode string `json:"testmode"` //
	Username string `json:"username"` //
	Password string `json:"password"` //
}

type SectigoParams struct {
	Csr               string `json:"csr,omitempty"`                 //
	ProductId         string `json:"product,omitempty"`             //
	Days              int    `json:"days,omitempty"`                //
	Domains           string `json:"domainNames,omitempty"`         //
	PrimaryDomain     string `json:"primaryDomainName,omitempty"`   //
	Validation        string `json:"validationTokens,omitempty"`    // HTTPCSRHASH, CNAMECSRHASH, IPADDRESSPRE
	Dcvmethod         string `json:"dcvmethod,omitempty"`           // Email, HTTP_CSR_HASH, CNAME_CSR_HASH, IP_ADDRESS_PRE
	Test              string `json:"test,omitempty"`                //
	Login             string `json:"loginName,omitempty"`           //
	Password          string `json:"loginPassword,omitempty"`       //
	CustomerValidated string `json:"isCustomerValidated,omitempty"` //
	ServerSoftware    int    `json:"serverSoftware,omitempty"`      //
	Organization      string `json:"organizationName,omitempty"`    //
	Address           string `json:"streetAddress1,omitempty"`      //
	City              string `json:"localityName,omitempty"`        //
	State             string `json:"stateOrProvinceName,omitempty"` //
	Postcode          string `json:"postalCode,omitempty"`          //
	Country           string `json:"countryName,omitempty"`         //
	DcvEmail          string `json:"dcvEmailAddress,omitempty"`     //
	Prioritise        string `json:"prioritiseCSRValues,omitempty"` //
	RepEmail          string `json:"appRepEmailAddress,omitempty"`  //

}

func (a SectigoSSL) GenerateOrder(params url.Values) (url.Values, error) {
	if a.Testmode == "on" || a.Testmode == "yes" {
		params.Set("test", "Y")
	} else {
		params.Set("test", "N")
	}
	_result, err := a.apiCall(params, "!AutoApplyOrder")
	return _result, err
}

func (a SectigoSSL) GetOrderStatus(params url.Values) (url.Values, error) {
	result, err := a.apiCall(params, "!GetDetailedOrderStatus")
	return result, err
}

func (a SectigoSSL) Refund(params url.Values) (url.Values, error) {
	result, err := a.apiCall(params, "!AutoRefund")
	return result, err
}

func (a SectigoSSL) CollectCertificate(params url.Values) (url.Values, error) {
	result, err := a.apiCall(params, "download/CollectSSL")
	return result, err
}

func (a SectigoSSL) UpdateDCV(params url.Values) (url.Values, error) {
	result, err := a.apiCall(params, "!AutoUpdateDCV")
	return result, err
}

func (a SectigoSSL) Reissuecertificate(params url.Values) (url.Values, error) {
	result, err := a.apiCall(params, "!AutoReplaceSSL")
	return result, err
}

func (a SectigoSSL) AccountBalance(params url.Values) (url.Values, error) {
	_result, err := a.apiCall(params, "!getAccountBalance")
	return _result, err
}

func (a SectigoSSL) GeneratePrivateKey() (string, error) {
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return "", err
	}
	privateKeyBytes := encodePrivateKeyToPEM(privateKey)
	return string(privateKeyBytes), nil
}

func (a SectigoSSL) GenerateCSR(params map[string]string) (string, string, error) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, bitSize)

	emailAddress := params["emailAddress"]
	subj := pkix.Name{
		CommonName:         params["commonName"],
		Country:            []string{params["country"]},
		Province:           []string{params["state"]},
		Locality:           []string{params["city"]},
		Organization:       []string{params["company"]},
		OrganizationalUnit: []string{params["unit"]},
	}
	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{Type: oidEmailAddress, Value: emailAddress},
	})

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	csr := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	privateKeyBytes := encodePrivateKeyToPEM(keyBytes)
	return string(csr), string(privateKeyBytes), nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func (a SectigoSSL) apiCall(request url.Values, postUri string) (url.Values, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	client := http.Client{
		Timeout:   100 * time.Second,
		Transport: tr,
	}
	reqUrl := liveUrl

	if postUri != "" {
		reqUrl = reqUrl + postUri
	} else {
		reqUrl = reqUrl + "!AutoApplySSL"
	}

	request.Set("loginName", a.Username)
	request.Set("loginPassword", a.Password)

	if a.Testmode == "on" {
		fmt.Printf("%s", request)
	}
	resp, err := client.PostForm(reqUrl, request)

	// check for response error
	if err != nil {
		return url.Values{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		decodeCSRResponse, _ := ioutil.ReadAll(resp.Body)
		if postUri != "download/CollectSSL" {
			result, err := url.ParseQuery(string(decodeCSRResponse))
			if err != nil {
				return url.Values{}, err
			}
			_statusCode := strings.Join(result["errorCode"], "")
			_errMessage := strings.Join(result["errorMessage"], "")
			if _statusCode != "0" {
				return url.Values{}, errors.New(_errMessage)
			}
			return result, nil
		} else {
			r := make(url.Values)
			s := string(decodeCSRResponse)
			r.Set("certificate", s[2:])
			return r, nil
		}
	}

	return url.Values{}, nil
}
