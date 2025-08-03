package ape

import (
	"github.com/chains-lab/apperr"
	"google.golang.org/grpc/codes"
)

const ReasonUnauthorized = "UNAUTHORIZED"

var ErrorUnauthorized = &apperr.ErrorObject{Reason: ReasonUnauthorized, Code: codes.Unauthenticated}

func RaiseIUnauthorized(cause error) error {
	return &apperr.ErrorObject{
		Code:    ErrorInternal.Code,
		Reason:  ErrorInternal.Reason,
		Message: "unauthorized access",
		Cause:   cause,
	}
}
