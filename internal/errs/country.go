package errs

import (
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCountryAlreadyExists = ape.Declare("COUNTRY_ALREADY_EXISTS")

func RaiseCountryAlreadyExists(cause error, countryName string) error {
	return ErrorCountryAlreadyExists.Raise(
		cause,
		responses.CountryAlreadyExists(countryName),
	)
}

var ErrorCountryNotFound = ape.Declare("COUNTRY_NOT_FOUND")

func RaiseCountryNotFoundByID(cause error, ID uuid.UUID) error {
	return ErrorCountryNotFound.Raise(
		cause,
		responses.CountryNotFoundByID(ID),
	)
}

var ErrorInvalidCountryStatus = ape.Declare("INVALID_COUNTRY_STATUS")

func RaiseInvalidCountryStatus(cause error, status enum.CountryStatus) error {
	return ErrorInvalidCountryStatus.Raise(
		cause,
		responses.InvalidCountryStatus(status),
	)
}

var ErrorCountryStatusIsNotApplicable = ape.Declare("COUNTRY_STATUS_IS_NOT_APPLICABLE")

func RaiseCountryStatusIsNotApplicable(cause error, countryID uuid.UUID, expectedStatus, curStatus enum.CountryStatus) error {
	return ErrorCountryStatusIsNotApplicable.Raise(
		cause,
		responses.CountryStatusNotApplicable(countryID, curStatus, expectedStatus),
	)
}
