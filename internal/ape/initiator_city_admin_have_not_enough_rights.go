package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonInitiatorCityAdminHaveNotEnoughRights = "INITIATOR_CITY_ADMIN_HAVE_NOT_ENOUGH_RIGHTS"

var ErrorInitiatorCityAdminHaveNotEnoughRights = &apperr.ErrorObject{Reason: ReasonInitiatorCityAdminHaveNotEnoughRights, Code: codes.PermissionDenied}

func RaiseInitiatorCityAdminHaveNotEnoughRights(userID, cityID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorInitiatorCityAdminHaveNotEnoughRights.Code,
		Reason:  ErrorInitiatorCityAdminHaveNotEnoughRights.Reason,
		Message: "initiator does not have enough rights to perform this action",
		Cause:   fmt.Errorf("initiator does not have enough rights to perform this action"),
		Details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorInitiatorCityAdminHaveNotEnoughRights.Reason,
				Subject:     fmt.Sprintf("cities_admins?user_id=%s&city_id=%s", userID, cityID),
				Description: fmt.Sprintf("user with ID '%s' does not have enough rights to perform this action", userID),
			}}},
		},
	}
}
