package foree_email_service

import (
	"bytes"
	"fmt"
	"html/template"
)

//Check this one for attachment: https://gist.github.com/douglasmakey/90753ecf37ac10c25873825097f46300

type EmailMsg struct {
	From        string
	To          string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

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
	SupportEmail      string
}

type ServiceConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

const (
	emailVerifyCodeTPLName      = "EMAIL_VERIFY_CODE"
	contactAddedKeyTPLName      = "CONTACT_ADDED"
	contactRemovedKeyTPLName    = "CONTACT_REMOVED"
	transactionInitiatedTPLName = "TRANSACTION_INITIATED"
	transactionPickupTPLName    = "TRANSACTION_PICKUP"
	transactionCompletedTPLName = "TRANSACTION_COMPLETED"
	transactionCancelledTPLName = "TRANSACTION_CANCELLED"
)

func NewEmailService() *EmailService {
	ret := &EmailService{}
	ret.templates = map[string]emailTemplate{
		emailVerifyCodeTPLName:      buildTemplate(emailVerifyCodeTPLName, "Verify your email to activate your Foree Remittance account", emailVerifyCodeHTML, rootLayoutTemplateHTML),
		contactAddedKeyTPLName:      buildTemplate(contactAddedKeyTPLName, "Foree Remittance - Contact added", ContactAddHTML, rootLayoutTemplateHTML),
		contactRemovedKeyTPLName:    buildTemplate(contactRemovedKeyTPLName, "Foree Remittance - Contact removed", contactRemoveHTML, rootLayoutTemplateHTML),
		transactionInitiatedTPLName: buildTemplate(transactionInitiatedTPLName, "Foree Remittance - New Transaction Initiated", transactionInitiatedHTML, rootLayoutTemplateHTML),
		transactionPickupTPLName:    buildTemplate(transactionPickupTPLName, "Foree Remittance - Your cash transaction is available for pick-up", transactionPickupHTML, rootLayoutTemplateHTML),
		transactionCompletedTPLName: buildTemplate(transactionCompletedTPLName, "Foree Remittance - Transaction completed", transactionCompletedHTML, rootLayoutTemplateHTML),
		transactionCancelledTPLName: buildTemplate(transactionCancelledTPLName, "Foree Remittance - Transaction cancelled", transactionCancelledHTML, rootLayoutTemplateHTML),
	}

	return ret
}

type EmailService struct {
	basicTemplateCfg BasicTemplateCfg
	serviceConfig    ServiceConfig
	templates        map[string]emailTemplate
}

func (e *EmailService) buildEmailMsg(tplName string, data templateData, from, to string) (*EmailMsg, error) {
	data.AppName = e.basicTemplateCfg.AppName
	data.AppLink = e.basicTemplateCfg.AppLink
	data.LogoImg = e.basicTemplateCfg.LogoImg
	data.SendTo = e.basicTemplateCfg.SendTo
	data.SupportAddress = e.basicTemplateCfg.SupportAddress
	data.PrivacyUrl = e.basicTemplateCfg.PrivacyUrl
	data.PrivacyLabel = e.basicTemplateCfg.PrivacyLabel
	data.TermsAndCondLink = e.basicTemplateCfg.TermsAndCondLink
	data.TermsAndCondLabel = e.basicTemplateCfg.TermsAndCondLabel
	data.ContactEmail = e.basicTemplateCfg.ContactEmail
	data.AboutLink = e.basicTemplateCfg.AboutLink

	tpl, ok := e.templates[tplName]
	if !ok {
		return nil, fmt.Errorf("email template `%v` not found", tplName)
	}

	buf := new(bytes.Buffer)
	err := tpl.contentTpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	data.Content = template.HTML(buf.String())
	buf.Reset()
	err = tpl.layoutTpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	emailBody := buf.String()

	return &EmailMsg{
		From:    from,
		To:      to,
		Subject: tpl.subject,
		Body:    emailBody,
	}, nil
}

func (e *EmailService) sendWithTemplate(tplName string, data templateData, from, to string) error {
	eMsg, err := e.buildEmailMsg(tplName, data, from, to)
	if err != nil {
		return err
	}
	if err := e.send(eMsg.Subject, eMsg.Body, eMsg.From, eMsg.To); err != nil {
		return err
	}
	return nil
}

func (e *EmailService) send(subject, body, from, to string) error {
	return nil
}

func (e *EmailService) SendEmailVerifyCode(emailVerifyCode, from, to string) error {
	data := templateData{
		EmailVerifyCode: emailVerifyCode,
	}
	if err := e.sendWithTemplate(emailVerifyCodeTPLName, data, from, to); err != nil {
		return err
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
