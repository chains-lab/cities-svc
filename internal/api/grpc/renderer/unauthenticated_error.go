package renderer

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnauthenticatedError(
	message string, //for logging
	requestID *uuid.UUID,
) error {
	st := status.New(
		codes.Unauthenticated,
		"bad credentials",
	)

	info := &errdetails.ErrorInfo{
		Reason: "UNAUTHENTICATED",
		Domain: ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	var ri *errdetails.RequestInfo
	if requestID != nil {
		ri = &errdetails.RequestInfo{
			RequestId: requestID.String(),
		}
	}

	st, err := st.WithDetails(info, ri)

	if err != nil {
		return st.Err()
	}

	return st.Err()
}
