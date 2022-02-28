package takedown

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Incidents []struct {
	ID                   string `json:"id"`
	GroupID              string `json:"group_id"`
	AttackURL            string `json:"attack_url"`
	ReportedURL          string `json:"reported_url"`
	IP                   string `json:"ip"`
	CountryCode          string `json:"country_code"`
	DateSubmitted        string `json:"date_submitted"`
	LastUpdated          string `json:"last_updated"`
	Region               string `json:"region"`
	TargetBrand          string `json:"target_brand"`
	Authgiven            string `json:"authgiven"`
	Host                 string `json:"host"`
	Registrar            string `json:"registrar"`
	CustomerLabel        string `json:"customer_label"`
	CustomerTag          string `json:"customer_tag"`
	DateAuthed           string `json:"date_authed"`
	StopMonitoringDate   string `json:"stop_monitoring_date"`
	Domain               string `json:"domain"`
	Language             string `json:"language"`
	DateFirstActioned    string `json:"date_first_actioned"`
	Escalated            string `json:"escalated"`
	AttackType           string `json:"attack_type"`
	DeceptiveDomainScore string `json:"deceptive_domain_score"`
	DomainRiskRating     string `json:"domain_risk_rating"`
	FinalOutage          string `json:"final_outage"`
	FinalResolved        string `json:"final_resolved"`
	FirstOutage          string `json:"first_outage"`
	FirstResolved        string `json:"first_resolved"`
	FwdOwner             string `json:"fwd_owner"`
	HasPhishingKit       string `json:"has_phishing_kit"`
	Hostname             string `json:"hostname"`
	HostnameDdssScore    string `json:"hostname_ddss_score"`
	EvidenceURL          string `json:"evidence_url"`
	DomainAttack         string `json:"domain_attack"`
	FalsePositive        bool   `json:"false_positive"`
	HostnameAttack       string `json:"hostname_attack"`
	MalwareCategory      string `json:"malware_category"`
	MalwareFamily        string `json:"malware_family"`
	ReportSource         string `json:"report_source"`
	Reporter             string `json:"reporter"`
	RevOwner             string `json:"rev_owner"`
	ReverseDNS           string `json:"reverse_dns"`
	ScreenshotURL        string `json:"screenshot_url"`
	StatusChangeUptime   string `json:"status_change_uptime"`
	Status               string `json:"status"`
	TargetedURL          string `json:"targeted_url"`
	SiteRiskRating       string `json:"site_risk_rating"`
	WhoisServer          string `json:"whois_server"`
	AuthorisationSource  string `json:"authorisation_source"`
	EscalationSource     string `json:"escalation_source"`
	RestartDate          string `json:"restart_date"`
	Managed              bool   `json:"managed"`
	DateEscalated        string `json:"date_escalated"`
}

type Netcraft struct {
	Username string `json:"username"` //
	Password string `json:"password"` //
}

var reqUrl = "https://takedown.netcraft.com/apis/get-info.php"

func (c Netcraft) GetInfo() (Incidents, error) {
	var _result Incidents
	client := &http.Client{}
	method := "POST"
	var request []byte

	//log.Println(reqUrl)
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(request))
	if err != nil {
		return _result, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, errDo := client.Do(req)

	if errDo != nil {
		return _result, errDo
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	errjson := data.Decode(&_result)
	if errjson != nil {
		log.Println(errjson.Error())
	}

	return _result, nil
}
