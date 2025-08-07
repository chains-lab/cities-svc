package errs

import (
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCityNotFound = ape.Declare("CITY_NOT_FOUND")

func RaiseCityNotFoundByID(cause error, cityID uuid.UUID) error {
	return ErrorCityNotFound.Raise(
		cause,
		responses.CityNotFoundByID(cityID),
	)
}

func RaiseCityNotFoundByName(cause error, cityName string) error {
	return ErrorCityNotFound.Raise(
		cause,
		responses.CityNotFoundByName(cityName),
	)
}

var ErrorInvalidCityStatus = ape.Declare("INVALID_CITY_STATUS")

func RaiseInvalidCityStatus(cause error, status enum.CityStatus) error {
	return ErrorInvalidCityStatus.Raise(
		cause,
		responses.InvalidCityStatus(status),
	)
}

var ErrorCityStatusIsNotApplicable = ape.Declare("CITY_STATUS_IS_NOT_APPLICABLE")

func RaiseCityStatusIsNotApplicable(cause error, cityID uuid.UUID, expectedStatus, curStatus enum.CityStatus) error {
	return ErrorCityStatusIsNotApplicable.Raise(
		cause,
		responses.CityStatusNotApplicable(cityID, expectedStatus, curStatus),
	)
}
