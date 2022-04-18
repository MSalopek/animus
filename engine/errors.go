package engine

const (
	ErrNotFound      = "not found"
	ErrInternalError = "internal error"
	ErrUnauthorized  = "unauthorized"
	ErrForbidden     = "forbidden"

	ErrCouldNotLogin      = "could not log in"
	ErrCouldNotRegister   = "could not register"
	ErrInvalidCredentials = "invalid credentials"
	ErrNoAuthHeader       = "Authorization header not provided"
	ErrInvalidAuthToken   = "invalid Authorization token"

	ErrJWTExpired            = "JWT expired"
	ErrJWTClaimUnprocessable = "unprocessable JWT"

	ErrInvalidMeta       = "unprocessable meta field"
	ErrUserNotFound      = "user not found"
	ErrInvalidQueryParam = "invalid query param"
	ErrNotADirectory     = "not a directory"
)
