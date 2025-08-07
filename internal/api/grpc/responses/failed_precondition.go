package responses

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/enum"
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
			ResourceType: ResourceTypeCityAdmin,
			ResourceName: fmt.Sprintf("city_admin?city_id=%s&user_id=%s", cityID, ownerID),
			Owner:        ownerID.String(),
			Description:  fmt.Sprintf("cannot delete city owner %s for city %s", ownerID, cityID),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": fmt.Sprintf("%s", time.Now().UTC().Format(time.RFC3339Nano)),
			},
		})

	return response
}

func CityStatusNotApplicable(cityID uuid.UUID, expectedStatus, curStatus enum.CityStatus) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		FailedPreconditionMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?id=%s", cityID),
			Owner:        fmt.Sprintf("cities_admins?role=%s&city_id=%s", enum.CityOwner, cityID),
			Description:  fmt.Sprintf("city with ID '%s' is not %s, current status: %s", cityID, expectedStatus, curStatus),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func CountryStatusNotApplicable(countryID uuid.UUID, expectedStatus, curStatus enum.CountryStatus) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		FailedPreconditionMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCountry,
			ResourceName: fmt.Sprintf("countries?id=%s", countryID),
			Description:  fmt.Sprintf("country with ID '%s' is not %s, current status: %s", countryID, expectedStatus, curStatus),
		},
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: serviceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
