package mock

import "errors"

var ErrMock = errors.New("mock error")

type HTTPServer struct {
	ShouldError bool
}

func (s *HTTPServer) ListenAndServe() error {
	if s.ShouldError {
		return ErrMock
	}
	return nil
}

func (s *HTTPServer) ListenAndServeTLS(certFile, keyFile string) error {
	if s.ShouldError {
		return ErrMock
	}
	return nil
}
