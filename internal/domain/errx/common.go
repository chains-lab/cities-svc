package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorInvalidTimeZone = ape.DeclareError("INVALID_TIME_ZONE")

var ErrorInvalidPoint = ape.DeclareError("INVALID_POINT")

var ErrorInvalidCountryISO3ID = ape.DeclareError("INVALID_COUNTRY_ISO3_ID")
