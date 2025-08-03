package ape

import (
	"github.com/chains-lab/apperr"
	"google.golang.org/grpc/codes"
)

const ReasonInternal = "INTERNAL_ERROR"

var ErrorInternal = &apperr.ErrorObject{Reason: ReasonInternal, Code: codes.Internal}

func RaiseInternal(cause error) error {
	return &apperr.ErrorObject{
		Code:    ErrorInternal.Code,
		Reason:  ErrorInternal.Reason,
		Message: "unexpected internal error occurred",
		Cause:   cause,
	}
}
