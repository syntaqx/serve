package mock

import "errors"

// ErrMock is a mock error
var ErrMock = errors.New("mock error")

// HTTPServer is a mock http server
type HTTPServer struct {
	ShouldError bool
}

// ListenAndServe is a mock http server method
func (s *HTTPServer) ListenAndServe() error {
	if s.ShouldError {
		return ErrMock
	}
	return nil
}

// ListenAndServeTLS is a mock http server method
func (s *HTTPServer) ListenAndServeTLS(certFile, keyFile string) error {
	if s.ShouldError {
		return ErrMock
	}
	return nil
}
