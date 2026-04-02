package constant

// Success messages
const (
	SuccessCreated   = "Resource created successfully"
	SuccessUpdated   = "Resource updated successfully"
	SuccessDeleted   = "Resource deleted successfully"
	SuccessFetched   = "Resource fetched successfully"
	SuccessOperation = "Operation completed successfully"
)

// Client error messages (4xx)
const (
	ErrBadRequest      = "Bad request"
	ErrUnauthorized    = "Unauthorized"
	ErrForbidden       = "Forbidden"
	ErrNotFound        = "Resource not found"
	ErrConflict        = "Resource conflict"
	ErrValidation      = "Validation failed"
	ErrTooManyRequests = "Too many requests"
)

// Server error messages (5xx)
const (
	ErrInternalServer     = "Internal server error"
	ErrServiceUnavailable = "Service unavailable"
	ErrGatewayTimeout     = "Gateway timeout"
)
