package ape

import (
	"fmt"

	"github.com/chains-lab/apperr"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ReasonCityNotFound = "CITY_NOT_FOUND"

var ErrorCityNotFound = &apperr.ErrorObject{Reason: ReasonCityNotFound, Code: codes.NotFound}

func RaiseCityNotFoundByID(cause error, cityID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorCityNotFound.Code,
		Reason:  ErrorCityNotFound.Reason,
		Message: fmt.Sprintf("city with not found"),
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "city",
				ResourceName: fmt.Sprintf("city?id=:%s", cityID),
				Description:  fmt.Sprintf("city with ID '%s' does not exist", cityID),
			},
		},
	}
}

func RaiseCitiesNotFoundByNameAndCountryID(cause error, name string, countryID uuid.UUID) error {
	return &apperr.ErrorObject{
		Code:    ErrorCityNotFound.Code,
		Reason:  ErrorCityNotFound.Reason,
		Message: "city not found",
		Cause:   cause,
		Details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "city",
				ResourceName: fmt.Sprintf("city?name=%s&country_id=%s", name, countryID),
				Description:  fmt.Sprintf("city with name '%s' and country ID '%s' does not exist", name, countryID),
			},
		},
	}
}
