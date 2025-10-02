package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCountryNotFound = ape.DeclareError("COUNTRY_NOT_FOUND")

var ErrorInvalidCountryStatus = ape.DeclareError("INVALID_COUNTRY_STATUS")

var ErrorCountryAlreadyExistsWithThisName = ape.DeclareError("COUNTRY_ALREADY_EXISTS_WITH_THIS_NAME")

var ErrorCountryIsNotSupported = ape.DeclareError("COUNTRY_NOT_SUPPORTED")
