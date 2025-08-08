package responses

import (
	"context"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AppError(ctx context.Context, requestID uuid.UUID, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return internalError(ctx, requestID)
	}

	withReq, derr := st.WithDetails(
		&errdetails.RequestInfo{RequestId: requestID.String()},
	)
	if derr != nil {
		return status.Errorf(
			codes.Internal,
			"failed to attach request info: %v",
			derr,
		)
	}

	return withReq.Err()
}

func internalError(
	ctx context.Context,
	requestID uuid.UUID,
) error {
	st := status.New(codes.Internal, "internal server error")

	info := &errdetails.ErrorInfo{
		Reason: "INTERNAL",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID.String(),
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

type Violation struct {
	Field       string
	Description string
}

func InvalidArgumentError(
	ctx context.Context,
	requestID uuid.UUID,
	violations ...Violation,
) error {
	st := status.New(codes.InvalidArgument, "bad request")

	info := &errdetails.ErrorInfo{
		Reason: "INVALID_ARGUMENT",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	var fb []*errdetails.BadRequest_FieldViolation
	for _, v := range violations {
		fb = append(fb, &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Description,
		})
	}
	br := &errdetails.BadRequest{FieldViolations: fb}

	ri := &errdetails.RequestInfo{
		RequestId: requestID.String(),
	}

	st, err := st.WithDetails(info, br, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}
