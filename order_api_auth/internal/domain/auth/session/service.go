package session

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	pkgLogger "order_api_auth/pkg/logger"
	"order_api_auth/pkg/sms"
	"order_api_auth/pkg/utils"
	"strconv"
)

type Service interface {
	SendSessionID(*Session) error
	VerifyCode(sessionID, code string) (jwt string, err error)
}

type ServiceSession struct {
	repository Repository
}

func NewService(repository Repository) *ServiceSession {
	return &ServiceSession{repository: repository}
}

func (s *ServiceSession) SendSessionID(session *Session) error {
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("error generating sessionID")
		return fmt.Errorf("error generating sessionID: %w", err)
	}

	smsCode := utils.GetFakeSMSCode()
	session.SessionID = sessionID
	session.SMSCode = strconv.Itoa(smsCode)

	fmt.Printf("%+v\n", *session)
	fmt.Printf("%+v\n", session)
	fmt.Printf("%d %s", smsCode, sessionID)

	err = s.repository.CreateSession(session)
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("error creating session in DB")
		return fmt.Errorf("error creating session in DB: %w", err)
	}

	smsErr := sms.SendFakeSMS(session.Phone, session.SMSCode)
	if smsErr != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error":      err,
			"session_id": session.SessionID,
		}).Error("error sending fake sms code")

		deleteErr := s.repository.DeleteSession(session.SessionID)
		if deleteErr != nil {
			pkgLogger.Logger.WithFields(logrus.Fields{
				"error":      deleteErr.Error(),
				"session_id": session.SessionID,
			}).Error("error deleting session in DB after SMS failure")
		}

		return fmt.Errorf("error sending fake sms code: %w",
			errors.Join(smsErr, deleteErr))
	}

	return nil
}

func (s *ServiceSession) VerifyCode(sessionID, code string) (jwt string, err error) {
	//TODO: Implement method
	panic("implement me")
}
