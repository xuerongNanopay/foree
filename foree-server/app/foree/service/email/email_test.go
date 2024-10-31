package foree_email_service

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	foree_email_template "xue.io/go-pay/app/foree/service/email/template"
)

func TestHowHTMLTemplateWorks(t *testing.T) {

	type ForeeRootTemplateData struct {
		AppName           string
		AppLink           string
		LogoImg           string
		Outlet            template.HTML
		SendTo            string
		SupportAddress    string
		PrivacyUrl        string
		PrivacyLabel      string
		TermsAndCondLink  string
		TermsAndCondLabel string
		ContactEmail      string
		AboutLink         string
	}

	type ForeeTransactionCancelledTemplateData struct {
		GreetingName      string
		TransactionNumber string
		TermsAndCondLink  string
		SupportEmail      string
	}

	t.Run("Demo 1", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		_ = template.ExecuteTemplate(buf, "T", "<script>alert('you have been pwned')</script>")
		fmt.Println(buf.String())
	})

	t.Run("Test Foree Root Template", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("Foree Root Template").Parse(foree_email_template.RootLayoutTemplate)

		data := ForeeRootTemplateData{
			AppName:           "Foree",
			AppLink:           "http://www.foree.net",
			LogoImg:           "http://www.foree.net/logo",
			Outlet:            `<b>World</b>`,
			SendTo:            "Xuerong",
			SupportAddress:    "support@foree.net",
			PrivacyUrl:        "http://www.foree.net/privacy_url",
			PrivacyLabel:      "Privancy",
			TermsAndCondLink:  "http://www.foree.net/terms",
			TermsAndCondLabel: "Terms And Condition",
			ContactEmail:      "contact@foree.net",
			AboutLink:         "http://www.foree.net/anout",
		}
		_ = template.Execute(buf, data)
		fmt.Println(buf.String())
	})

	t.Run("Test Foree Transaction Cancelled Template", func(t *testing.T) {
		buf := new(bytes.Buffer)
		txTpl, _ := template.New("Foree Transaction Cancelled Template").Parse(foree_email_template.TransactionCancelledTemplate)

		txData := ForeeTransactionCancelledTemplateData{
			GreetingName:      "Xuerong Wu",
			TransactionNumber: "NP000000000001",
			TermsAndCondLink:  "http://www.foree.net/terms",
			SupportEmail:      "support@foree.net",
		}

		_ = txTpl.Execute(buf, txData)
		o := buf.String()
		outlet := template.HTML(o)
		buf = new(bytes.Buffer)
		template, _ := template.New("Foree Root Template").Parse(foree_email_template.RootLayoutTemplate)
		data := ForeeRootTemplateData{
			AppName:           "Foree",
			AppLink:           "http://www.foree.net",
			LogoImg:           "http://www.foree.net/logo",
			SendTo:            "Xuerong",
			Outlet:            outlet,
			SupportAddress:    "support@foree.net",
			PrivacyUrl:        "http://www.foree.net/privacy_url",
			PrivacyLabel:      "Privancy",
			TermsAndCondLink:  "http://www.foree.net/terms",
			TermsAndCondLabel: "Terms And Condition",
			ContactEmail:      "contact@foree.net",
			AboutLink:         "http://www.foree.net/anout",
		}
		_ = template.Execute(buf, data)

		fmt.Println(buf.String())
	})
}
