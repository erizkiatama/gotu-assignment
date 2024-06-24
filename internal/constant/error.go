package constant

var ErrorInternalServer = "internal server error"

// User module error messages
var (
	ErrorUserAlreadyExists = "email already exists"
	ErrorCreateUserFailed  = "failed to create user"
	ErrorGenerateToken     = "failed to generate token"
	ErrorUserNotFound      = "user with that email does not found"
	ErrorGetUserFailed     = "failed to get user details"
	ErrorPasswordNotMatch  = "password does not match"
)
