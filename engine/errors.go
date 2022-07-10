package engine

type APIError struct {
	ErrString string
}

func (e *APIError) Error() string { return e.ErrString }

var (
	ErrNotFound      = &APIError{"not found"}
	ErrInternalError = &APIError{"internal error"}
	ErrUnauthorized  = &APIError{"unauthorized"}
	ErrForbidden     = &APIError{"forbidden"}

	ErrCouldNotLogin      = &APIError{"could not log in"}
	ErrCouldNotRegister   = &APIError{"could not register"}
	ErrInvalidCredentials = &APIError{"invalid credentials"}
	ErrNoAuthHeader       = &APIError{"Authorization header not provided"}
	ErrInvalidAuthToken   = &APIError{"invalid Authorization token"}

	ErrInvalidClientAuth      = &APIError{"invalid client authorization"}
	ErrInvalidClientSignature = &APIError{"invalid client signature"}

	ErrJWTExpired            = &APIError{"JWT expired"}
	ErrJWTClaimUnprocessable = &APIError{"unprocessable JWT"}

	ErrInvalidMeta        = &APIError{"unprocessable meta field"}
	ErrUserNotFound       = &APIError{"user not found"}
	ErrInvalidQueryParam  = &APIError{"invalid query param"}
	ErrInvalidRequestBody = &APIError{"invalid request body"}
	ErrNotADirectory      = &APIError{"not a directory"}

	ErrUnprocessableFormFile      = &APIError{"unprocessable form file"}
	ErrUnprocessableMultipartForm = &APIError{"unprocessable multipart form"}
	ErrMissingFormDirName         = &APIError{"missing directory name"}

	// failure to upload to s3 bucket or store file info in DB
	ErrFileSaveFailed = &APIError{"failed saving file"}
	ErrDirSaveFailed  = &APIError{"directory saving failed"}

	ErrMaxKeyAllocReached = &APIError{"maximum key allocation reached"}

	ErrNoCID = &APIError{"CID is missing"}
)
