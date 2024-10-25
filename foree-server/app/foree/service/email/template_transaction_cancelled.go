package foree_email_service

const ForeeTransactionCancelledTemplate = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
		Dear {{.GreetingName}},
	</p> 
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
		This is a confirmation that your transaction with transaction reference "{{.TransactionNumber}}" has been cancelled. Please refer to Foree Remittance 
		<a href="{{.TermsAndCondLink}}" target="_blank">Terms of Service</a> for cancellations and refunds.
	</p>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance </p>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
		If you think you have received this email in error, please contact us at <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a>
	</p>
</main>
`
