package entities

import (
	"context"
	"errors"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/orb"
)

type City struct {
	queries *dbx.Queries
}

func NewCitySvc(db *pgxpool.Pool) *City {
	return &City{queries: dbx.New(db)}
}

// Create methods for city

func (c City) CreateCity(
	ctx context.Context,
	CountryID uuid.UUID,
	Icon, Slug, TimeZone string,
	Zone orb.MultiPolygon,
) (models.City, error) {
	city, err := c.queries.CreateCity(ctx, dbx.CreateCityParams{
		CountryID:         CountryID,
		Status:            dbx.CityStatusesUnsupported,
		StGeomfromgeojson: Zone,
		Icon:              Icon,
		Slug:              Slug,
		Timezone:          TimeZone,
	})
	if err != nil {
		return models.City{}, problems.RaiseInternal(ctx, err)
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    string(city.Status),
		Zone:      city.Zone,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt.Time,
		UpdatedAt: city.UpdatedAt.Time,
	}, nil
}

func (c City) CreateCityDetails(ctx context.Context, cityID uuid.UUID, name, language string) error {
	lang, err := enum.ParseLanguage(language)
	if err != nil {
		return err
	}

	err = c.queries.CreateCityDetails(ctx, dbx.CreateCityDetailsParams{
		CityID:   cityID,
		Language: dbx.CityLanguages(lang),
		Name:     name,
	})
	if err != nil {
		return problems.RaiseInternal(ctx, err)
	}

	return nil
}

// Read methods for city

func (c City) GetCityByID(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	city, err := c.queries.GetCityByID(ctx, cityID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return models.City{}, problems.RaiseCityNotFoundByID(ctx, err, cityID)
		default:
			return models.City{}, problems.RaiseInternal(ctx, err)
		}
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    string(city.Status),
		Zone:      city.Zone,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt.Time,
		UpdatedAt: city.UpdatedAt.Time,
	}, nil
}

func (c City) GetCityDetails(ctx context.Context, cityID uuid.UUID, language string) (models.CityDetail, error) {
	lang, err := enum.ParseLanguage(language)
	if err != nil {
		return models.CityDetail{}, err
	}

	// 1. Пытаемся получить на нужном языке
	details, err := c.queries.GetCityDetailsByCityIdAndLanguage(ctx, dbx.GetCityDetailsByCityIdAndLanguageParams{
		CityID:   cityID,
		Language: dbx.CityLanguages(lang),
	})
	if err == nil {
		return models.CityDetail{
			Name:     details.Name,
			Language: string(details.Language),
		}, nil
	}

	// 2. Пытаемся fallback на английский
	if errors.Is(err, pgx.ErrNoRows) && lang != enum.LanguageEnglish {
		details, err = c.queries.GetCityDetailsByCityIdAndLanguage(ctx, dbx.GetCityDetailsByCityIdAndLanguageParams{
			CityID:   cityID,
			Language: dbx.CityLanguages(enum.LanguageEnglish),
		})
		if err == nil {
			return models.CityDetail{
				Name:     details.Name,
				Language: string(details.Language),
			}, nil
		}
	}

	// 3. Последний fallback: любой доступный язык
	if errors.Is(err, pgx.ErrNoRows) {
		details, err = c.queries.GetCityDetailsByCityIdAndAnyLanguage(ctx, cityID)
		if err == nil {
			return models.CityDetail{
				Name:     details.Name,
				Language: string(details.Language),
			}, nil
		}
	}

	// Если вообще ничего не нашли — возвращаем проблему
	if errors.Is(err, pgx.ErrNoRows) {
		return models.CityDetail{}, problems.RaiseCityDetailsNotFound(ctx, cityID)
	}

	// Иначе — внутренняя ошибка
	return models.CityDetail{}, problems.RaiseInternal(ctx, err)
}

func (c City) GetNearestCity(ctx context.Context, lon, lat float64) (uuid.UUID, error) {
	cityID, err := c.queries.GetNearestCity(ctx, dbx.GetNearestCityParams{
		StMakepoint:   lon,
		StMakepoint_2: lat,
	})
	if err != nil {
		return uuid.Nil, problems.RaiseInternal(ctx, err)
	}

	return cityID, nil
}

func (c City) GetCityBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	cityID, err := c.queries.GetCityBySlug(ctx, slug)
	if err != nil {
		return uuid.Nil, problems.RaiseInternal(ctx, err)
	}

	return cityID, nil
}

func (c City) SearchCityDetails(
	ctx context.Context,
	namePattern string,
	cityStatuses []string,
	countryIDs []uuid.UUID,
	pagination pagi.Request,
) ([]models.CityDetail, pagi.Response, error) {
	var statuses []dbx.CityStatuses
	for _, s := range cityStatuses {
		status, err := enum.ParseCityStatus(s)
		if err != nil {
			return nil, pagi.Response{}, err
		}
		statuses = append(statuses, dbx.CityStatuses(status))
	}

	res, err := c.queries.SelectCityDetailsByNames(ctx, dbx.SelectCityDetailsByNamesParams{
		NamePattern: namePattern,
		Statuses:    statuses,
		CountryIds:  countryIDs,
		Page:        int64(pagination.Page),
		PageSize:    int64(pagination.Size),
	})
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(ctx, err)
	}

	var details []models.CityDetail
	for _, r := range res {
		details = append(details, models.CityDetail{
			CityID:   r.CityID,
			Name:     r.Name,
			Language: string(r.Language),
		})
	}

	total := 0
	if len(res) > 0 {
		total = int(res[0].TotalCount)
	}

	return details, pagi.Response{
		Page:  pagination.Page,
		Size:  pagination.Size,
		Total: uint64(total),
	}, nil
}

// Update methods for city

func (c City) UpdateCityStatus(ctx context.Context, cityID uuid.UUID, status string) error {
	cityStatus, err := enum.ParseCityStatus(status)
	if err != nil {
		return err
	}

	err = c.queries.UpdateCityStatus(ctx, dbx.UpdateCityStatusParams{
		ID:     cityID,
		Status: dbx.CityStatuses(cityStatus),
	})
	if err != nil {
		return problems.RaiseInternal(ctx, err)
	}

	return nil
}

func (c City) UpdateCityZone(ctx context.Context, cityID uuid.UUID, zone orb.MultiPolygon) error {
	err := c.queries.UpdateCityZone(ctx, dbx.UpdateCityZoneParams{
		ID:                cityID,
		StGeomfromgeojson: zone,
	})
	if err != nil {
		return problems.RaiseInternal(ctx, err)
	}

	return nil
}

func (c City) UpdateCityIcon(ctx context.Context, cityID uuid.UUID, icon string) error {
	err := c.queries.UpdateCityIcon(ctx, dbx.UpdateCityIconParams{
		ID:   cityID,
		Icon: icon,
	})
	if err != nil {
		return problems.RaiseInternal(ctx, err)
	}

	return nil
}

func (c City) UpdateCitySlug(ctx context.Context, cityID uuid.UUID, slug string) error {
	err := c.queries.UpdateCitySlug(ctx, dbx.UpdateCitySlugParams{
		ID:   cityID,
		Slug: slug,
	})
	if err != nil {
		return problems.RaiseInternal(ctx, err)
	}

	return nil
}
