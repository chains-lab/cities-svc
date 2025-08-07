package responses

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidArgumentCode    = codes.InvalidArgument
	InvalidArgumentReason  = "INVALID_ARGUMENT"
	InvalidArgumentMessage = "invalid argument"
)

func InvalidCityStatus(cityStatus enum.CityStatus) *status.Status {
	response, _ := status.New(
		InvalidArgumentCode,
		InvalidArgumentMessage,
	).WithDetails(
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "city_status",
					Description: fmt.Sprintf("invalid city status %s, must be one of %s", cityStatus, enum.GetAllCitiesStatuses()),
				},
			},
		},
		&errdetails.ErrorInfo{
			Reason: InvalidArgumentReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func InvalidCountryStatus(countryStatus enum.CountryStatus) *status.Status {
	response, _ := status.New(
		InvalidArgumentCode,
		InvalidArgumentMessage,
	).WithDetails(
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "country_status",
					Description: fmt.Sprintf("invalid country status %s, must be one of %s", countryStatus, enum.GetAllCountriesStatuses()),
				},
			},
		},
		&errdetails.ErrorInfo{
			Reason: InvalidArgumentReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func InvalidCityAdminRole(cityAdminRole enum.CityAdminRole) *status.Status {
	response, _ := status.New(
		InvalidArgumentCode,
		InvalidArgumentMessage,
	).WithDetails(
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "city_admin_role",
					Description: fmt.Sprintf("invalid city admin role %s, must be one of %s", cityAdminRole, enum.GetAllCitiesAdminsRoles()),
				},
			},
		},
		&errdetails.ErrorInfo{
			Reason: InvalidArgumentReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
