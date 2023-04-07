package auth

import (
	"context"
	"fmt"
)

type SMSLogSender struct {
}

// Send an SMS
func (s SMSLogSender) Send(_ context.Context, number, text string) error {
	fmt.Println("sms sent to:", number, "contents:", text)
	return nil
}
