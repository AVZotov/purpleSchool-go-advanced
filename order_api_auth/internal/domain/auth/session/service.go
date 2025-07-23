package session

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgLogger "order_api_auth/pkg/logger"
	"order_api_auth/pkg/sms"
	"order_api_auth/pkg/utils"
	"strconv"
)

type Service interface {
	CreateSession(*http.Request, *Session) error
	VerifySession(sessionID, code string) (jwt string, err error)
}

type ServiceSession struct {
	repository Repository
}

func NewService(repository Repository) *ServiceSession {
	return &ServiceSession{repository: repository}
}

func (s *ServiceSession) CreateSession(r *http.Request, session *Session) error {
	pkgLogger.InfoWithRequestID(r, "request for session passed to service layer", logrus.Fields{})

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "error generating sessionID", logrus.Fields{
			"error": err,
		})
		return fmt.Errorf("error generating sessionID: %w", err)
	}

	smsCode := utils.GetFakeSMSCode()
	session.SessionID = sessionID
	session.SMSCode = strconv.Itoa(smsCode)

	err = s.repository.CreateSession(r, session)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "error creating session in DB", logrus.Fields{
			"error": err,
		})
		return fmt.Errorf("error creating session in DB: %w", err)
	}

	smsErr := sms.SendFakeSMS(session.Phone, session.SMSCode)
	if smsErr != nil {
		pkgLogger.ErrorWithRequestID(r, "error sending fake sms code", logrus.Fields{
			"error": smsErr.Error(),
		})
		deleteErr := s.repository.DeleteSession(session.SessionID)
		if deleteErr != nil {
			pkgLogger.ErrorWithRequestID(r, "error sending fake sms code", logrus.Fields{
				"error": deleteErr.Error(),
			})
		}

		return fmt.Errorf("error sending fake sms code: %w",
			errors.Join(smsErr, deleteErr))
	}

	return nil
}

func (s *ServiceSession) VerifySession(sessionID, code string) (jwt string, err error) {
	//TODO: Implement method
	panic("implement me")
}
