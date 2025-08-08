package renderer

import (
	"context"

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
