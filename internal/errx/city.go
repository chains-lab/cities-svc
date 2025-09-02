package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCityNotFound = ape.DeclareError("CITY_NOT_FOUND")

var ErrorCityAlreadyExists = ape.DeclareError("CITY_ALREADY_EXISTS")

var ErrorCityAlreadyExistsWithThisSlug = ape.DeclareError("CITY_ALREADY_EXISTS_WITH_THIS_SLUG")

var ErrorInvalidCityStatus = ape.DeclareError("INVALID_CITY_STATUS")

var ErrorInvalidSlug = ape.DeclareError("INVALID_SLUG")

var ErrorInvalidTimeZone = ape.DeclareError("INVALID_TIME_ZONE")

var ErrorInvalidPoint = ape.DeclareError("INVALID_POINT")

var ErrorInvalidCityName = ape.DeclareError("INVALID_CITY_NAME")
