package http

// APIErr is the standard return of error.
type APIErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
