package session

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	pkgJWT "order_api_auth/pkg/jwt"
	pkgLogger "order_api_auth/pkg/logger"
	"order_api_auth/pkg/sms"
	"order_api_auth/pkg/utils"
	"strconv"
)

type Service interface {
	CreateSession(*http.Request, *Session) error
	VerifySession(r *http.Request, session *Session) (string, error)
}

type ServiceSession struct {
	repository Repository
	secret     string
}

func NewService(repository Repository, secret string) *ServiceSession {
	return &ServiceSession{repository: repository, secret: secret}
}

func (s *ServiceSession) CreateSession(r *http.Request, session *Session) error {
	pkgLogger.InfoWithRequestID(r, "request for session passed to service layer", logrus.Fields{})

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, ErrGeneratingSessionID.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return errors.Join(ErrGeneratingSessionID, err)
	}

	smsCode := utils.GetFakeSMSCode()
	session.SessionID = sessionID
	session.SMSCode = strconv.Itoa(smsCode)

	err = s.repository.CreateSession(r, session)
	if err != nil {
		return errors.Join(ErrCreatingSession, err)
	}

	smsErr := sms.SendFakeSMS(session.Phone, session.SMSCode)
	if smsErr != nil {
		pkgLogger.ErrorWithRequestID(r, ErrSendingSMS.Error(), logrus.Fields{
			"error": smsErr.Error(),
		})
		deleteErr := s.repository.DeleteSession(r, session)
		if deleteErr != nil {
			pkgLogger.ErrorWithRequestID(r, ErrDeletingSession.Error(), logrus.Fields{
				"error": deleteErr.Error(),
			})
		}
		return errors.Join(ErrSendingSMS, smsErr, ErrDeletingSession, deleteErr)
	}

	return nil
}

func (s *ServiceSession) VerifySession(r *http.Request, session *Session) (string, error) {
	pkgLogger.InfoWithRequestID(r, "request for validation passed to service layer", logrus.Fields{})

	var requestedSession Session
	if err := s.repository.GetSession(r, &requestedSession); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Join(ErrSessionNotFound, err)
		}
		return "", errors.Join(ErrDBInternalError, err)
	}

	if session.SMSCode != requestedSession.SMSCode {
		pkgLogger.ErrorWithRequestID(r, ErrInvalidSMSCode.Error(), logrus.Fields{
			"error": ErrInvalidSMSCode.Error(),
		})
		return "", ErrInvalidSMSCode
	}

	jwtString, err := pkgJWT.Create(s.secret, requestedSession.Phone)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, ErrGeneratingToken.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return "", errors.Join(ErrInternalError, ErrGeneratingToken)
	}

	return jwtString, nil
}
