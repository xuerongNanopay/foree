package foree_email_service

const TransactionCreateTemplate EmailTemplate = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
	Dear {{.CustomerName}},
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
	Thank you for using Foree Remittance!
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
	Congratulations! Your payment to {{.ContactName}} in the amount of {{.Amount}} has been initiated. Your transaction reference # is "{{.TransactionNumber}}".
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
	Please provide the transaction reference to your contact for any queries. In case you have sent a cash pickup payment, this transaction reference number will be required to collect the funds along with a valid government issued photo ID.
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">To track your transaction:</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">1. Sign in to Foree Remittance app</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">2. You can have a quick view from the Recent Transactions on the dashboard OR from the menu select, Transactions</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">3. You can view the status of your transaction or click on a transaction for details</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Through the Foree Remittance Referral Rewards, you can earn a <b>$20 CAD</b> credit when you refer someone using your referral link following these easy steps:</p>
	<br>
	<ol type="1">
	<li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">While logged into the Foree app, click on 'Referral Rewards' from the left menu</p></li>
	<li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Click on 'Share' (if using the mobile app) and share it through your preferred means of communication with your contacts in Canada. If you're using the desktop browser version of the Foree app, click on 'Copy' to copy the referral link for sharing.</p></li>
	<li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Your referral will use your unique link to Sign-Up on Foree Remittance and receive <b>$20 CAD</b> credit which will be applied to their first transaction.</p></li>
	<li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">When they complete the first transaction to their beneficiary contact in Pakistan, you will receive <b>$20 CAD</b> as a reward balance which can be viewed from the 'Referral Rewards' menu. This balance will be applied towards your next transaction.</p></li>
	<li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">You can earn a <b>$20 CAD</b> credit for each referral that signs up using your link!</p></li>
	</ol>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">You can start earning rewards through this program right away!</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance </p>
	</p> 
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
	If you think you have received this email in error, please contact us at <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a>
	</p> 
</main>
`

const TransactionPickupTemplate = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
	Dear {{.CustomerName}},
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
	This is to confirm that funds against your transaction to {{.ContactName}} in the amount of {{.Amount}} are now available for pick-up by your contact, {{.ContactName}}. Your transaction reference # is "{{.TransactionNumber}}".
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
	Please provide the transaction reference to your contact for any queries and to be provided at the time of cash pick-up from NBP locations.
	</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">To track your transaction:</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">1. Sign in to Foree Remittance app</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">2. You can have a quick view from the Recent Transactions on the dashboard OR from the menu select, Transactions</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">3. You can view the status of your transaction or click on a transaction for details</p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance</p>
	</p> 
	<br>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
	If you think you have received this email in error, please contact us at <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a>
	</p> 
</main>
`

const TransactionCompletedTemplate = `
    <include template = 'header'>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
      Dear {{.CustomerName}},
    </p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px; line-height: 1.5; margin: 0;">
      This is to confirm that funds against your transaction with reference "{{.TransactionNumber}}" to {{.ContactName}} in the amount of {{.Amount}} have been paid.
    </p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">To track your transaction:</p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">1. Sign in to Foree Remittance app</p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">2. You can have a quick view from the Recent Transactions on the dashboard OR from the menu select, Transactions</p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">3. You can view the status of your transaction or click on a transaction for details</p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Through the Foree Remittance Referral Rewards, you can earn a <b>$20 CAD</b> credit when you refer someone using your referral link following these easy steps:</p>
    <br>
    <ol type="1">
      <li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">While logged into the Foree app, click on 'Referral Rewards' from the left menu</p></li>
      <li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Click on 'Share' (if using the mobile app) and share it through your preferred means of communication with your contacts in Canada. If you're using the desktop browser version of the Foree app, click on 'Copy' to copy the referral link for sharing.</p></li>
      <li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Your referral will use your unique link to Sign-Up on Foree Remittance receive <b>$20 CAD</b> credit which will be applied to their first transaction.</p></li>
      <li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">When they complete the first transaction to their beneficiary contact in Pakistan, you will receive <b>$20 CAD</b> as a reward balance which can be viewed from the 'Referral Rewards' menu. This balance will be applied towards your next transaction.</p></li>
      <li><p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">You can earn a <b>$20 CAD</b> credit for each referral that signs up using your link!</p></li>
    </ol>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">You can start earning rewards through this program right away!</p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Thanks, </p>
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">Team Foree Remittance</p>
    </p> 
    <br>
    <p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
      If you think you have received this email in error, please contact us at <a href="mailto:{{.SupportEmail}}" style="text-decoration: none;">{{.SupportEmail}}</a>
    </p> 
  </include>
`

const TransactionCancelledTemplate = `
<main>
	<p style="font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;  font-size: 16px;  line-height: 1.5; margin: 0;">
		Dear {{.CustomerName}},
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
