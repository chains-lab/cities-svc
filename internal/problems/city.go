package problems

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorCityNotFound = ape.Declare("CITY_NOT_FOUND")

func RaiseCityNotFoundByID(ctx context.Context, cause error, cityID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city %s not found", cityID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorCityNotFound.Raise(cause, st)
}

func RaiseCityNotFoundByName(ctx context.Context, cause error, cityName string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city with name %q not found", cityName))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorCityNotFound.Raise(cause, st)
}

var ErrorInvalidCityStatus = ape.Declare("INVALID_CITY_STATUS")

func RaiseInvalidCityStatus(ctx context.Context, cause error, cityStatus string) error {
	st := status.New(codes.InvalidArgument, fmt.Sprintf("invalid city status: %s", cityStatus))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidCityStatus.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorInvalidCityStatus.Raise(cause, st)
}

var ErrorCityDetailsNotFound = ape.Declare("CITY_DETAILS_NOT_FOUND")

func RaiseCityDetailsNotFound(ctx context.Context, cause error, cityID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("city details not found: city=%s", cityID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCityDetailsNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorCityDetailsNotFound.Raise(cause, st)
}
