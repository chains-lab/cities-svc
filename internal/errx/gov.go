package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInvalidGovRole = ape.DeclareError("INVALID_CITY_GOV_ROLE")

var ErrorCityGovNotFound = ape.DeclareError("CITY_GOV_NOT_FOUND")

var ErrorInitiatorIsNotThisCityGov = ape.DeclareError("INITIATOR_IS_NOT_IN_THIS_CITY_GOV")

var ErrorInitiatorRoleHaveNotEnoughRights = ape.DeclareError("INITIATOR_ROLE_HAVE_NOT_ENOUGH_RIGHTS")

var ErrorGovAlreadyExists = ape.DeclareError("CITY_GOV_ALREADY_EXISTS")

var ErrorCannotRefuseMayor = ape.DeclareError("CANNOT_REFUSE_MAYOR")

var ErrorInitiatorIsNotActiveCityGov = ape.DeclareError("NOT_ACTIVE_CITY_GOV_INITIATOR")
