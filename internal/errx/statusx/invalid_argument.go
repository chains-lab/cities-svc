package statusx

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidArgumentCode    = codes.InvalidArgument
	InvalidArgumentReason  = "INVALID_ARGUMENT"
	InvalidArgumentMessage = "invalid argument"
)

// InvalidCityStatus returns a status indicating that the provided city status is invalid.
// It must use when status is not expected
// example u have 3 statuses for city: "active", "suspended", "deleted
// and the user tries to set "archived" status, this function will return an error.
func InvalidCityStatus(cityStatus string) *status.Status {
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
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

// InvalidCountryStatus returns a status indicating that the provided country status is invalid.
// It must use when status is not expected
// example u have 3 statuses for country: "active", "suspended", "deleted"
// and the user tries to set "archived" status, this function will return an error.
func InvalidCountryStatus(countryStatus string) *status.Status {
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
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

// InvalidCityAdminRole returns a status indicating that the provided city admin role is invalid.
// It must use when role is not expected
// example u have 3 roles for city admin: "owner", "manager", "viewer
// and the user tries to set "editor" role, this function will return an error.
func InvalidCityAdminRole(cityAdminRole string) *status.Status {
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
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
