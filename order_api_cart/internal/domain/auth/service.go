package auth

import (
	"context"
	"errors"
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
	repository RepositoryAuth
	secret     string
}

func NewService(repository RepositoryAuth, secret string) *ServiceAuth {
	return &ServiceAuth{repository: repository, secret: secret}
}

func (s *ServiceAuth) CreateSession(ctx context.Context, session *Session) error {
	pkgLog.InfoWithRequestID(ctx, "session request in service layer", logrus.Fields{})

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrGeneratingSessionID.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return pkgErr.ErrGeneratingSessionID
	}

	smsCode := utils.GetFakeCode()
	session.SessionID = sessionID
	session.SMSCode = smsCode

	if err = pkgValidator.ValidateStruct(session); err != nil {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrInvalidInput.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return pkgErr.ErrInvalidInput
	}

	var model models.Session
	if err = utils.ConvertToModel(model, session); err != nil {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrConvertingToModel.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return pkgErr.ErrConvertingToModel
	}

	if err = s.repository.CreateSession(ctx, &model); err != nil {
		return errors.Join(pkgErr.ErrTransactionFailed, err)
	}

	smsErr := sms.SendFakeSMS(session.Phone, strconv.Itoa(session.SMSCode))
	if smsErr != nil {
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrSendingSMS.Error(), logrus.Fields{
			"error": smsErr.Error(),
		})
		deleteErr := s.repository.DeleteSession(ctx, &model)
		if deleteErr != nil {
			pkgLog.ErrorWithRequestID(ctx, deleteErr.Error(), logrus.Fields{
				"error": deleteErr.Error(),
			})
		}
		return errors.Join(pkgErr.ErrSendingSMS, smsErr, deleteErr)
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
		pkgLog.ErrorWithRequestID(ctx, pkgErr.ErrCreatingToken.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return "", pkgErr.ErrCreatingToken
	}

	return jwtString, nil
}
