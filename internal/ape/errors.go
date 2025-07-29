package ape

import (
	"google.golang.org/grpc/codes"
)

const ServiceName = "cities-dir-svc"

var ErrorInternal = &Error{reason: ReasonInternal, code: codes.Internal}

func RaiseInternal(cause error) error {
	return &Error{
		code:    ErrorInternal.code,
		reason:  ErrorInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}
