package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInvalidCityAdminRole = ape.DeclareError("INVALID_CITY_GOV_ROLE")

var ErrorCityAdminNotFound = ape.DeclareError("CITY_MOD_NOT_FOUND")

var ErrorInitiatorIsNotCityAdmin = ape.DeclareError("INITIATOR_IS_NOT_CITY_MODER")

var ErrorInitiatorIsNotThisCityAdmin = ape.DeclareError("INITIATOR_IS_NOT_THIS_CITY_MODER")

var ErrorUserIsAlreadyCityAdmin = ape.DeclareError("USER_IS_ALREADY_CITY_MODER")
