package errx

import "github.com/chains-lab/ape"

var ErrorInvalidCityAdminRole = ape.DeclareError("INVALID_CITY_ADMIN_ROLE")

var ErrorCityAdminNotFound = ape.DeclareError("CITY_ADMIN_NOT_FOUND")

var ErrorInitiatorHasNoRights = ape.DeclareError("INITIATOR_HAS_NO_RIGHTS")

var ErrorInitiatorIsNotCityAdmin = ape.DeclareError("INITIATOR_IS_NOT_CITY_ADMIN")

var ErrorCityAdminAlreadyExists = ape.DeclareError("CITY_ADMIN_ALREADY_EXISTS")

var ErrorCityAdminTechLeadCannotRefuseOwn = ape.DeclareError("CITY_ADMIN_TECH_LEAD_CANNOT_REFUSE_OWN")
