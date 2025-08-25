package problems

import (
	"github.com/chains-lab/cities-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorCityNotFound = ape.Declare("CITY_NOT_FOUND")

func RaiseCityNotFoundByID(cause error, message string) error {
	st := status.New(codes.NotFound, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)
	return ErrorCityNotFound.Raise(cause, st)
}

func RaiseCityNotFoundByName(cause error, message string) error {
	st := status.New(codes.NotFound, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)
	return ErrorCityNotFound.Raise(cause, st)
}

var ErrorInvalidCityStatus = ape.Declare("INVALID_CITY_STATUS")

func RaiseInvalidCityStatus(cause error, message string) error {
	st := status.New(codes.InvalidArgument, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCityStatus.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)
	return ErrorInvalidCityStatus.Raise(cause, st)
}

var ErrorCityDetailsNotFound = ape.Declare("CITY_DETAILS_NOT_FOUND")

func RaiseCityDetailsNotFound(cause error, message string) error {
	st := status.New(codes.NotFound, message)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityDetailsNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)
	return ErrorCityDetailsNotFound.Raise(cause, st)
}
