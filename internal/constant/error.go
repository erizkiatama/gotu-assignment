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

// Book module error messages
var (
	ErrorListBooksFailed = "failed to list all books"
)

// Order module error messages
var (
	ErrorCreateOrderFailed       = "failed to create order"
	ErrorCreateOrderDetailFailed = "failed to create order detail"
	ErrorGetAllOrderFailed       = "failed to get order list"
	ErrorGetOrderDetailFailed    = "failed to get order detail"
	ErrorOrderNotFound           = "order not found"
)
