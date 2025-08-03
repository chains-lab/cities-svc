package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonInvalidCityAdminRole = "INVALID_CITY_ADMIN_ROLE"

var ErrorInvalidCityAdminRole = &apperr.ErrorObject{Reason: ReasonInvalidCityAdminRole, Code: codes.InvalidArgument}

func RaiseInvalidCityAdminRole(role string) error {
	return &apperr.ErrorObject{
		Code:    ErrorInvalidCityAdminRole.Code,
		Reason:  ErrorInvalidCityAdminRole.Reason,
		Message: "invalid city admin role",
		Cause:   fmt.Errorf("role '%s' is not one of the predefined values", role),
		Details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "role",
				Description: fmt.Sprintf("role must be one of the predefined values: %s", enum.GetAllCitiesAdminsRoles()),
				Reason:      ErrorInvalidCityAdminRole.Reason,
			}}},
		},
	}
}
