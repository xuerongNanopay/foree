package foree_email_service

const EmailVerifyCodeTemplate EmailTemplate = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 24px; font-weight: 900; margin: 40px 0;">Verify your email address</p>

	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">Hello {{.CustomerName}},</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">Welcome to {{.AppName}}. We provide fast and low cost money transfers from Canada to Pakistan in collaboration with our payout partner, the National Bank of XXXXXå. It is great to have you on board!</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">Please verify your email using the following code:</p>
	<div style="width: 100%; font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif; margin: 20px 0; height: 48px; border-radius: 3px; font-size: 32px; color: gray; text-align: center; font-weight: bold; letter-spacing: 1rem;">
		<p style = "padding: 14px 0 0; margin: 0;">{{.EmailVerifyCode}}</p>
	</div>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">Your code will expire in 30 minutes.</p>
</main>
`
