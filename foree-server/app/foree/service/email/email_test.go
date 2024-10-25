package foree_email_service

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestHowHTMLTemplateWorks(t *testing.T) {

	type ForeeRootTemplateData struct {
		AppName           string
		AppLink           string
		LogoImg           string
		Outlet            string
		SendTo            string
		SupportAddress    string
		PrivacyUrl        string
		PrivacyLabel      string
		TermsAndCondLink  string
		TermsAndCondLabel string
		ContactEmail      string
		AboutLink         string
	}

	t.Run("Demo 1", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		_ = template.ExecuteTemplate(buf, "T", "<script>alert('you have been pwned')</script>")
		fmt.Println(buf.String())
	})

	t.Run("Test Foree Root Template", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("Foree Root Template").Parse(ForeeRootTemplate)

		data := ForeeRootTemplateData{
			AppName:           "Foree",
			AppLink:           "http://www.foree.net",
			LogoImg:           "http://www.foree.net/logo",
			Outlet:            "<h1>AAA</h1>",
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
}
