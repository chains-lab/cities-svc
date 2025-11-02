package errx

import "github.com/chains-lab/ape"

var ErrorCityStatusNotFound = ape.DeclareError("CITY_STATUS_NOT_FOUND")

var ErrorStatusAlreadyExists = ape.DeclareError("CITY_STATUS_ALREADY_EXISTS")

var ErrorCityStatusNotAllowedAdmin = ape.DeclareError("CITY_STATUS_NOT_ALLOWED_ADMIN")

var ErrorCityStatusNotAccessible = ape.DeclareError("CITY_STATUS_NOT_ACCESSIBLE")

var ErrorStatusInUse = ape.DeclareError("CITY_STATUS_IN_USE")
