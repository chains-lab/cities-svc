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

var ErrorCityAdminNotFound = ape.Declare("CITY_ADMIN_NOT_FOUND")

func RaiseCityAdminNotFound(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city admin not found: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityAdminNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityAdminNotFound.Raise(cause, st)
}

// --- CITY_ADMIN_HAVE_NOT_ENOUGH_RIGHTS ---

var ErrorCityAdminHaveNotEnoughRights = ape.Declare("CITY_ADMIN_HAVE_NOT_ENOUGH_RIGHTS")

func RaiseCityAdminHaveNotEnoughRights(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.PermissionDenied, fmt.Sprintf("city admin have not enough rights: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityAdminHaveNotEnoughRights.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityAdminHaveNotEnoughRights.Raise(cause, st)
}

// --- INITIATOR_IS_NOT_CITY_ADMIN ---

var ErrorInitiatorIsNotCityAdmin = ape.Declare("INITIATOR_IS_NOT_CITY_ADMIN")

func RaiseInitiatorIsNotCityAdmin(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.PermissionDenied, fmt.Sprintf("initiator is not city admin: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorIsNotCityAdmin.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInitiatorIsNotCityAdmin.Raise(cause, st)
}

// --- CITY_OWNER_ALREADY_EXISTS ---

var ErrorCityOwnerAlreadyExits = ape.Declare("CITY_OWNER_ALREADY_EXISTS")

func RaiseCityOwnerAlreadyExits(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("city owner already exists: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityOwnerAlreadyExits.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityOwnerAlreadyExits.Raise(cause, st)
}

// --- USER_IS_ALREADY_CITY_ADMIN ---

var ErrorUserIsAlreadyCityAdmin = ape.Declare("USER_IS_ALREADY_CITY_ADMIN")

func RaiseUserIsAlreadyCityAdmin(ctx context.Context, cause error, cityID, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("user is already city admin: city=%s user=%s", cityID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserIsAlreadyCityAdmin.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorUserIsAlreadyCityAdmin.Raise(cause, st)
}

// --- CITY_ADMIN_ROLE_IS_INVALID ---

var ErrorInvalidCityAdminRole = ape.Declare("CITY_ADMIN_ROLE_IS_INVALID")

func RaiseInvalidCityAdminRole(ctx context.Context, cause error, role string) error {
	st := status.New(codes.InvalidArgument, fmt.Sprintf("invalid city admin role: %s", role))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCityAdminRole.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInvalidCityAdminRole.Raise(cause, st)
}

// --- CANNOT_DELETE_CITY_OWNER ---

var ErrorCannotDeleteCityOwner = ape.Declare("CANNOT_DELETE_CITY_OWNER")

func RaiseCannotDeleteCityOwner(ctx context.Context, cause error, cityID, ownerID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("cannot delete city owner: city=%s owner=%s", cityID, ownerID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCannotDeleteCityOwner.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCannotDeleteCityOwner.Raise(cause, st)
}
