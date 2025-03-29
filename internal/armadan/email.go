package armadan

const RESET_PASSWORD_EMAIL_TEMPLATE_TEXT = `
	Återställ ditt lösenord

	Vi har fått en begäran om att återställa ditt lösenord. Följ länken nedan för att välja ett nytt lösenord:

	%s

	Om du inte begärde en lösenordsåterställning kan du ignorera detta meddelande.

	Med vänliga hälsningar,  
	Armadan support
`

const RESET_PASSWORD_EMAIL_TEMPLATE_HTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Återställ ditt lösenord</title>
</head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
    <div class="container" style="max-width: 600px; background: #ffffff; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
        <h2>Återställ ditt lösenord</h2>
        <p>Vi har fått en begäran om att återställa ditt lösenord. Klicka på knappen nedan för att välja ett nytt lösenord.</p>
        <a href="%s" class="button" style="display: inline-block; padding: 10px 20px; color: #ffffff; background: #000000; text-decoration: none; border-radius: 5px;">Återställ lösenord</a>
        <p>Om du inte begärde en lösenordsåterställning kan du ignorera detta meddelande.</p>
        <div class="footer" style="margin-top: 20px; font-size: 12px; color: #666666;">
            <p>Med vänliga hälsningar,<br>Armadan support</p>
        </div>
    </div>
</body>
</html>
`

type EmailService interface {
	SendResetPassword(string, string) error
}

type Senders struct {
	ResetPassword string
}
