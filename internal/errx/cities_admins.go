package errx

import (
	"github.com/chains-lab/cities-dir-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCityAdminNotFound = ape.Declare("CITY_ADMIN_NOT_FOUND")

func RaiseCityAdminNotFound(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityAdminNotFound.Raise(
		cause,
		statusx.CityAdminNotFound(cityID, userID),
	)
}

var ErrorCityAdminHaveNotEnoughRights = ape.Declare("CITY_ADMIN_HAVE_NOT_ENOUGH_RIGHTS")

func RaiseCityAdminHaveNotEnoughRights(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityAdminHaveNotEnoughRights.Raise(
		cause,
		statusx.CityAdminHaveNotEnoughRights(cityID, userID, cause.Error()),
	)
}

var ErrorInitiatorIsNotCityAdmin = ape.Declare("INITIATOR_IS_NOT_CITY_ADMIN")

func RaiseInitiatorIsNotCityAdmin(cause error, cityID, userID uuid.UUID) error {
	return ErrorInitiatorIsNotCityAdmin.Raise(
		cause,
		statusx.InitiatorIsNotCityAdmin(cityID, userID, cause.Error()),
	)
}

var ErrorCityOwnerAlreadyExits = ape.Declare("CITY_OWNER_ALREADY_EXISTS")

func RaiseCityOwnerAlreadyExits(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityOwnerAlreadyExits.Raise(
		cause,
		statusx.CityOwnerAlreadyExists(cityID, userID),
	)
}

var ErrorUserIsAlreadyCityAdmin = ape.Declare("USER_IS_ALREADY_CITY_ADMIN")

func RaiseUserIsAlreadyCityAdmin(cause error, cityID, userID uuid.UUID) error {
	return ErrorUserIsAlreadyCityAdmin.Raise(
		cause,
		statusx.CityAdminAlreadyExists(cityID, userID),
	)
}

var ErrorInvalidCityAdminRole = ape.Declare("CITY_ADMIN_ROLE_IS_INVALID")

func RaiseInvalidCityAdminRole(cause error, role string) error {
	return ErrorInvalidCityAdminRole.Raise(
		cause,
		statusx.InvalidCityAdminRole(role),
	)
}

var ErrorCannotDeleteCityOwner = ape.Declare("CANNOT_DELETE_CITY_OWNER")

func RaiseCannotDeleteCityOwner(cause error, cityID, OwnerID uuid.UUID) error {
	return ErrorCannotDeleteCityOwner.Raise(
		cause,
		statusx.CannotDeleteCityOwner(cityID, OwnerID),
	)
}
