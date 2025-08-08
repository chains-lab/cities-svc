package statusx

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	FailedPreconditionCode    = codes.FailedPrecondition
	FailedPreconditionReason  = "FAILED_PRECONDITION"
	FailedPreconditionMessage = "failed precondition: the operation cannot be performed in the current state of the resource"
)

func CannotDeleteCityOwner(cityID, ownerID uuid.UUID) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		FailedPreconditionMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCityAdmin,
			ResourceName: fmt.Sprintf("city_admin?city_id=%s&user_id=%s", cityID, ownerID),
			Owner:        ownerID.String(),
			Description:  fmt.Sprintf("cannot delete city owner %s for city %s", ownerID, cityID),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": fmt.Sprintf("%s", time.Now().UTC().Format(time.RFC3339Nano)),
			},
		})

	return response
}

// CityStatusNotApplicable returns a status indicating that the city is not in the expected status.
// This function must use when the city is not in the expected status and the operation cannot be
// performed. For example, when trying to update a city that status is "suspended"
// but the expected status for this operation is "active".
func CityStatusNotApplicable(cityID uuid.UUID, expectedStatus, curStatus string) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		FailedPreconditionMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?id=%s", cityID),
			Owner:        fmt.Sprintf("cities_admins?role=%s&city_id=%s", enum.CityAdminRoleOwner, cityID),
			Description:  fmt.Sprintf("city with ID '%s' is not %s, current status: %s", cityID, expectedStatus, curStatus),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

// CountryStatusNotApplicable returns a status indicating that the country is not in the expected status.
// This function must use when the country is not in the expected status and the operation cannot be performed.
// for example when trying to update a country that status is "suspended"
// but the expected status for this operation is "active".
func CountryStatusNotApplicable(countryID uuid.UUID, expectedStatus, curStatus string) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		FailedPreconditionMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCountry,
			ResourceName: fmt.Sprintf("countries?id=%s", countryID),
			Description:  fmt.Sprintf("country with ID '%s' is not %s, current status: %s", countryID, expectedStatus, curStatus),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
