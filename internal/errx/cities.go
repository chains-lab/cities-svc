package errx

import (
	"github.com/chains-lab/cities-dir-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCityNotFound = ape.Declare("CITY_NOT_FOUND")

func RaiseCityNotFoundByID(cause error, cityID uuid.UUID) error {
	return ErrorCityNotFound.Raise(
		cause,
		statusx.CityNotFoundByID(cityID),
	)
}

func RaiseCityNotFoundByName(cause error, cityName string) error {
	return ErrorCityNotFound.Raise(
		cause,
		statusx.CityNotFoundByName(cityName),
	)
}

var ErrorInvalidCityStatus = ape.Declare("INVALID_CITY_STATUS")

func RaiseInvalidCityStatus(cause error, status string) error {
	return ErrorInvalidCityStatus.Raise(
		cause,
		statusx.InvalidCityStatus(status),
	)
}

var ErrorCityStatusIsNotApplicable = ape.Declare("CITY_STATUS_IS_NOT_APPLICABLE")

func RaiseCityStatusIsNotApplicable(cause error, cityID uuid.UUID, expectedStatus, curStatus string) error {
	return ErrorCityStatusIsNotApplicable.Raise(
		cause,
		statusx.CityStatusNotApplicable(cityID, expectedStatus, curStatus),
	)
}
