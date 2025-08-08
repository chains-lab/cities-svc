package errx

import (
	"github.com/chains-lab/cities-dir-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCountryAlreadyExists = ape.Declare("COUNTRY_ALREADY_EXISTS")

func RaiseCountryAlreadyExists(cause error, countryName string) error {
	return ErrorCountryAlreadyExists.Raise(
		cause,
		statusx.CountryAlreadyExists(countryName),
	)
}

var ErrorCountryNotFound = ape.Declare("COUNTRY_NOT_FOUND")

func RaiseCountryNotFoundByID(cause error, ID uuid.UUID) error {
	return ErrorCountryNotFound.Raise(
		cause,
		statusx.CountryNotFoundByID(ID),
	)
}

var ErrorInvalidCountryStatus = ape.Declare("INVALID_COUNTRY_STATUS")

func RaiseInvalidCountryStatus(cause error, status string) error {
	return ErrorInvalidCountryStatus.Raise(
		cause,
		statusx.InvalidCountryStatus(status),
	)
}

var ErrorCountryStatusIsNotApplicable = ape.Declare("COUNTRY_STATUS_IS_NOT_APPLICABLE")

func RaiseCountryStatusIsNotApplicable(cause error, countryID uuid.UUID, expectedStatus, curStatus string) error {
	return ErrorCountryStatusIsNotApplicable.Raise(
		cause,
		statusx.CountryStatusNotApplicable(countryID, curStatus, expectedStatus),
	)
}
