package response

// Response is the representation of http general response
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// ServiceError is error returned by the service(s)
type ServiceError struct {
	Code int
	Msg  string
	Err  error
}

func (s *ServiceError) Error() string {
	return s.Err.Error()
}
