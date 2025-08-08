package statusx

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/constant"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PermissionDeniedCode    = codes.PermissionDenied
	PermissionDeniedReason  = "PERMISSION_DENIED"
	PermissionDeniedMessage = "permission denied"
)

func CityAdminHaveNotEnoughRights(cityID, userID uuid.UUID, description string) *status.Status {
	response, _ := status.New(
		PermissionDeniedCode,
		PermissionDeniedMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?id=%s", cityID.String()),
			Owner:        fmt.Sprintf("cities_admins?role=%s&city_id=%s", enum.CityAdminRoleOwner, cityID.String()),
			Description:  fmt.Sprintf("user with ID '%s' does not have enough rights to perform this action on city with ID '%s': %s", userID.String(), cityID.String(), description),
		},
		&errdetails.ErrorInfo{
			Reason: PermissionDeniedReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func InitiatorIsNotCityAdmin(cityID, userID uuid.UUID, description string) *status.Status {
	response, _ := status.New(
		PermissionDeniedCode,
		PermissionDeniedMessage,
	).WithDetails(
		&errdetails.ResourceInfo{
			ResourceType: constant.ResourceTypeCity,
			ResourceName: fmt.Sprintf("cities?id=%s", cityID.String()),
			Owner:        fmt.Sprintf("cities_admins?user_id=%s&city_id=%s", userID.String(), cityID.String()),
			Description:  fmt.Sprintf("initiator with ID '%s' is not a city admin for city with ID '%s': %s", userID.String(), cityID.String(), description),
		},
		&errdetails.ErrorInfo{
			Reason: PermissionDeniedReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}

func PermissionDeniedByRole(userID uuid.UUID, role roles.Role) *status.Status {
	response, _ := status.New(
		PermissionDeniedCode,
		PermissionDeniedMessage,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: PermissionDeniedReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
