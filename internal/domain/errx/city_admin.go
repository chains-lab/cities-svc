package errx

import "github.com/chains-lab/ape"

var ErrorInvalidCityAdminRole = ape.DeclareError("INVALID_CITY_ADMIN_ROLE")

var ErrorCityAdminNotFound = ape.DeclareError("CITY_ADMIN_NOT_FOUND")

var ErrorInitiatorIsNotCityAdmin = ape.DeclareError("INITIATOR_IS_NOT_CITY_ADMIN")
