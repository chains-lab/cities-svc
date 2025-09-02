package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInvalidGovRole = ape.DeclareError("INVALID_CITY_GOV_ROLE")

var ErrorCityGovNotFound = ape.DeclareError("CITY_GOV_NOT_FOUND")

var ErrorInitiatorIsNotCityGov = ape.DeclareError("INITIATOR_IS_NOT_CITY_GOV")

var ErrorGovAlreadyExists = ape.DeclareError("CITY_GOV_ALREADY_EXISTS")

var ErrorCannotDeleteMayor = ape.DeclareError("CANNOT_DELETE_MAYOR")

var ErrorCannotRefuseMayor = ape.DeclareError("CANNOT_REFUSE_MAYOR")

var ErrorAdvisorMaxNumberReached = ape.DeclareError("ADVISOR_MAX_NUMBER_REACHED")

var ErrorCannotDeleteCityAdmin = ape.DeclareError("CANNOT_DELETE_CITY_ADMIN")

var ErrorNotActiveCityGovInitiator = ape.DeclareError("NOT_ACTIVE_CITY_GOV_INITIATOR")

var ErrorCannotUpdateInactiveGov = ape.DeclareError("CANNOT_UPDATE_INACTIVE_GOV")

var ErrorCannotUpdateSelfGov = ape.DeclareError("CANNOT_UPDATE_SELF_GOV")

var ErrorCannotUpdateMayorGovByOther = ape.DeclareError("CANNOT_UPDATE_MAYOR_GOV_BY_OTHER_GOV")

var ErrorInvalidGovStatus = ape.DeclareError("INVALID_GOV_STATUS")
