package foree_email_service

import "html/template"

type emailTemplate struct {
	name       string
	subject    string
	contentTpl *template.Template
	layoutTpl  *template.Template
}

type BasicTemplateCfg struct {
	AppName           string
	AppLink           string
	LogoImg           string
	SendTo            string
	SupportAddress    string
	PrivacyUrl        string
	PrivacyLabel      string
	TermsAndCondLink  string
	TermsAndCondLabel string
	ContactEmail      string
	AboutLink         string
}

type ServiceConfig struct {
}

func NewEmailService() *EmailService {
	ret := &EmailService{}
	ret.compileTemplates()
	return ret
}

type EmailService struct {
	basicTemplateCfg BasicTemplateCfg
	templates        map[string]emailTemplate
}

func (e *EmailService) sendEmail(template emailTemplate) {

}

func (e *EmailService) compileTemplates() {

}

func (e *EmailService) EmailTransactionCancelled(greetingName, transactionNumber string) {

}

func (e *EmailService) Email(templateName string, data any) error {
	return nil
}

func (e *EmailService) EmailAsync(templateName string, data any) {

}

func buildTemplate(name, subject, content, layout string) emailTemplate {
	contentTpl, err := template.New(name).Parse(content)
	if err != nil {
		panic(err)
	}
	layoutTpl, err := template.New("").Parse(layout)
	if err != nil {
		panic(err)
	}
	return emailTemplate{
		name:       name,
		subject:    subject,
		contentTpl: contentTpl,
		layoutTpl:  layoutTpl,
	}
}
