package problems

import (
	"time"

	"github.com/chains-lab/cities-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func nowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(cause error) error {
	res, _ := status.New(codes.Internal, "internal server error").WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInternal.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInternal.Raise(cause, res)
}

var ErrorPermissionDenied = ape.Declare("PERMISSION_DENIED")

func RaisePermissionDenied(cause error, message string) error {
	res, _ := status.New(codes.PermissionDenied, message).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPermissionDenied.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInternal.Raise(cause, res)
}

var ErrorUnauthenticated = ape.Declare("UNAUTHENTICATED")

func RaiseUnauthenticated(cause error, message string) error {
	res, _ := status.New(codes.Unauthenticated, message).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUnauthenticated.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInternal.Raise(cause, res)
}

var ErrorInvalidArgument = ape.Declare("INVALID_ARGUMENT")

func RaiseInvalidArgument(cause error, message string, details ...*errdetails.BadRequest_FieldViolation) error {
	res, _ := status.New(codes.InvalidArgument, message).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidArgument.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.BadRequest{
			FieldViolations: details,
		},
	)

	return ErrorInternal.Raise(cause, res)
}

var ErrorNotFound = ape.Declare("NOT_FOUND")

func RaiseNotFound(cause error, message string) error {
	res, _ := status.New(codes.NotFound, message).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInternal.Raise(cause, res)
}
