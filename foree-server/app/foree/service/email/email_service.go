package foree_email_service

type EmailTemplate string

type BasicTplCfg struct {
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

func NewEmailService() *EmailService {
	ret := &EmailService{}
	ret.compileTemplates()
	return ret
}

type EmailService struct {
	tplCfg BasicTplCfg
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
