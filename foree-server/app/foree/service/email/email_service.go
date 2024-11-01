package foree_email_service

import (
	"bytes"
	"html/template"
)

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
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService() *EmailService {
	ret := &EmailService{}
	ret.templates = map[string]emailTemplate{
		"EMAIL_VERIFY_CODE":     buildTemplate("EMAIL_VERIFY_CODE", "Verify your email to activate your Foree Remittance account", emailVerifyCodeHTML, rootLayoutTemplateHTML),
		"CONTACT_ADDED":         buildTemplate("CONTACT_ADDED", "Foree Remittance - Contact added", ContactAddHTML, rootLayoutTemplateHTML),
		"CONTACT_REMOVED":       buildTemplate("CONTACT_REMOVED", "Foree Remittance - Contact removed", contactRemoveHTML, rootLayoutTemplateHTML),
		"TRANSACTION_INITIATED": buildTemplate("TRANSACTION_INITIATED", "Foree Remittance - New Transaction Initiated", transactionInitiatedHTML, rootLayoutTemplateHTML),
		"TRANSACTION_PICKUP":    buildTemplate("TRANSACTION_PICKUP", "Foree Remittance - Your cash transaction is available for pick-up", transactionPickupHTML, rootLayoutTemplateHTML),
		"TRANSACTION_COMPLETED": buildTemplate("TRANSACTION_COMPLETED", "Foree Remittance - Transaction completed", transactionCompletedHTML, rootLayoutTemplateHTML),
		"TRANSACTION_CANCELLED": buildTemplate("TRANSACTION_CANCELLED", "Foree Remittance - Transaction cancelled", transactionCancelledHTML, rootLayoutTemplateHTML),
	}

	return ret
}

type EmailService struct {
	basicTemplateCfg BasicTemplateCfg
	serviceConfig    ServiceConfig
	templates        map[string]emailTemplate
}

func (e *EmailService) sendEmail(name string, data templateData) error {
	tpl, _ := e.templates[name]

	buf := new(bytes.Buffer)
	err := tpl.contentTpl.Execute(buf, data)
	if err != nil {
		return err
	}

	data.Content = template.HTML(buf.String())
	buf.Reset()
	err = tpl.layoutTpl.Execute(buf, data)
	if err != nil {
		return err
	}

	emailBody := buf.String()
	if emailBody == "" {
		//TODO: send email use goEmail
	}

	return nil
}

// We can put all template variable into one struct
type templateData struct {
	Content           template.HTML
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
	CustomerName      string
	TransactionNumber string
	SupportEmail      string
	EmailVerifyCode   string
	ContactName       string
	Amount            string
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
