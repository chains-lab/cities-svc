package errx

import "github.com/chains-lab/ape"

var ErrorCityNotFound = ape.DeclareError("CITY_NOT_FOUND")

var ErrorInvalidSlug = ape.DeclareError("INVALID_SLUG")

var ErrorCityAlreadyExistsWithThisSlug = ape.DeclareError("CITY_ALREADY_EXISTS_WITH_THIS_SLUG")

var ErrorInvalidCityName = ape.DeclareError("INVALID_CITY_NAME")
