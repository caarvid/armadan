package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/caarvid/armadan/internal/armadan"
	"golang.org/x/time/rate"
)

const dailyLimit = 1000

var (
	dailyCalls int
	dailyReset = time.Now().Truncate(time.Hour * 24).Add(time.Hour * 24)
	limiter    = rate.NewLimiter(rate.Every(time.Second), 14)
	mu         sync.Mutex
)

func resetDailyLimit() {
	mu.Lock()
	defer mu.Unlock()

	if time.Now().After(dailyReset) {
		dailyCalls = 0
		dailyReset = time.Now().Truncate(time.Hour * 24).Add(time.Hour * 24)
	}
}

func sendEmail(email, sender, subject, htmlBody, textBody string, awsConfig aws.Config) error {
	resetDailyLimit()

	mu.Lock()
	if dailyCalls >= dailyLimit {
		mu.Unlock()
		return errors.New("daily rate limit exceeded")
	}

	dailyCalls += 1
	mu.Unlock()

	if err := limiter.Wait(context.Background()); err != nil {
		return errors.New("rate limiter error")
	}

	svc := sesv2.NewFromConfig(awsConfig)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(sender),
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("utf-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(htmlBody),
						Charset: aws.String("utf-8"),
					},
					Text: &types.Content{
						Data:    aws.String(textBody),
						Charset: aws.String("utf-8"),
					},
				},
			},
		},
	}

	fmt.Printf("Sending email to %s from %s", email, sender)

	_, err := svc.SendEmail(context.Background(), input)

	return err
}

type emailService struct {
	senders       armadan.Senders
	awsConfig     aws.Config
	emailOverride string
}

func NewEmailService(senders armadan.Senders, cfg aws.Config, emailOverride string) *emailService {
	return &emailService{senders: senders, awsConfig: cfg, emailOverride: emailOverride}
}

func (es *emailService) SendResetPassword(email, token string) error {
	subject := "Armadan - Återställ ditt lösenord"
	htmlBody := fmt.Sprintf(armadan.RESET_PASSWORD_EMAIL_TEMPLATE_HTML, token)
	textBody := fmt.Sprintf(armadan.RESET_PASSWORD_EMAIL_TEMPLATE_TEXT, token)

	if len(es.emailOverride) > 0 {
		email = es.emailOverride
	}

	return sendEmail(email, es.senders.ResetPassword, subject, htmlBody, textBody, es.awsConfig)
}
