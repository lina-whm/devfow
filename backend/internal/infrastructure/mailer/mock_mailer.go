package mailer

import (
	"context"
	"log/slog"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type MockMailer struct {
	SentEmails []Email
}

type Email struct {
	To      string
	Subject string
	Body    string
}

func NewMockMailer() *MockMailer {
	return &MockMailer{SentEmails: make([]Email, 0)}
}

func (m *MockMailer) Send(ctx context.Context, to, subject, body string) error {
	m.SentEmails = append(m.SentEmails, Email{To: to, Subject: subject, Body: body})
	slog.Info("mock email sent", "to", to, "subject", subject)
	return nil
}