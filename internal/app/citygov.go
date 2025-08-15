package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type cityAdminQ interface {
	New() dbx.CityAdminQ

	Insert(ctx context.Context, input dbx.CityAdminModel) error
	Update(ctx context.Context, input dbx.UpdateCityAdmin) error
	Get(ctx context.Context) (dbx.CityAdminModel, error)
	Select(ctx context.Context) ([]dbx.CityAdminModel, error)
	Delete(ctx context.Context) error

	FilterUserID(UserID uuid.UUID) dbx.CityAdminQ
	FilterCityID(cityID uuid.UUID) dbx.CityAdminQ
	FilterRole(role string) dbx.CityAdminQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CityAdminQ
}

func (a App) GetCityAdmin(ctx context.Context, cityID uuid.UUID) (models.CityGov, error) {
	cityAdmin, err := a.adminsQ.New().FilterCityID(cityID).FilterRole(enum.CityAdminRoleAdmin).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.RaiseCityAdminNotFound(ctx, err, cityID)
		default:
			return models.CityGov{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityAdminModel(cityAdmin), nil
}

func (a App) DeleteCityAdmin(ctx context.Context, cityID uuid.UUID) error {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return err
	}

	cityAdmin, err := a.GetCityAdmin(ctx, cityID)
	if err != nil {
		return err
	}

	err = a.adminsQ.New().FilterUserID(cityAdmin.UserID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseCityAdminNotFound(ctx, err, cityID)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}
	return nil
}

func (a App) CreateCityGovAdmin(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	user, err := a.GetCityGov(ctx, cityID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			// No existing admin, proceed to create a new one
		default:
			return models.CityGov{}, err
		}
	}

	if user != (models.CityGov{}) || err == nil {
		return models.CityGov{}, errx.RaiseUserIsAlreadyCityGov(
			ctx,
			fmt.Errorf("user with user_id: %s is already city admin for city_id: %s",
				userID,
				cityID,
			),
			cityID,
			userID,
		)
	}

	admin := dbx.CityAdminModel{
		CityID:    cityID,
		UserID:    userID,
		Role:      enum.CityAdminRoleAdmin,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
	err = a.adminsQ.New().Insert(ctx, admin)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.RaiseCityAlreadyHaveAdmin(
				ctx,
				err,
				cityID,
			)
		default:
			return models.CityGov{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityAdminModel(admin), nil
}

func (a App) CreateCityGovModer(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	user, err := a.GetCityGov(ctx, cityID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			// No existing admin, proceed to create a new one
		default:
			return models.CityGov{}, err
		}
	}

	if user != (models.CityGov{}) || err == nil {
		return models.CityGov{}, errx.RaiseUserIsAlreadyCityGov(
			ctx,
			fmt.Errorf("user with user_id: %s is already city admin for city_id: %s",
				userID,
				cityID,
			),
			cityID,
			userID,
		)
	}

	admin := dbx.CityAdminModel{
		CityID:    cityID,
		UserID:    userID,
		Role:      enum.CityAdminRoleModerator,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
	err = a.adminsQ.New().Insert(ctx, admin)
	if err != nil {
		switch {
		default:
			return models.CityGov{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityAdminModel(admin), nil
}

// Read methods for citygov

// getInitiatorCityGov retrieves the city admin for the given initiator and city.
func (a App) getInitiatorCityGov(ctx context.Context, cityID, initiatorID uuid.UUID) (dbx.CityAdminModel, error) {
	initiator, err := a.adminsQ.New().FilterUserID(initiatorID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return dbx.CityAdminModel{}, errx.RaiseInitiatorIsNotCityGov(
				ctx,
				fmt.Errorf("initiator: %s, not city admin for cit: %s", initiatorID, cityID),
				initiatorID,
				cityID,
			)
		default:
			return dbx.CityAdminModel{}, errx.RaiseInternal(ctx, err)
		}
	}
	return initiator, nil
}

func (a App) GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	cityAdmin, err := a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			_, err := a.GetCityByID(ctx, cityID)
			if err != nil {
				return models.CityGov{}, err
			}

			return models.CityGov{}, errx.RaiseCityGovNotFound(
				ctx,
				fmt.Errorf("city admin not found, cityID: %s, userID: %s, cause %s", cityID, userID, err),
				cityID,
				userID,
			)
		default:
			return models.CityGov{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityAdminModel(cityAdmin), nil
}

func (a App) GetCityGovs(ctx context.Context, cityID uuid.UUID, pag pagination.Request) ([]models.CityGov, pagination.Response, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return nil, pagination.Response{}, err
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	cityAdmins, err := a.adminsQ.New().FilterCityID(cityID).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.CityGov{}, pagination.Response{}, nil
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	total, err := a.adminsQ.New().FilterCityID(cityID).Count(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			total = 0
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	res, pagRes := cityAdminsArray(cityAdmins, limit, offset, total)

	return res, pagRes, nil
}

// Update methods for citygov

// TransferCityAdminRight transfers the admin of a city to a new owner.
func (a App) TransferCityAdminRight(ctx context.Context, cityID, initiatorID, newOwnerID uuid.UUID) error {
	//TODO implement transfer ownership logic

	return nil
}

func (a App) RefuseOwnCityGovRights(ctx context.Context, cityID, userID uuid.UUID) error {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return err
	}

	cityAdmin, err := a.GetCityGov(ctx, cityID, userID)
	if err != nil {
		return err
	}

	if cityAdmin.Role == enum.CityAdminRoleAdmin {
		return errx.RaiseCannotDeleteCityAdmin(
			ctx,
			fmt.Errorf("city admin with user_id:%s cannot delete city admin with user_id: %s, city_id: %s",
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
			return errx.RaiseCityGovNotFound(
				ctx,
				fmt.Errorf("city admin not found, city_id: %s, user_id: %s, cause %s", cityID, userID, err),
				cityID,
				userID,
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	return nil
}

// Delete methods for citygov

func (a App) DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return err
	}

	cityGov, err := a.GetCityGov(ctx, cityID, userID)
	if err != nil {
		return err
	}

	if cityGov.Role == enum.CityAdminRoleAdmin {
		return errx.RaiseCannotDeleteCityAdmin(
			ctx,
			fmt.Errorf("city admin with user_id:%s cannot delete city admin with user_id: %s, city_id: %s",
				userID,
				cityGov.UserID,
				cityGov.CityID,
			),
			userID,
			cityGov.CityID,
		)
	}

	err = a.adminsQ.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		switch {
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	return nil
}

// Helper functions for citygov

func cityAdminModel(cityAdmin dbx.CityAdminModel) models.CityGov {
	return models.CityGov{
		UserID:    cityAdmin.UserID,
		CityID:    cityAdmin.CityID,
		Role:      cityAdmin.Role,
		UpdatedAt: cityAdmin.UpdatedAt,
		CreatedAt: cityAdmin.CreatedAt,
	}
}

func cityAdminsArray(cityAdmins []dbx.CityAdminModel, limit, offset, total uint64) ([]models.CityGov, pagination.Response) {
	res := make([]models.CityGov, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		res[i] = cityAdminModel(cityAdmin)
	}

	pag := pagination.Response{
		Page:  offset/limit + 1,
		Size:  limit,
		Total: total,
	}

	return res, pag
}
