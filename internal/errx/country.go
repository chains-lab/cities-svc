package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCountryAlreadyExists = ape.DeclareError("COUNTRY_ALREADY_EXISTS")

var ErrorCountryNotFound = ape.DeclareError("COUNTRY_NOT_FOUND")

var ErrorInvalidCountryStatus = ape.DeclareError("INVALID_COUNTRY_STATUS")
