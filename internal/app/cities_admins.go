package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	errs2 "github.com/chains-lab/cities-dir-svc/internal/errs"
	"github.com/google/uuid"
)

type citiesAdminsQ interface {
	New() dbx.CitiesAdminsQ

	Insert(ctx context.Context, input dbx.CityAdminModel) error
	Update(ctx context.Context, input dbx.UpdateCityAdmin) error
	Get(ctx context.Context) (dbx.CityAdminModel, error)
	Select(ctx context.Context) ([]dbx.CityAdminModel, error)
	Delete(ctx context.Context) error

	FilterUserID(UserID uuid.UUID) dbx.CitiesAdminsQ
	FilterCityID(cityID uuid.UUID) dbx.CitiesAdminsQ
	FilterRole(role enum.CityAdminRole) dbx.CitiesAdminsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesAdminsQ
}

// GetInitiatorCityAdmin retrieves the city admin for the given initiator and city.
func (a App) GetInitiatorCityAdmin(ctx context.Context, cityID, initiatorID uuid.UUID) (dbx.CityAdminModel, error) {
	initiator, err := a.adminsQ.New().FilterUserID(initiatorID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return dbx.CityAdminModel{}, errs2.RaiseInitiatorIsNotCityAdmin(
				fmt.Errorf("initiator: %s, not city admin for cit: %s", initiatorID, cityID),
				initiatorID,
				cityID,
			)
		default:
			return dbx.CityAdminModel{}, errs2.RaiseInternal(err)
		}
	}
	return initiator, nil
}

func (a App) CreateCityOwner(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error) {
	ID := uuid.New()

	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	owner, err := a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).FilterRole(enum.CityOwner).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// No existing owner, proceed to create a new one
			break
		default:
			return models.CityAdmin{}, errs2.RaiseInternal(err)
		}
	}
	if err == nil || owner != (dbx.CityAdminModel{}) {
		return models.CityAdmin{}, errs2.RaiseCityOwnerAlreadyExits(
			fmt.Errorf("city owner already exists for user_id: %s, city_id: %s",
				userID,
				cityID,
			),
			owner.UserID,
			cityID,
		)
	}

	cityAdmin := dbx.CityAdminModel{
		ID:        ID,
		CityID:    cityID,
		UserID:    userID,
		Role:      enum.CityOwner,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = a.adminsQ.New().Insert(ctx, cityAdmin)
	if err != nil {
		switch {
		default:
			return models.CityAdmin{}, errs2.RaiseInternal(err)
		}
	}

	return cityAdminDbxToModel(cityAdmin), nil
}

func (a App) TransferCityOwnership(ctx context.Context, initiatorID, newOwnerID, cityID uuid.UUID) error {
	//TODO implement transfer ownership logic

	return nil
}

func (a App) GetCityAdmin(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error) {
	cityAdmin, err := a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityAdmin{}, errs2.RaiseCityAdminNotFound(
				fmt.Errorf("city admin not found, cityID: %s, userID: %s, cause %s", cityID, userID, err),
				cityID,
				userID,
			)
		default:
			return models.CityAdmin{}, errs2.RaiseInternal(err)
		}
	}

	return cityAdminDbxToModel(cityAdmin), nil
}

func (a App) GetCityAdminForCity(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error) {
	cityAdmin, err := a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityAdmin{}, errs2.RaiseCityAdminNotFound(
				fmt.Errorf("city admin not found for city_id: %s, user_id: %s, cause %s", cityID, userID, err),
				cityID,
				userID,
			)
		default:
			return models.CityAdmin{}, errs2.RaiseInternal(err)
		}
	}

	return cityAdminDbxToModel(cityAdmin), nil
}

func (a App) GetUserCitiesAdmins(ctx context.Context, userID uuid.UUID, limit, page uint64) ([]models.CityAdmin, error) {
	cityAdmins, err := a.adminsQ.New().FilterUserID(userID).Page(limit, page).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.CityAdmin{}, nil
		default:
			return nil, errs2.RaiseInternal(err)
		}
	}

	return cityArrayAdminDbxToModel(cityAdmins), nil
}

type CreateCityAdminInput struct {
	Role enum.CityAdminRole
}

