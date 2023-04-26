package bwnotification

import (
	"encoding/json"
	"errors"
	nats "github.com/nats-io/nats.go"
	"os"

	validator "github.com/go-playground/validator/v10"
)

func (a BWNotification) Publish(params NotificationParams) error {
	validate := validator.New()
	err := validate.Var(params.Action, "required")

	if err != nil {
		return errors.New("action is missing")
	}

	if params.Action == "mail" {
		err := validate.Var(params.EmailData.PartnerId, "required")

		if err != nil {
			return errors.New("partner id is missing")
		}

		if !params.EmailData.IsTemplate {
			err = validate.Var(params.EmailData.Content, "required")

			if err != nil {
				return errors.New("content is missing")
			}

			err = validate.Var(params.EmailData.Subject, "required")

			if err != nil {
				return errors.New("subject is missing")
			}
		}

		if len(params.EmailData.Recipient) == 0 {
			return errors.New("recipient is missing")
		}
	}

	if params.Action == "sms" {
		err := validate.Var(params.APPData.PartnerId, "required")

		if err != nil {
			return errors.New("partner id is missing")
		}

		if len(params.SMSData.Recipient) == 0 {
			return errors.New("recipient is missing")
		}
	}

	if params.Action == "app" {
		err := validate.Var(params.APPData.PartnerId, "required")

		if err != nil {
			return errors.New("partner id is missing")
		}

		if !params.APPData.IsTemplate {
			err = validate.Var(params.APPData.Content, "required")

			if err != nil {
				return errors.New("content is missing")
			}

			err = validate.Var(params.APPData.Subject, "required")

			if err != nil {
				return errors.New("subject is missing")
			}
		}

		if params.Provider == "slack" && params.APPData.ChannelId == "" {
			return errors.New("channel is missing")
		}
	}

	marshalData, _ := json.Marshal(params)
	natsPublish("bwnotification", string(marshalData))
	return nil
}

func natsPublish(channel string, param string) {
	natsUrl := os.Getenv("NATS_SERVICE")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}
	nc, _ := nats.Connect(natsUrl)

	nc.Publish(channel, []byte(param))
	nc.Flush()

}

type BWNotification struct{}

type NotificationParams struct {
	Action    string    `json:"action,omitempty"`   //
	Provider  string    `json:"provider,omitempty"` //
	EmailData EMailData `json:"data,omitempty"`     //
	SMSData   SMSData   `json:"sms_data,omitempty"` //
	APPData   APPData   `json:"app_data,omitempty"` //
}

type EmailContact struct {
	Name  string `json:"name,omitempty"`  //
	Email string `json:"email,omitempty"` //
}

type EMailData struct {
	Recipient        []EmailContact         `json:"recipient,omitempty"`        //
	Cc               []EmailContact         `json:"cc,omitempty"`               //
	Bcc              []EmailContact         `json:"bcc,omitempty"`              //
	Content          string                 `json:"content,omitempty"`          //
	Subject          string                 `json:"subject,omitempty" `         //
	Sender           *EmailContact          `json:"sender,omitempty"`           //
	ReplyTo          *EmailContact          `json:"replyTo,omitempty"`          //
	IsTemplate       bool                   `json:"isTemplate,omitempty"`       //
	InternalTemplate bool                   `json:"internalTemplate,omitempty"` //
	TemplateId       string                 `json:"templateId,omitempty"`       //
	PartnerId        string                 `json:"partnerId,omitempty"`        //
	SettingsId       int                    `json:"settingsId,omitempty"`       //
	ClientId         int                    `json:"clientId,omitempty"`         //
	Params           map[string]interface{} `json:"params,omitempty"`           //
	Attachment       struct {
		Content string `json:"content,omitempty"` //
		Name    string `json:"name,omitempty"`    //
	} `json:"attachment,omitempty"` //

}

type SMSContact struct {
	Number string `json:"number,omitempty"` //
}

type SMSData struct {
	Recipient  []SMSContact           `json:"recipient,omitempty"`  //
	Content    string                 `json:"content,omitempty"`    //
	PartnerId  string                 `json:"partnerId,omitempty"`  //
	IsTemplate bool                   `json:"isTemplate,omitempty"` //
	TemplateId string                 `json:"templateId,omitempty"` //
	ClientId   int                    `json:"clientId,omitempty"`   //
	Params     map[string]interface{} `json:"params,omitempty"`     //
}

type APPData struct {
	ChannelId  string                 `json:"channelId,omitempty"`  //
	Subject    string                 `json:"subject,omitempty"`    //
	Content    string                 `json:"content,omitempty"`    //
	PartnerId  string                 `json:"partnerId,omitempty"`  //
	IsTemplate bool                   `json:"isTemplate,omitempty"` //
	TemplateId string                 `json:"templateId,omitempty"` //
	Params     map[string]interface{} `json:"params,omitempty"`     //
}
