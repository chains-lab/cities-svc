package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonCountryNotFound = "COUNTRY_NOT_FOUND"

var ErrorCountryNotFound = &apperr.ErrorObject{Reason: ReasonCountryNotFound, Code: codes.NotFound}

func RaiseCountryNotFoundByID(cause error, countryID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorCountryNotFound.Code,
		Reason:  ErrorCountryNotFound.Reason,
		Message: "country not found",
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "country",
				ResourceName: fmt.Sprintf("country?id=%s", countryID),
				Description:  fmt.Sprintf("country with ID '%s' does not exist", countryID),
			},
		},
	}
}

func RaiseCountryNotFoundByName(cause error, name string) error {
	return &apperr.ErrorObject{
		Code:    ErrorCountryNotFound.Code,
		Reason:  ErrorCountryNotFound.Reason,
		Message: "country not found",
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "country",
				ResourceName: fmt.Sprintf("country?name=%s", name),
				Description:  fmt.Sprintf("country with name '%s' does not exist", name),
			},
		},
	}
}
