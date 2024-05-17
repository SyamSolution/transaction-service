package util

const (
	SUCCESS_RESPONSE_CODE = 3001
	SUCCESS_RESPONSE_MSG  = "success"

	ERROR_BASE_CODE          = 4001
	ERROR_BASE_MSG           = "something is wrong, report to support team"
	ERROR_NOT_FOUND_CODE     = 4002
	ERROR_NOT_FOUND_MSG      = "attribute not found"
	ERROR_INVALID_PARAM_CODE = 4101
	ERROR_INVALID_PARAM_MSG  = "failed"
	ERROR_UNAUTHORIZE_CODE   = 4102
	ERROR_UNAUTHORIZE_MSG    = "unauthorize access"
	ERROR_NOTACCEPTABLE_CODE = 4103
	ERROR_NOTACCEPTABLE_MSG  = "not accetable value"
	ERROR_DELETED_POST_MSG   = "Sorry, the post is deleted. Explore other interesting content!"
)

const DEFAULT_BUSINESS_ERROR_CODE = ERROR_BASE_CODE
const DEFAULT_BUSINESS_ERROR_MESSAGE = ERROR_BASE_MSG

const DEFAULT_LIMIT_PAGINATION = 10

// date & time
const (
	TIMESTAMP_DEFAULT_FORMAT = "2006-01-02T15:04:05-0700"
	TIMESTAMP_DATE_FORMAT    = "2006-01-02 15:04:05"
)

