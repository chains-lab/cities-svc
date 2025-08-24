package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type citiesQ interface {
	New() dbx.CitiesQ

	Insert(ctx context.Context, input dbx.City) error
	Update(ctx context.Context, input map[string]any) error
	Get(ctx context.Context) (dbx.City, error)
	Select(ctx context.Context) ([]dbx.City, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CitiesQ
	FilterCountryID(countryID uuid.UUID) dbx.CitiesQ
	FilterStatus(status string) dbx.CitiesQ
	FilterSlug(slug string) dbx.CitiesQ

	FilterWithinRadiusMeters(lon, lat float64, radiusM uint64) dbx.CitiesQ
	OrderByDistance(lon, lat float64, asc bool) dbx.CitiesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesQ

	Transaction(fn func(ctx context.Context) error) error
}

type CityDetailsQ interface {
	New() dbx.CityDetailsQ

	Insert(ctx context.Context, input dbx.CityDetail) error
	Update(ctx context.Context, input map[string]any) error
	Get(ctx context.Context) (dbx.CityDetail, error)
	Select(ctx context.Context) ([]dbx.CityDetail, error)
	Delete(ctx context.Context) error

	FilterCityID(cityID uuid.UUID) dbx.CityDetailsQ
	FilterLanguage(language string) dbx.CityDetailsQ
	SearchName(like string) dbx.CityDetailsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CityDetailsQ
}

type City struct {
	city    citiesQ
	details CityDetailsQ
}

func NewCity(cfg config.Config) (City, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return City{}, err
	}

	return City{
		city:    dbx.NewCitiesQ(pg),
		details: dbx.NewCityDetailsQ(pg),
	}, nil
}

// Create methods for city

type CreateCityInput struct {
	CountryID uuid.UUID
	Name      string
	Status    string
	Icon      string
	Slug      string
	TimeZone  string
	Center    orb.Point        // [lon, lat]
	Boundary  orb.MultiPolygon // многоугольники границы

	Details []models.CityDetail
}

func (c City) CreateCity(ctx context.Context, input CreateCityInput) error {
	status, err := enum.ParseCityStatus(input.Status)
	if err != nil {
		return problems.RaiseInvalidCityStatus(ctx, err, input.Status)
	}

	ID := uuid.New()

	city := dbx.City{
		ID:        ID,
		CountryID: input.CountryID,
		Status:    status,
		Icon:      input.Icon,
		Slug:      input.Slug,
		Timezone:  input.TimeZone,
		Center:    input.Center,
		Boundary:  input.Boundary,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if len(input.Details) < 1 {
		return err //TODO
	}

	txErr := c.city.Transaction(func(txCtx context.Context) error {
		err = c.city.New().Insert(ctx, city)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return problems.RaiseInternal(ctx, err)
			}
		}

		for _, el := range input.Details {
			detail := dbx.CityDetail{
				CityID: ID,
				Name:   el.Name,
			}

			lang, err := enum.ParseLanguage(el.Language)
			if err != nil {
				return err //TODO
			}

			if el.Description != nil {
				detail.Description = el.Description
			}
			detail.Language = lang

			err = c.details.New().Insert(ctx, detail)
			if err != nil {
				switch {
				case errors.Is(err, sql.ErrNoRows):
					return problems.RaiseInternal(ctx, err)
				}
			}
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

// Read methods for city

func (c City) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	city, err := c.city.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, problems.RaiseCityNotFoundByID(
				ctx,
				fmt.Errorf("city with ID '%s' not found cause: %s", ID, err),
				ID,
			)
		default:
			return models.City{}, problems.RaiseInternal(ctx, err)
		}
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    city.Status,
		Center:    city.Center,
		Boundary:  city.Boundary,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,

		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}, nil
}

func (c City) GetCityDetails(ID uuid.UUID, language string) (models.CityDetail, error) {
	details, err := c.details.New().FilterCityID(ID).FilterLanguage(language).Get(context.Background())
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			details, err = c.details.New().FilterCityID(ID).FilterLanguage("en").Get(context.Background())
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return models.CityDetail{}, nil
			}
		default:
			return models.CityDetail{}, problems.RaiseInternal(context.Background(), err)
		}
	}

	res := models.CityDetail{
		Name:     details.Name,
		Language: details.Language,
	}

	if details.Description != nil {
		res.Description = details.Description
	}

	return res, nil
}

