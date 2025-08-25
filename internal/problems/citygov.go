package problems

import (
	"github.com/chains-lab/cities-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInvalidCityGovRole = ape.Declare("INVALID_CITY_GOV_ROLE")

func RaiseInvalidCityGovRole(cause error, message string) error {
	st := status.New(codes.InvalidArgument, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCityGovRole.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInvalidCityGovRole.Raise(cause, st)
}

var ErrorCityGovNotFound = ape.Declare("CITY_GOV_NOT_FOUND")

func RaiseCityGovNotFound(cause error, message string) error {
	st := status.New(codes.NotFound, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityGovNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCityGovNotFound.Raise(cause, st)
}

var ErrorInitiatorIsNotCityGov = ape.Declare("INITIATOR_IS_NOT_CITY_GOV")

func RaiseInitiatorIsNotCityGov(cause error, message string) error {
	st := status.New(codes.PermissionDenied, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorIsNotCityGov.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInitiatorIsNotCityGov.Raise(cause, st)
}

var ErrorUserIsAlreadyCityGov = ape.Declare("USER_IS_ALREADY_CITY_GOV")

func RaiseUserIsAlreadyCityGov(cause error, message string) error {
	st := status.New(codes.AlreadyExists, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserIsAlreadyCityGov.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorUserIsAlreadyCityGov.Raise(cause, st)
}

var ErrorCannotDeleteCityAdmin = ape.Declare("CANNOT_DELETE_CITY_ADMIN")

func RaiseCannotDeleteCityAdmin(cause error, message string) error {
	st := status.New(codes.FailedPrecondition, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCannotDeleteCityAdmin.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCannotDeleteCityAdmin.Raise(cause, st)
}
