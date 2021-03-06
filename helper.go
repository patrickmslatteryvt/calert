package main

import (
	"bytes"
	"html/template"
	"strings"

	alerttemplate "github.com/prometheus/alertmanager/template"
	"github.com/spf13/viper"
)

func sendMessageToChat(alerts []alerttemplate.Alert, notif *Notifier, webHookURL string) error {
	var (
		message = ChatNotification{}
		str     strings.Builder
	)
	// common funcs used in template
	templateFuncMap := template.FuncMap{
		"Title":   strings.Title,
		"toUpper": strings.ToUpper,
	}
	// read template file
	tmpl, err := template.New("message.tmpl").Funcs(templateFuncMap).ParseFiles(viper.GetString("app.template_file"))
	if err != nil {
		errLog.Printf("Error reading template %s", err)
		return err
	}
	// loop through list of alerts and append the data in template
	for _, a := range alerts {
		var to bytes.Buffer
		err = tmpl.Execute(&to, a)
		if err != nil {
			errLog.Printf("Error parsing values in template %s", err)
			return err
		}
		str.WriteString(to.String())
		str.WriteString("\n")
	}
	// prepare request payload for Google chat webhook endpoint
	message.Text = str.String()

	return notif.PushNotification(message, webHookURL)
}
