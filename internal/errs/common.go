package errs

import (
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(cause error) error {
	return ErrorInternal.Raise(
		cause,
		status.New(codes.Internal, "internal server error"),
	)
}
