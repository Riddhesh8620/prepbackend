package handlers

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type otpVerifyReq struct {
	Email string `json:"email"` // React sends { "email": "..." }
	Code  string `json:"code"`  // React sends { "code": "..." }
}
type otpLoginReq struct {
	TargetEmail string `json:"email"`
}

// Simple in-memory store (Use Redis for production!)
var otpStore = make(map[string]string)
var mu sync.Mutex

// Generate a random 6-digit string
func generateOTP() string {
	max := 6
	b := make([]byte, max)
	n, _ := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "123456"
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func SendOTPEmail(c *fiber.Ctx) error {
	var body otpLoginReq

	GMAIL_ID := os.Getenv("GMAIL_ID")
	APP_PASSWORD := os.Getenv("APP_PASSWORD")

	otp := generateOTP()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Store OTP with the email as the key
	mu.Lock()
	otpStore[body.TargetEmail] = otp
	mu.Unlock()

	// Gmail Config
	from := GMAIL_ID
	pass := APP_PASSWORD
	subject := "Subject: 🔑 Your KodingKraze Access Code\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7ff; padding: 40px 20px;">
  <div style="max-width: 460px; margin: 0 auto; background: #ffffff; border-radius: 24px; overflow: hidden; box-shadow: 0 20px 40px rgba(79, 70, 229, 0.1); border: 1px solid #e5e7eb;">
    
    <div style="background: linear-gradient(135deg, #8b5cf6 0%%, #6d28d9 100%%); padding: 40px 30px; text-align: center;">
      <div style="background: rgba(255, 255, 255, 0.2); width: 60px; height: 60px; line-height: 60px; border-radius: 16px; margin: 0 auto 20px; color: white; font-size: 32px; font-weight: bold; border: 1px solid rgba(255, 255, 255, 0.3);">K</div>
      <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 700; letter-spacing: -0.5px;">KodingKraze6</h1>
    </div>

    <div style="padding: 40px 30px; text-align: center;">
      <h2 style="color: #1f2937; font-size: 22px; font-weight: 700; margin-bottom: 12px;">Verification Code</h2>
      <p style="color: #6b7280; font-size: 15px; line-height: 24px; margin-bottom: 32px;">
        Use the code below to securely sign in to your premium learning dashboard. It's valid for 10 minutes.
      </p>
      
      <div style="background-color: #f5f3ff; border: 2px dashed #ddd6fe; padding: 24px; border-radius: 20px; margin-bottom: 32px;">
        <span style="font-size: 42px; font-weight: 800; letter-spacing: 12px; color: #7c3aed; font-family: 'Courier New', monospace;">
          %s
        </span>
      </div>
      
      <p style="color: #9ca3af; font-size: 13px; margin: 0;">
        Didn't request this? Please ignore this email or contact support.
      </p>
    </div>

    <div style="background: #f9fafb; padding: 20px; text-align: center; border-top: 1px solid #f3f4f6;">
      <p style="color: #a5b4fc; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; margin: 0;">
        Premium Courses • Real Skills • Your Pace
      </p>
    </div>
  </div>
</body>
</html>`, otp)
	msg := []byte(subject + mime + htmlBody)

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{body.TargetEmail}, []byte(msg))
}

func VerifyOTP(c *fiber.Ctx) error {
	var body otpVerifyReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	mu.Lock()
	defer mu.Unlock() // Use defer to ensure unlock happens even if we return early

	// Check if the email exists in the map first
	if _, exists := otpStore[body.Email]; !exists {
		return c.Status(401).JSON(fiber.Map{"error": "No OTP found for this email"})
	}

	// Now delete it
	delete(otpStore, body.Email)
	return LoginInternal(&body.Email, c)
}
