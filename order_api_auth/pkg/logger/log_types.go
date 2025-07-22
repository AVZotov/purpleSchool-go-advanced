package logger

const (
	HttpRequestStart = "http_request_start"
	HttpRequestEnd   = "http_request_end"
	HttpRequestError = "http_request_error"
	DBMigrationStart = "db_migration"
	DBMigrationEnd   = "db_migration_end"
	DBError          = "db_error"
	DBWarn           = "db_warning"
	DBSuccess        = "db_success"
	ValidationError  = "validation_error"
	JsonError        = "json_error"
	RepositoryError  = "repository_error"
	HandlerStart     = "handler_start"
	HandlerEnd       = "handler_end"
	HandlerError     = "handler_error"
)
