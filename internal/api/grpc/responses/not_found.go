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
	NotFoundCode    = codes.NotFound
	NotFoundReason  = "NOT_FOUND"
	NotFoundMessage = "The requested resource was not found"
)

func CityNotFoundByID(cityID uuid.UUID) *status.Status {
	response, _ := status.New(
		NotFoundCode,
		NotFoundMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?id=%s", cityID),
			Owner:        fmt.Sprintf("cities_admins?role=%s&city_id=%s", enum.CityOwner, cityID),
			Description:  fmt.Sprintf("city with ID '%s' not found", cityID),
		},
		&errdetails.ErrorInfo{
			Reason: NotFoundReason,
			Domain: serviceName, //TODO thinking about it
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func CityNotFoundByName(name string) *status.Status {
	response, _ := status.New(
		NotFoundCode,
		NotFoundMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?name=%s", name),
			Description:  fmt.Sprintf("city with name '%s' not found", name),
		},
		&errdetails.ErrorInfo{
			Reason: NotFoundReason,
			Domain: serviceName, //TODO thinking about it
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func CountryNotFoundByID(countryID uuid.UUID) *status.Status {
	response, _ := status.New(
		NotFoundCode,
		NotFoundMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCountry,
			ResourceName: fmt.Sprintf("countries?id=%s", countryID),
			Description:  fmt.Sprintf("country with ID '%s' not found", countryID),
		},
		&errdetails.ErrorInfo{
			Reason: NotFoundReason,
			Domain: serviceName, //TODO thinking about it
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func CityAdminNotFound(cityID, userID uuid.UUID) *status.Status {
	response, _ := status.New(
		NotFoundCode,
		NotFoundMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: ResourceTypeCityAdmin,
			ResourceName: fmt.Sprintf("cities_admins?city_id=%s&user_id=%s", cityID, userID),
			Description:  fmt.Sprintf("city admin for city id: %s with user_id '%s' not found", cityID, userID),
		},
		&errdetails.ErrorInfo{
			Reason: NotFoundReason,
			Domain: serviceName, //TODO thinking about it
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
