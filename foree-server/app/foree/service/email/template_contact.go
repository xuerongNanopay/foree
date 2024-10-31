package foree_email_service

var contactAddTemplate = buildTemplate("CONTACT_ADDED", "Foree Remittance - Contact added", ContactAddHTML, rootLayoutTemplateHTML)

const ContactAddHTML = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Dear {{.CustomerName}}, </p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">This is a confirmation that {{.ContactName}} has been added as a Contact to your Foree Remittance account.</p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance</p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">If you think you have received this email in error, please contact us at: <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a></p>
</main>
`

var contactRemoveTemplate = buildTemplate("CONTACT_REMOVED", "Foree Remittance - Contact removed", contactRemoveHTML, rootLayoutTemplateHTML)

const contactRemoveHTML = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Dear {{.CustomerName}}, </p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">This is a confirmation that {{.ContactName}} has been removed from Contacts on your Foree Remittance account.</p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance</p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">If you think you have received this email in error, please contact us at: <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a></p>
</main>
`
