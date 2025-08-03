package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonCityAdminNotFound = "CITY_ADMIN_NOT_FOUND"

var CityAdminNotFound = &apperr.ErrorObject{Reason: ReasonCityAdminNotFound, Code: codes.NotFound}

func RaiseCityAdminNotFound(err error, userID, cityID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    CityAdminNotFound.Code,
		Reason:  CityAdminNotFound.Reason,
		Message: "city admin not found",
		Cause:   err,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "city_admin",
				ResourceName: fmt.Sprintf("city_admin?user_id=%s&city_id=%s", userID, cityID),
				Owner:        fmt.Sprintf("user?id=%s", userID),
				Description:  fmt.Sprintf("city admin with user_id '%s' for city with ID '%s' does not exist", userID, cityID),
			},
		},
	}
}
