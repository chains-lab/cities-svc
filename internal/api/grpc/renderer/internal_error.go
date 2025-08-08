package renderer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func internalError(
	ctx context.Context,
	requestID uuid.UUID,
) error {
	st := status.New(codes.Internal, "internal server error")

	info := &errdetails.ErrorInfo{
		Reason: "INTERNAL",
		Domain: ServiceName,
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