func (a App) CreateCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID, input CreateCityAdminInput) (models.CityAdmin, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	initiator, err := a.GetInitiatorCityAdmin(ctx, cityID, initiatorID)

	user, err := a.GetCityAdmin(ctx, cityID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errs2.ErrorCityAdminNotFound):
			// No existing admin, proceed to create a new one
		default:
			return models.CityAdmin{}, err
		}
	}
	if user != (models.CityAdmin{}) {
		return models.CityAdmin{}, errs2.RaiseUserIsAlreadyCityAdmin(
			fmt.Errorf("user with user_id: %s is already city admin for city_id: %s",
				userID,
				cityID,
			),
			cityID,
			userID,
		)
	}

	role, ok := enum.ParseCityAdminRole(string(input.Role))
	if !ok {
		return models.CityAdmin{}, errs2.RaiseInvalidCityAdminRole(input.Role)
	}

	if enum.CompareCityAdminRole(initiator.Role, role) < 1 {
		return models.CityAdmin{}, errs2.RaiseCityAdminHaveNotEnoughRights(
			fmt.Errorf("initiator user_id: %s has role %s, but trying to create admin with user_id: %s, city_id: %s, role %s",
				initiatorID,
				input.Role,
				userID,
				cityID,
				role),
			userID,
			cityID,
		)
	}

	admin := dbx.CityAdminModel{
		CityID:    cityID,
		UserID:    userID,
		Role:      role,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
	err = a.adminsQ.New().Insert(ctx, admin)
	if err != nil {
		switch {
		default:
			return models.CityAdmin{}, errs2.RaiseInternal(err)
		}
	}

	return cityAdminDbxToModel(admin), nil
}

func (a App) UpdateCityAdminRole(ctx context.Context, initiatorID, cityID, userID uuid.UUID, role enum.CityAdminRole) error {
	_, ok := enum.ParseCityAdminRole(string(role))
	if !ok {
		return errs2.RaiseInvalidCityAdminRole(role)
	}

	initiator, err := a.GetInitiatorCityAdmin(ctx, cityID, initiatorID)
	if err != nil {
		return err
	}

	user, err := a.GetCityAdminForCity(ctx, cityID, userID)
	if err != nil {
		return err
	}

	if enum.CompareCityAdminRole(initiator.Role, user.Role) < 1 {
		return errs2.RaiseCityAdminHaveNotEnoughRights(
			fmt.Errorf("initiator user_id: %s has role %s, but trying to update admin with user_id: %s, city_id: %s, role %s",
				initiatorID,
				initiator.Role,
				userID,
				cityID,
				role),
			cityID,
			userID)
	}

	updateInput := dbx.UpdateCityAdmin{
		Role:      &role,
		UpdatedAt: time.Now().UTC(),
	}

	err = a.adminsQ.New().FilterUserID(userID).Update(ctx, updateInput)
	if err != nil {
		switch {
		default:
			return errs2.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) RefuseOwnAdminRights(ctx context.Context, cityID, userID uuid.UUID) error {
	cityAdmin, err := a.GetCityAdmin(ctx, cityID, userID)
	if err != nil {
		return err
	}

	if cityAdmin.Role == enum.CityOwner {
		return errs2.RaiseCannotDeleteCityOwner(
			fmt.Errorf("city admin with user_id:%s cannot delete city owner with user_id: %s, city_id: %s",
				userID,
				cityAdmin.UserID,
				cityAdmin.CityID,
			),
			userID,
			cityAdmin.CityID,
		)
	}

	//TODO: add more safety checks here, for example use email confirmation, realize it in the future with kafka
	err = a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errs2.RaiseCityAdminNotFound(
				fmt.Errorf("city admin not found, city_id: %s, user_id: %s, cause %s", cityID, userID, err),
				cityID,
				userID,
			)
		default:
			return errs2.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) DeleteCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID) error {
	initiator, err := a.GetInitiatorCityAdmin(ctx, cityID, initiatorID)
	if err != nil {
		return err
	}

	user, err := a.GetCityAdminForCity(ctx, cityID, userID)
	if err != nil {
		return err
	}

	if enum.CompareCityAdminRole(initiator.Role, user.Role) < 1 {
		return errs2.RaiseCityAdminHaveNotEnoughRights(
			fmt.Errorf("initiator user_id: %s has role %s, but trying to delete admin with user_id: %s, city_id: %s, role %s",
				initiatorID,
				initiator.Role,
				userID,
				cityID,
				user.Role,
			),
			initiatorID,
			userID,
		)
	}

	err = a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		switch {
		default:
			return errs2.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) GetCityAdmins(ctx context.Context, cityID uuid.UUID, limit, page uint64) ([]models.CityAdmin, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return nil, err
	}

	cityAdmins, err := a.adminsQ.New().FilterCityID(cityID).Page(limit, page).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.CityAdmin{}, nil
		default:
			return nil, errs2.RaiseInternal(err)
		}
	}

	return cityArrayAdminDbxToModel(cityAdmins), nil
}

func cityAdminDbxToModel(cityAdmin dbx.CityAdminModel) models.CityAdmin {
	return models.CityAdmin{
		UserID:    cityAdmin.UserID,
		CityID:    cityAdmin.CityID,
		Role:      cityAdmin.Role,
		UpdatedAt: cityAdmin.UpdatedAt,
		CreatedAt: cityAdmin.CreatedAt,
	}
}

func cityArrayAdminDbxToModel(cityAdmins []dbx.CityAdminModel) []models.CityAdmin {
	modelsCityAdmins := make([]models.CityAdmin, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		modelsCityAdmins[i] = cityAdminDbxToModel(cityAdmin)
	}
	return modelsCityAdmins
}
