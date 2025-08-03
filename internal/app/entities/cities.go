package entities

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/ape"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type citiesQ interface {
	New() dbx.CitiesQ

	Insert(ctx context.Context, input dbx.CityModels) error
	Update(ctx context.Context, input dbx.CityUpdate) error
	Get(ctx context.Context) (dbx.CityModels, error)
	Select(ctx context.Context) ([]dbx.CityModels, error)
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.CitiesQ
	FilterCountryID(countryID uuid.UUID) dbx.CitiesQ
	FilterStatus(status string) dbx.CitiesQ
	FilterName(name string) dbx.CitiesQ

	SearchName(name string) dbx.CitiesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesQ
}

type Cities struct {
	queries citiesQ
}

func NewCities(db *sql.DB) (Cities, error) {
	return Cities{
		queries: dbx.NewCities(db),
	}, nil
}

type CreateCityInput struct {
	Name      string
	CountryID uuid.UUID
	Status    string
}

func (c Cities) Create(ctx context.Context, input CreateCityInput) error {
	city := dbx.CityModels{
		ID:        uuid.New(),
		CountryID: input.CountryID,
		Name:      input.Name,
		Status:    input.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := c.queries.New().Insert(ctx, city)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c Cities) GetByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	city, err := c.queries.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, ape.RaiseCityNotFoundByID(err, ID)
		default:
			return models.City{}, ape.RaiseInternal(err)
		}
	}

	return CityDbxToModel(city), nil
}

func (c Cities) SearchByNameInCountry(ctx context.Context, prompt string, countryID uuid.UUID, page, limit uint64) ([]models.City, error) {
	cities, err := c.queries.New().FilterCountryID(countryID).SearchName(prompt).Page(page, limit).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.RaiseCitiesNotFoundByNameAndCountryID(err, prompt, countryID)
		default:
			return nil, ape.RaiseInternal(err)
		}
	}

	return CityArrayDbxToModel(cities), nil
}

type UpdateCityInput struct {
	Name      *string
	Status    *string
	CountryID *uuid.UUID
}

func (c Cities) Update(ctx context.Context, ID uuid.UUID, input UpdateCityInput) error {
	cityUpdate := dbx.CityUpdate{
		Name:      input.Name,
		Status:    input.Status,
		CountryID: input.CountryID,
		UpdatedAt: time.Now().UTC(),
	}

	if input.Status != nil {
		if !enum.CheckCityStatus(*input.Status) {
			return ape.RaiseCitiesStatusIsIncorrect(*input.Status)
		}
	}

	err := c.queries.New().FilterID(ID).Update(ctx, cityUpdate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c Cities) Delete(ctx context.Context, ID uuid.UUID) error {
	err := c.queries.New().FilterID(ID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.RaiseCityNotFoundByID(err, ID)
		case err != nil:
			return ape.RaiseInternal(errors.Wrap(err, "failed to delete city"))
		}
	}

	return nil
}

func CityDbxToModel(city dbx.CityModels) models.City {
	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Name:      city.Name,
		Status:    city.Status,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}
}

func CityArrayDbxToModel(cities []dbx.CityModels) []models.City {
	res := make([]models.City, 0, len(cities))
	for _, city := range cities {
		res = append(res, CityDbxToModel(city))
	}
	return res
}
