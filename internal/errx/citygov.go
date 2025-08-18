package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
)

var ErrorInvalidCityGovRole = ape.Declare("INVALID_CITY_GOV_ROLE")

func RaiseInvalidCityGovRole(ctx context.Context, cause error, role string) error {
	st := status.New(codes.InvalidArgument, fmt.Sprintf("invalid city government role: %s", role))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCityGovRole.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInvalidCityGovRole.Raise(cause, st)
}

var ErrorCityGovNotFound = ape.Declare("CITY_GOV_NOT_FOUND")

func RaiseCityGovNotFound(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city government not found: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityGovNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityGovNotFound.Raise(cause, st)
}

func RaiseCityGovAdminNotFound(ctx context.Context, cause error, cityID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city government admin not found: city=%s", cityID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityGovNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityGovNotFound.Raise(cause, st)
}

var ErrorInitiatorIsNotCityGov = ape.Declare("INITIATOR_IS_NOT_CITY_GOV")

func RaiseInitiatorIsNotCityGov(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.PermissionDenied, fmt.Sprintf("initiator is not city government: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorIsNotCityGov.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInitiatorIsNotCityGov.Raise(cause, st)
}

var ErrorUserIsAlreadyCityGov = ape.Declare("USER_IS_ALREADY_CITY_GOV")

func RaiseUserIsAlreadyCityGov(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("user is already city government: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserIsAlreadyCityGov.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorUserIsAlreadyCityGov.Raise(cause, st)
}

var ErrorCannotDeleteCityAdmin = ape.Declare("CANNOT_DELETE_CITY_ADMIN")

func RaiseCannotDeleteCityAdmin(ctx context.Context, cause error, cityID, UserID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("cannot delete city admin: city=%s admin=%s", cityID, UserID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCannotDeleteCityAdmin.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCannotDeleteCityAdmin.Raise(cause, st)
}
