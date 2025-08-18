package errx

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(ctx context.Context, cause error) error {
	st := status.New(codes.Internal, "internal server error")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInternal.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInternal.Raise(cause, st)
}

var ErrorNoPermissions = ape.Declare("NO_PERMISSIONS")

func RaiseNoPermissions(ctx context.Context, cause error) error {
	res, _ := status.New(codes.PermissionDenied, cause.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorNoPermissions.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorNoPermissions.Raise(
		cause,
		res,
	)
}

var ErrorUnauthenticated = ape.Declare("UNAUTHENTICATED")

func RaiseUnauthenticated(ctx context.Context, cause error) error {
	res, _ := status.New(codes.Unauthenticated, cause.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUnauthenticated.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorUnauthenticated.Raise(
		cause,
		res,
	)
}
