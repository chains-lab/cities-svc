package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInvalidCityGovRole = ape.DeclareError("INVALID_CITY_GOV_ROLE")

var ErrorCityGovNotFound = ape.DeclareError("CITY_GOV_NOT_FOUND")

var ErrorInitiatorIsNotCityGov = ape.DeclareError("INITIATOR_IS_NOT_CITY_GOV")

var ErrorUserIsAlreadyCityGov = ape.DeclareError("USER_IS_ALREADY_CITY_GOV")

var ErrorCannotDeleteCityAdmin = ape.DeclareError("CANNOT_DELETE_CITY_ADMIN")
