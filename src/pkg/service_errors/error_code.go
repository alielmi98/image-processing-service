package service_errors

const (
	// Token
	UnExpectedError     = "Expected error"
	ClaimsNotFound      = "Claims not found"
	TokenRequired       = "token required"
	TokenExpired        = "token expired"
	TokenInvalid        = "token invalid"
	InvalidRefreshToken = "invalid refresh token"
	InvalidRolesFormat  = "invalid roles format"
	// User
	EmailExists               = "Email exists"
	UsernameExists            = "Username exists"
	PermissionDenied          = "Permission denied"
	UsernameOrPasswordInvalid = "username or password invalid"
	// Validation
	ValidationError = "validation error"
	UserIdNotFound  = "failed to get user ID from context"
	UserNotOwner    = "user is not the owner of this workout"
	InvalidStatus   = "invalid status. Status must be 'active' or 'completed' or 'canceled'"

	// DB
	RecordNotFound = "record not found"
	UnknownError   = "unknown error"
)
