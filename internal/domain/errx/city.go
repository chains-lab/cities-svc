package errx

import (
	"github.com/chains-lab/ape"
)

// TODO check in controller
var ErrorInvalidCityStatus = ape.DeclareError("INVALID_CITY_STATUS")

var ErrorInvalidSlug = ape.DeclareError("INVALID_SLUG")

var ErrorInvalidTimeZone = ape.DeclareError("INVALID_TIME_ZONE")

var ErrorInvalidPoint = ape.DeclareError("INVALID_POINT")

var ErrorInvalidCityName = ape.DeclareError("INVALID_CITY_NAME")

var ErrorCityNotFound = ape.DeclareError("CITY_NOT_FOUND")

var ErrorCityAlreadyExistsWithThisSlug = ape.DeclareError("CITY_ALREADY_EXISTS_WITH_THIS_SLUG")

// TODO: remove maybe should deleted an replecated ErrorCountryISNOtSupported
var ErrorCannotUpdateCityStatusInUnsupportedCountry = ape.DeclareError("CANNOT_UPDATE_CITY_STATUS_IN_UNSUPPORTED_COUNTRY")

var ErrorCityIsNotSupported = ape.DeclareError("CITY_IS_NOT_SUPPORTED")
