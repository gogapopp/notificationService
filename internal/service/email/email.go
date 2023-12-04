package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogapopp/notificationService/internal/models"
	"go.uber.org/zap"
)

type repository interface {
	GetSubscribedUsers(ctx context.Context) ([]models.UserSub, error)
}

type MailService struct {
	repo repository
	log  *zap.SugaredLogger
}

func NewMailService(repo repository, log *zap.SugaredLogger) *MailService {
	return &MailService{repo: repo, log: log}
}

func (m *MailService) SendEMails(msgBytes []byte) {
	const op = "service.email.SendEmails"
	var msg models.Message
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		m.log.Errorf("%s: error unmarshalling message:", op, err)
		return
	}
	users, err := m.repo.GetSubscribedUsers(context.Background())
	if err != nil {
		m.log.Errorf("%s: error getting subscribed users:", op, err)
		return
	}
	for _, user := range users {
		err = m.sendEmail(user.Email, msg)
		if err != nil {
			m.log.Errorf("%s: error sending email:", op, err)
		}
	}
}

func (m *MailService) sendEmail(to string, msg models.Message) error {
	fmt.Println(to, msg.Message, msg.Timestamp, msg.UserID)
	return nil
}
