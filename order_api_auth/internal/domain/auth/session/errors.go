package session

import "errors"

var (
	ErrDBInternalError     = errors.New("error database internal error")
	ErrCreatingSession     = errors.New("error creating session")
	ErrGettingSession      = errors.New("error getting session")
	ErrDeletingSession     = errors.New("error deleting session")
	ErrSessionNotFound     = errors.New("error session not found")
	ErrInvalidSMSCode      = errors.New("error invalid SMS code")
	ErrInternalError       = errors.New("error internal server error")
	ErrGeneratingToken     = errors.New("error generating token")
	ErrSendingSMS          = errors.New("error sending sms")
	ErrGeneratingSessionID = errors.New("error generating sessionID")
)
