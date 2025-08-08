package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_api_cart/pkg/db/models"
	pkgErr "order_api_cart/pkg/errors"
	pkgJWT "order_api_cart/pkg/jwt"
	pkgLog "order_api_cart/pkg/logger"
	"order_api_cart/pkg/sms"
	"order_api_cart/pkg/utils"
	pkgValidator "order_api_cart/pkg/validator"
	"strconv"
)

type Service interface {
	CreateSession(context.Context, *Session) error
	VerifySession(context.Context, *Session) (string, error)
}

type ServiceAuth struct {
	repository *RepositoryAuth
	secret     string
}

func NewService(repository *RepositoryAuth, secret string) *ServiceAuth {
	return &ServiceAuth{repository: repository, secret: secret}
}

func (s *ServiceAuth) CreateSession(ctx context.Context, session *Session) error {
	pkgLog.InfoWithRequestID(ctx, "session request in service layer", logrus.Fields{})

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return fmt.Errorf("session creation failed: %w", err)
	}

	smsCode := utils.GetFakeCode()
	session.SessionID = sessionID
	session.SMSCode = smsCode

	if err = pkgValidator.ValidateStruct(session); err != nil {
		return fmt.Errorf("validation failed: %w", pkgErr.ErrValidation)
	}

	var model models.Session
	if err = utils.ConvertToModel(&model, session); err != nil {
		return fmt.Errorf("model conversion failed: %w", err)
	}

	if err = s.repository.CreateSession(ctx, &model); err != nil {
		return fmt.Errorf("database operation failed: %w", err)
	}

	if err = sms.SendFakeSMS(session.Phone, strconv.Itoa(session.SMSCode)); err != nil {
		_ = s.repository.DeleteSession(ctx, &model)
		return fmt.Errorf("notification service failed: %w", pkgErr.ErrServiceUnavailable)
	}

	return nil
}

func (s *ServiceAuth) VerifySession(ctx context.Context, session *Session) (string, error) {
	pkgLog.InfoWithRequestID(ctx, "validation request in service layer", logrus.Fields{})

	var requestedSession models.Session
	if err := s.repository.GetSession(ctx, &requestedSession, session.SessionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Join(pkgErr.ErrRecordNotFound, err)
		}
		return "", errors.Join(pkgErr.ErrQueryFailed, err)
	}

	if session.SMSCode != requestedSession.SMSCode {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrValidation.Error(), logrus.Fields{
			"error": pkgErr.ErrValidation.Error(),
		})
		return "", pkgErr.ErrValidation
	}

	if err := s.repository.DeleteSession(ctx, &requestedSession); err != nil {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrQueryFailed.Error(), logrus.Fields{
			"error": err.Error(),
		})
	}

	jwtString, err := pkgJWT.Create(s.secret, requestedSession.Phone)
	if err != nil {
		return "", pkgErr.ErrServiceUnavailable
	}

	return jwtString, nil
}
