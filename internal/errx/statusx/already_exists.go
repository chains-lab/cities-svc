package statusx

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AlreadyExistsCode    = codes.AlreadyExists
	AlreadyExistsReason  = "ALREADY_EXISTS"
	AlreadyExistsMessage = "The resource already exists"
)

func CountryAlreadyExists(countryName string) *status.Status {
	responses, _ := status.New(
		AlreadyExistsCode,
		AlreadyExistsMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCountry,
			ResourceName: fmt.Sprintf("country?name=%s", countryName),
			Description:  fmt.Sprintf("country with name '%s' already exists", countryName),
		},
		&errdetails.ErrorInfo{
			Reason: AlreadyExistsReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return responses
}

func CityOwnerAlreadyExists(cityID, userID uuid.UUID) *status.Status {
	responses, _ := status.New(
		AlreadyExistsCode,
		AlreadyExistsMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: "city_owner",
			ResourceName: fmt.Sprintf("city_owner?user_id=%s&city_id=%s", userID, cityID),
			Description:  fmt.Sprintf("city owner with userID '%s' already exists for cityID '%s'", userID, cityID),
		},
		&errdetails.ErrorInfo{
			Reason: AlreadyExistsReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return responses
}

func CityAdminAlreadyExists(cityID, userID uuid.UUID) *status.Status {
	responses, _ := status.New(
		AlreadyExistsCode,
		AlreadyExistsMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCityAdmin,
			ResourceName: fmt.Sprintf("city_admin?user_id=%s&city_id=%s", userID, cityID),
			Description:  fmt.Sprintf("city admin with userID '%s' already exists for cityID '%s'", userID, cityID),
		},
		&errdetails.ErrorInfo{
			Reason: AlreadyExistsReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return responses
}
