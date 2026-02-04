package email

// EmailBodyTemplates holds the raw strings for our emails
type EmailBodyTemplates struct {
	PaymentSuccess string
	OTPEmail       string
}

// Global instance to access them
var Templates = EmailBodyTemplates{
	PaymentSuccess: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; background-color: #f1f5f9; font-family: 'Inter', -apple-system, sans-serif;">
    <table width="100%" border="0" cellpadding="0" cellspacing="0">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table class="card" width="600" border="0" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 32px; border: 1px solid #e2e8f0; box-shadow: 0 10px 25px rgba(0,0,0,0.03);">
                    
                    <tr>
                        <td align="center" style="padding: 50px 40px 30px 40px;">
                            <div style="font-size: 28px; letter-spacing: -1px;">
                                <span style="color: #6366f1; font-weight: 900;">KODING</span><span style="color: #a855f7; font-weight: 900;">KRAZE</span>
                            </div>
                            <div style="height: 3px; width: 40px; background: linear-gradient(to right, #6366f1, #a855f7); margin-top: 10px; border-radius: 10px;"></div>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding: 0 50px 35px 50px;">
                            <h1 style="color: #1e293b; font-size: 26px; font-weight: 800; margin: 0; line-height: 1.2;">Purchase Successful</h1>
                            <p style="color: #64748b; font-size: 15px; margin-top: 10px;">Your new content has been added to your library.</p>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 0 40px 40px 40px;">
                            <table width="100%" border="0" cellpadding="0" cellspacing="0" style="background-color: #f8fafc; border-radius: 24px; border: 1px solid #f1f5f9;">
                                <tr>
                                    <td style="padding: 10px 30px 30px 30px;">
                                        <table width="100%" border="0" cellpadding="0" cellspacing="0">
                                            {{content}}
                                        </table>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding: 0 40px 60px 40px;">
                            <a href="{{dashboardUrl}}" style="background: linear-gradient(to right, #6366f1, #a855f7); background-color: #6366f1; color: #ffffff; display: inline-block; padding: 18px 45px; border-radius: 16px; text-decoration: none; font-weight: 800; font-size: 16px;">Access My Library</a>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding: 30px; border-top: 1px solid #f1f5f9; background-color: #fcfcfc; border-bottom-left-radius: 32px; border-bottom-right-radius: 32px;">
                            <p style="color: #94a3b8; font-size: 12px; margin: 0; font-weight: 500;">&copy; 2026 KodingKraze. Building the future of dev.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`,
}
