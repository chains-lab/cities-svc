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

var ErrorCountryAlreadyExists = ape.Declare("COUNTRY_ALREADY_EXISTS")

func RaiseCountryAlreadyExists(ctx context.Context, cause error, countryName string) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("country %q already exists", countryName))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCountryAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCountryAlreadyExists.Raise(cause, st)
}

// --- COUNTRY_NOT_FOUND ---

var ErrorCountryNotFound = ape.Declare("COUNTRY_NOT_FOUND")

func RaiseCountryNotFoundByID(ctx context.Context, cause error, id uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("country not found: id=%s", id))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCountryNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCountryNotFound.Raise(cause, st)
}

// --- INVALID_COUNTRY_STATUS ---

var ErrorInvalidCountryStatus = ape.Declare("INVALID_COUNTRY_STATUS")

func RaiseInvalidCountryStatus(ctx context.Context, cause error, statusStr string) error {
	st := status.New(codes.InvalidArgument, fmt.Sprintf("invalid country status: %s", statusStr))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCountryStatus.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInvalidCountryStatus.Raise(cause, st)
}

// --- COUNTRY_STATUS_IS_NOT_APPLICABLE ---

var ErrorCountryStatusIsNotApplicable = ape.Declare("COUNTRY_STATUS_IS_NOT_APPLICABLE")

func RaiseCountryStatusIsNotApplicable(ctx context.Context, cause error, countryID uuid.UUID, expectedStatus, curStatus string) error {
	msg := fmt.Sprintf("status change is not applicable: country=%s expected=%s current=%s", countryID, expectedStatus, curStatus)
	st := status.New(codes.FailedPrecondition, msg)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCountryStatusIsNotApplicable.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCountryStatusIsNotApplicable.Raise(cause, st)
}

git add .
git commit -m 'smal fix'
git push