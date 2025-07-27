package session

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	pkgJWT "order_api_auth/pkg/jwt"
	pkgLogger "order_api_auth/pkg/logger"
	"order_api_auth/pkg/sms"
	"order_api_auth/pkg/utils"
	"strconv"
)

type Service interface {
	CreateSession(context.Context, *Session) error
	VerifySession(context.Context, *Session) (string, error)
}

type ServiceSession struct {
	repository Repository
	secret     string
}

func NewService(repository Repository, secret string) *ServiceSession {
	return &ServiceSession{repository: repository, secret: secret}
}

func (s *ServiceSession) CreateSession(ctx context.Context, session *Session) error {
	pkgLogger.InfoWithRequestID(ctx, "request for session passed to service layer", logrus.Fields{})

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrGeneratingSessionID.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return errors.Join(ErrGeneratingSessionID, err)
	}

	smsCode := utils.GetFakeCode()
	session.SessionID = sessionID
	session.SMSCode = smsCode

	err = s.repository.CreateSession(ctx, session)
	if err != nil {
		return errors.Join(ErrCreatingSession, err)
	}

	smsErr := sms.SendFakeSMS(session.Phone, strconv.Itoa(session.SMSCode))
	if smsErr != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrSendingSMS.Error(), logrus.Fields{
			"error": smsErr.Error(),
		})
		deleteErr := s.repository.DeleteSession(ctx, session)
		if deleteErr != nil {
			pkgLogger.ErrorWithRequestID(ctx, ErrDeletingSession.Error(), logrus.Fields{
				"error": deleteErr.Error(),
			})
		}
		return errors.Join(ErrSendingSMS, smsErr, ErrDeletingSession, deleteErr)
	}

	return nil
}

func (s *ServiceSession) VerifySession(ctx context.Context, session *Session) (string, error) {
	pkgLogger.InfoWithRequestID(ctx, "request for validation passed to service layer", logrus.Fields{})

	var requestedSession Session
	if err := s.repository.GetSession(ctx, &requestedSession, session.SessionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Join(ErrSessionNotFound, err)
		}
		return "", errors.Join(ErrDBInternalError, err)
	}

	if session.SMSCode != requestedSession.SMSCode {
		pkgLogger.ErrorWithRequestID(ctx, ErrInvalidSMSCode.Error(), logrus.Fields{
			"error": ErrInvalidSMSCode.Error(),
		})
		return "", ErrInvalidSMSCode
	}

	if err := s.repository.DeleteSession(ctx, &requestedSession); err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrDeletingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})
	}

	jwtString, err := pkgJWT.Create(s.secret, requestedSession.Phone)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrGeneratingToken.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return "", errors.Join(ErrInternalError, ErrGeneratingToken)
	}

	return jwtString, nil
}
