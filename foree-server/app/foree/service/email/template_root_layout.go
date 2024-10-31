package foree_email_service

const rootLayoutTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1" />
	<title>{{.AppName}}</title>
	<style type="text/css" media="screen">
	@media screen {
		* { font-family: /*%FONT1%*/ Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif !important; }
	}
	</style>
	<style type="text/css">
	* {
		font-family: /*%FONT1%*/, Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif;
	}
	body {
		-webkit-font-smoothing: antialiased;
		-webkit-text-size-adjust: none;
		width: 100%;
		margin: 0;
		height: 100%;
	}
	a {
		text-decoration: none !important;
	}
	p {
		margin: 0;
	}
	.p-body {
		font-size: 1rem;
		line-height: 1.5;
		font-style: normal;
		margin:  0;
	}
	.subheading {
		font-size: 1.25rem;
		line-height: 1.6;
		font-style: normal;
		font-weight: 700;
		margin:  0;
	}
	.legal {
		font-size: 0.75rem;
		color: #8D8D8D;
		text-align: center;
	}
	.link {
		color: #406DEA;
	}
	.wrapper {
		padding: 0.4rem;
	}
	@media only screen and (min-width: 768px) {
		.wrapper {
		padding: 1.5rem;
		}
	}
	</style>
</head>
<body style = "background: #EAECED;font-family: 'Source Sans Pro', Roboto, 'Helvetica Neue', Helvetica, Arial, sans-serif; padding:1rem;">
		<div class="wrapper" style="background: white; border-radius: 0.25rem; margin-bottom: 1.25rem;" >
		<div class="appLogo" style="margin-bottom: 1.5rem;">
			<a href="{{.AppLink}}" title="{{.AppName}} Logo" target="_blank" style="cursor:default"><img src= "{{.LogoImg}}"  alt="{{.AppName}} Logo" style="max-width: 115px;" /></a>
		</div>

		<div class="content">
			{{.Outlet}}
		</div>
		</div>

		<div style="text-align:center">
		<div style="display:inline-block;">
			<a href="https://facebook.com/foree.remit" title="Facebook logo" target="_blank" style="cursor:default">
			<img src="{{.AppLink}}/images/media/facebook-logo.png"  alt="Foree Facebook" style="max-width: 115px;" />
			</a>

			<a href="https://www.instagram.com/foree_remit/" title="Instagram logo" target="_blank" style="cursor:default">
			<img src="{{.AppLink}}/images/media/instagram-logo.png"  alt="Foree Instagram" style="max-width: 115px;" />
			</a>

			<a href="https://www.linkedin.com/company/foree-remit/" title="LinkedIn logo" target="_blank" style="cursor:default">
			<img src="{{.AppLink}}/images/media/linkedin-logo.png"  alt="Foree LinkedIn" style="max-width: 115px;" />
			</a>
		</div>

		<div class="legal">
			<div style="margin: 4px;">To manage notification settings, please <a href="{{.AppLink}}/#notification-settings">click here</a>.</div>
			<div style="margin: 4px;">This message was sent to {{.SendTo}}</div>
			<div style="margin: 4px;">© Foree Remittance, {{.SupportAddress}}</div>
		</div>

		<div style="text-align:center;">
			<a href="{{.PrivacyUrl}}" >{{.PrivacyLabel}}</a> |
			<a href="{{.TermsAndCondLink}}" >{{.TermsAndCondLabel}}</a>
			<br>
			<a href="mailto:{{.ContactEmail}}" > Contact Us </a> |
			<a href="{{.AboutLink}}" > About Us </a>
		</div>

		</div>
</body>
</html>
`
