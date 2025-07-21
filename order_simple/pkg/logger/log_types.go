package logger

const (
	HttpRequestStart    = "http_request_start"
	HttpRequestEnd      = "http_request_end"
	DBMigration         = "db_migration"
	DBError             = "db_error"
	DBWarn              = "db_warning"
	DBSuccess           = "db_success"
	ValidationError     = "validation_error"
	JSONError           = "json_error"
	RepositoryError     = "repository_error"
	HandlerRequestStart = "handler_request_start"
	HandlerRequestEnd   = "handler_request_end"
	HandlerError        = "handler_error"
)
