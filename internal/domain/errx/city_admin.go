package errx

import "github.com/chains-lab/ape"

var ErrorCityAdminNotFound = ape.DeclareError("CITY_ADMIN_NOT_FOUND")

var ErrorCityAdminAlreadyExists = ape.DeclareError("CITY_ADMIN_ALREADY_EXISTS")

var ErrorInvalidCityAdminRole = ape.DeclareError("INVALID_CITY_ADMIN_ROLE")

var ErrorCityAdminTechLeadCannotRefuseOwn = ape.DeclareError("CITY_ADMIN_TECH_LEAD_CANNOT_REFUSE_OWN")

var ErrorCannotDeleteYourself = ape.DeclareError("CANNOT_DELETE_YOURSELF")
