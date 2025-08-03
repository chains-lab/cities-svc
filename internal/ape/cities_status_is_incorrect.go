package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonCitiesStatusIsIncorrect = "CITIES_STATUS_IS_INCORRECT"

var CitiesStatusIsIncorrect = &apperr.ErrorObject{Reason: ReasonCitiesStatusIsIncorrect, Code: codes.InvalidArgument}

func RaiseCitiesStatusIsIncorrect(status string) error {
	return &apperr.ErrorObject{
		Code:    CitiesStatusIsIncorrect.Code,
		Reason:  CitiesStatusIsIncorrect.Reason,
		Message: "cities status is incorrect",
		Cause:   fmt.Errorf("status '%s' is not one of the predefined values", status),
		Details: []protoadapt.MessageV1{
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       "status",
				Description: fmt.Sprintf("status must be one of the predefined values: %s", enum.GetAllCitiesStatuses()),
				Reason:      CitiesStatusIsIncorrect.Reason,
			}}},
		},
	}
}
