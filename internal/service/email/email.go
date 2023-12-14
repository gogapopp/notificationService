package email

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/gogapopp/notificationService/internal/models"
	"go.uber.org/zap"
)

type repository interface {
	GetSubscribedUsers(ctx context.Context) ([]models.UserSub, error)
}

type MailService struct {
	repo   repository
	config *config.Config
	auth   smtp.Auth
	log    *zap.SugaredLogger
}

func NewMailService(repo repository, config *config.Config, log *zap.SugaredLogger) *MailService {
	emailAuth := smtp.PlainAuth("", config.SMTP.Login, config.SMTP.Password, config.SMTP.Provider)
	return &MailService{repo: repo, config: config, auth: emailAuth, log: log}
}

func (m *MailService) SendEMails(msgBytes []byte) {
	const op = "service.email.SendEmails"
	var msg models.Message
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		m.log.Errorf("%v: error unmarshalling message:", op, err)
		return
	}
	users, err := m.repo.GetSubscribedUsers(context.Background())
	if err != nil {
		m.log.Errorf("%v: error getting subscribed users:", op, err)
		return
	}
	for _, user := range users {
		err = m.sendEmail(user.Email, msg)
		if err != nil {
			m.log.Errorf("%v: error sending email:", op, err)
		}
	}
}

func (m *MailService) sendEmail(to string, msg models.Message) error {
	addr := fmt.Sprintf("%s:%s", m.config.SMTP.Provider, m.config.SMTP.Port)
	subject := "Subject: Hello\n"
	body := fmt.Sprintf("To: %s\r\n%s\r\n%s", to, subject, msg.Message)
	err := smtp.SendMail(addr, m.auth, m.config.SMTP.Login, []string{to}, []byte(body))
	return err
}