func (c City) GetCityWithDetails(ctx context.Context, ID uuid.UUID, language string) (models.City, error) {
	res, err := c.GetCityByID(ctx, ID)
	if err != nil {
		return models.City{}, err
	}

	details, err := c.GetCityDetails(ID, language)
	if err != nil {
		return models.City{}, err
	}

	res.Details = details
	return res, nil
}

func (c City) GetNearestCity(
	ctx context.Context,
	lon, lat float64,
	radiusM, limit, offset uint64,
) (models.City, error) {
	city, err := c.city.New().
		FilterWithinRadiusMeters(lon, lat, radiusM).
		OrderByDistance(lon, lat, true).
		Page(limit, offset).
		Get(ctx)
	if err != nil {
		return models.City{}, err //TODO
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    city.Status,
		Center:    city.Center,
		Boundary:  city.Boundary,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,

		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}, nil
}

func (c City) SearchCitiesByName(
	ctx context.Context,
	name string,
	pag pagi.Request,
) ([]models.City, pagination.Response, error) {
	limit, offset := pagi.CalculateLimitOffset(pag)

	details, err := c.details.SearchName(name).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.City{}, pagination.Response{}, nil
		default:
			return nil, pagination.Response{}, problems.RaiseInternal(ctx, err)
		}
	}

	cities := make([]models.City, 0, len(details))

	for _, detail := range details {
		city, err := c.GetCityByID(ctx, detail.CityID)
		if err != nil {
			return nil, pagination.Response{}, err
		}

		city.Details = models.CityDetail{
			Name:        detail.Name,
			Description: detail.Description,
			Language:    detail.Language,
		}
	}

	count, err := c.details.SearchName(name).Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, err
	}

	return cities, pagination.Response{
		Total: count,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}

func (c City) GetCountryCities(
	ctx context.Context,
	countryID uuid.UUID,
	pag pagi.Request,
) ([]models.City, pagination.Response, error) {
	limit, offset := pagi.CalculateLimitOffset(pag)

	count, err := c.city.New().FilterCountryID(countryID).Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, err
	}

	cities, err := c.city.New().FilterCountryID(countryID).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.City{}, pagination.Response{}, nil
		default:
			return nil, pagination.Response{}, problems.RaiseInternal(ctx, err)
		}
	}

	res := make([]models.City, 0, len(cities))
	for _, city := range cities {
		res = append(res, models.City{
			ID:        city.ID,
			CountryID: city.CountryID,
			Status:    city.Status,
			Center:    city.Center,
			Boundary:  city.Boundary,
			Icon:      city.Icon,
			Slug:      city.Slug,
			Timezone:  city.Timezone,

			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		})
	}

	return res, pagination.Response{
		Total: count,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}

// Update methods for city

func (c City) UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	_, err := enum.ParseCityStatus(status)
	if err != nil {
		return models.City{}, problems.RaiseInvalidCityStatus(ctx, err, status)
	}

	err = c.city.New().FilterID(cityID).Update(ctx, map[string]any{
		"status": status,
	})
	if err != nil {
		return models.City{}, err
	}

	return c.GetCityByID(ctx, cityID)
}

func (c City) UpdateCitiesDetails(ctx context.Context, cityID uuid.UUID, details models.CityDetail) (models.CityDetail, error) {
	_, err := c.GetCityByID(ctx, cityID)
	if err != nil {
		return models.CityDetail{}, err
	}

	lang, err := enum.ParseLanguage(details.Language)
	if err != nil {
		return models.CityDetail{}, err
	}

	updateData := map[string]any{
		"name": details.Name,
	}
	if details.Description != nil {
		updateData["description"] = details.Description
	}

	err = c.details.New().FilterCityID(cityID).FilterLanguage(lang).Update(ctx, updateData)
	if err != nil {
		return models.CityDetail{}, err
	}

	res, err := c.GetCityDetails(cityID, details.Language)
	if err != nil {
		return models.CityDetail{}, err
	}

	return res, nil
}

func (c City) Transaction(fn func(ctx context.Context) error) error {
	return c.city.New().Transaction(fn)
}
