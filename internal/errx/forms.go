package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorFormNotFound = ape.Declare("FORM_NOT_FOUND")

func RaiseFormNotFoundByID(ctx context.Context, cause error, formId uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("form not found: %s", formId))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorFormNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)

	return ErrorFormNotFound.Raise(cause, st)
}
