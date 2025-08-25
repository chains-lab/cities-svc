package problems

import (
	"github.com/chains-lab/cities-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorCountryAlreadyExists = ape.Declare("COUNTRY_ALREADY_EXISTS")

func RaiseCountryAlreadyExists(cause error, message string) error {
	st := status.New(codes.AlreadyExists, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCountryAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCountryAlreadyExists.Raise(cause, st)
}

var ErrorCountryNotFound = ape.Declare("COUNTRY_NOT_FOUND")

func RaiseCountryNotFound(cause error, message string) error {
	st := status.New(codes.NotFound, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCountryNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCountryNotFound.Raise(cause, st)
}

var ErrorInvalidCountryStatus = ape.Declare("INVALID_COUNTRY_STATUS")

func RaiseInvalidCountryStatus(cause error, message string) error {
	st := status.New(codes.InvalidArgument, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCountryStatus.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInvalidCountryStatus.Raise(cause, st)
}
