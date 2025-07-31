package context

type contextKey string

const (
	CtxRequestId contextKey = "request_id"
	CtxUserPhone contextKey = "phone"
)
