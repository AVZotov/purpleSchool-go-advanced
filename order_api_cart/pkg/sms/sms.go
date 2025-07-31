package sms

import (
	"github.com/sirupsen/logrus"
	pkgLogger "order_api_cart/pkg/logger"
)

func SendFakeSMS(phone, code string) error {
	pkgLogger.Logger.WithFields(logrus.Fields{
		"phone": phone,
		"code":  code,
	}).Info("SendFakeSMS")
	return nil
}
