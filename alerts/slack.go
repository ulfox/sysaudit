package alerts

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

type KDAuditSlackAlert struct {
	Webhook string
}

func NewSlackAlert(w string) KDAuditSlackAlert {
	kdAuditSlackAlert := KDAuditSlackAlert{
		Webhook: w,
	}

	return kdAuditSlackAlert
}

func (k *KDAuditSlackAlert) SendSlackNotification(msg string) error {
	if k.Webhook == "" {
		return nil
	}

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, k.Webhook, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	slackClient := &http.Client{}

	resp, err := slackClient.Do(req)
	if err != nil {
		return err
	}

	readBodyResponse := new(bytes.Buffer)
	readBodyResponse.ReadFrom(resp.Body)

	if readBodyResponse.String() != "ok" {
		return errors.New(resp.Status)
	}
	return nil
}
