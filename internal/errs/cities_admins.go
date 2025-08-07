package errs

import (
	"fmt"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorCityAdminNotFound = ape.Declare("CITY_ADMIN_NOT_FOUND")

func RaiseCityAdminNotFound(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityAdminNotFound.Raise(
		cause,
		responses.CityAdminNotFound(cityID, userID),
	)
}

var ErrorCityAdminHaveNotEnoughRights = ape.Declare("CITY_ADMIN_HAVE_NOT_ENOUGH_RIGHTS")

func RaiseCityAdminHaveNotEnoughRights(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityAdminHaveNotEnoughRights.Raise(
		cause,
		responses.CityAdminHaveNotEnoughRights(cityID, userID, cause.Error()),
	)
}

var ErrorInitiatorIsNotCityAdmin = ape.Declare("INITIATOR_IS_NOT_CITY_ADMIN")

func RaiseInitiatorIsNotCityAdmin(cause error, cityID, userID uuid.UUID) error {
	return ErrorInitiatorIsNotCityAdmin.Raise(
		cause,
		responses.InitiatorIsNotCityAdmin(cityID, userID, cause.Error()),
	)
}

var ErrorCityOwnerAlreadyExits = ape.Declare("CITY_OWNER_ALREADY_EXISTS")

func RaiseCityOwnerAlreadyExits(cause error, cityID, userID uuid.UUID) error {
	return ErrorCityOwnerAlreadyExits.Raise(
		cause,
		responses.CityOwnerAlreadyExists(cityID, userID),
	)
}

var ErrorUserIsAlreadyCityAdmin = ape.Declare("USER_IS_ALREADY_CITY_ADMIN")

func RaiseUserIsAlreadyCityAdmin(cause error, cityID, userID uuid.UUID) error {
	return ErrorUserIsAlreadyCityAdmin.Raise(
		cause,
		responses.CityAdminAlreadyExists(cityID, userID),
	)
}

var ErrorInvalidCityAdminRole = ape.Declare("CITY_ADMIN_ROLE_IS_INVALID")

func RaiseInvalidCityAdminRole(role enum.CityAdminRole) error {
	return ErrorInvalidCityAdminRole.Raise(
		fmt.Errorf("invalid city admin role: %s, must be one of %s", role, enum.GetAllCitiesAdminsRoles()),
		responses.InvalidCityAdminRole(role),
	)
}

var ErrorCannotDeleteCityOwner = ape.Declare("CANNOT_DELETE_CITY_OWNER")

func RaiseCannotDeleteCityOwner(cause error, cityID, OwnerID uuid.UUID) error {
	return ErrorCannotDeleteCityOwner.Raise(
		cause,
		responses.CannotDeleteCityOwner(cityID, OwnerID),
	)
}
