package entities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/ape"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
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
	FilterRole(role string) dbx.CitiesAdminsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesAdminsQ
}

type CitiesAdmins struct {
	queries citiesAdminsQ
}

func NewCitiesAdmins(db *sql.DB) (CitiesAdmins, error) {
	return CitiesAdmins{
		queries: dbx.NewCitiesAdmins(db),
	}, nil
}

type AddCityAdminInput struct {
	Role string
}

func (c CitiesAdmins) Add(ctx context.Context, initiatorID uuid.UUID, UserID uuid.UUID, CityID uuid.UUID, input AddCityAdminInput) error {
	if !enum.CheckRole(input.Role) {
		return ape.RaiseInvalidCityAdminRole(input.Role)
	}

	initiator, err := c.queries.New().FilterUserID(initiatorID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	_, err = c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	if initiator.Role <= input.Role {
		return ape.RaiseInitiatorCityAdminHaveNotEnoughRights(UserID, CityID)
	}

	cityAdmin := dbx.CityAdminModel{
		CityID:    CityID,
		UserID:    UserID,
		Role:      input.Role,
		CreatedAt: time.Now().UTC(),
	}

	err = c.queries.New().Insert(ctx, cityAdmin)
	if err != nil {
		switch {
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c CitiesAdmins) UpdateRole(ctx context.Context, initiatorID uuid.UUID, UserID uuid.UUID, CityID uuid.UUID, role string) error {
	if !enum.CheckRole(role) {
		return ape.RaiseInvalidCityAdminRole(role)
	}

	initiator, err := c.queries.New().FilterUserID(initiatorID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	updatedUser, err := c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	if initiator.Role <= updatedUser.Role || initiator.Role <= role {
		return ape.RaiseInitiatorCityAdminHaveNotEnoughRights(UserID, CityID)
	}

	err = c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Update(ctx, dbx.UpdateCityAdmin{
		Role:      &role,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		switch {
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c CitiesAdmins) Delete(ctx context.Context, initiatorID uuid.UUID, UserID uuid.UUID, CityID uuid.UUID) error {
	initiator, err := c.queries.New().FilterUserID(initiatorID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	updatedUser, err := c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInitiatorIsNotAdmin(initiatorID, CityID)
		default:
			return ape.RaiseInternal(err)
		}
	}

	if initiator.Role <= updatedUser.Role {
		return ape.RaiseInitiatorCityAdminHaveNotEnoughRights(UserID, CityID)
	}

	err = c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Delete(ctx)
	if err != nil {
		switch {
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c CitiesAdmins) GetByUserIDAndCityID(ctx context.Context, UserID uuid.UUID, CityID uuid.UUID) (models.CityAdmin, error) {
	cityAdmin, err := c.queries.New().FilterUserID(UserID).FilterCityID(CityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityAdmin{}, ape.RaiseCityAdminNotFound(err, UserID, CityID)
		default:
			return models.CityAdmin{}, ape.RaiseInternal(err)
		}
	}

	return models.CityAdmin{
		UserID:    cityAdmin.UserID,
		CityID:    cityAdmin.CityID,
		Role:      cityAdmin.Role,
		UpdatedAt: cityAdmin.UpdatedAt,
		CreatedAt: cityAdmin.CreatedAt,
	}, nil
}

func (c CitiesAdmins) GetByCityID(ctx context.Context, CityID uuid.UUID, page, limit uint64) ([]models.CityAdmin, error) {
	cityAdmins, err := c.queries.New().FilterCityID(CityID).Page(page, limit).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.RaiseCityNotFoundByID(err, CityID)
		default:
			return nil, ape.RaiseInternal(err)
		}
	}

	res := make([]models.CityAdmin, 0, len(cityAdmins))
	for _, cityAdmin := range cityAdmins {
		res = append(res, models.CityAdmin{
			UserID:    cityAdmin.UserID,
			CityID:    cityAdmin.CityID,
			Role:      cityAdmin.Role,
			UpdatedAt: cityAdmin.UpdatedAt,
			CreatedAt: cityAdmin.CreatedAt,
		})
	}

	return res, nil
}
