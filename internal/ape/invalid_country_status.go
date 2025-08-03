package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonInvalidCountryStatus = "INVALID_COUNTRY_STATUS"

var ErrorInvalidCountryStatus = &apperr.ErrorObject{Reason: ReasonInvalidCountryStatus, Code: codes.InvalidArgument}

func RaiseInvalidCountryStatus(status string) error {
	return &apperr.ErrorObject{
		Code:    ErrorInvalidCountryStatus.Code,
		Reason:  ErrorInvalidCountryStatus.Reason,
		Message: "invalid country status",
		Cause:   fmt.Errorf("status '%s' is not one of the predefined values"), // No specific cause provided
		Details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "status",
				Description: fmt.Sprintf("status must be one of the predefined values: %s", enum.GetAllCountriesStatuses()),
				Reason:      ErrorInvalidCountryStatus.Reason,
			}}},
		},
	}
}
