package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonInitiatorIsNotAdmin = "INITIATOR_IS_NOT_ADMIN"

var ErrorInitiatorIsNotAdmin = &apperr.ErrorObject{Reason: ReasonInitiatorIsNotAdmin, Code: codes.PermissionDenied}

func RaiseInitiatorIsNotAdmin(userID, cityID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorInitiatorIsNotAdmin.Code,
		Reason:  ErrorInitiatorIsNotAdmin.Reason,
		Message: "initiator is not an admin for the city",
		Cause:   fmt.Errorf("initiator with user_id:%s is not an admin for the city:%s", userID, cityID),
		Details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{
				Type:        ErrorInitiatorIsNotAdmin.Reason,
				Subject:     fmt.Sprintf("cities_admins?user_id=%s&city_id=%s", userID, cityID),
				Description: fmt.Sprintf("user with ID '%s' is not an admin for the city with ID '%s'", userID, cityID),
			}}},
		},
	}
}
