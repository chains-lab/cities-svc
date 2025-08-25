package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb"
)

type City struct {
	queries dbx.CitiesQ
}

func NewCitySvc(db *sql.DB) City {
	return City{
		queries: dbx.NewCitiesQ(db),
	}
}

// Create methods for city

func (c City) CreateCity(
	ctx context.Context,
	CountryID uuid.UUID,
	Icon, Slug, TimeZone string,
	Zone orb.MultiPolygon,
) (models.City, error) {
	cityID := uuid.New()
	now := time.Now().UTC()

	err := c.queries.Insert(ctx, dbx.City{
		ID:        cityID,
		CountryID: CountryID,
		Status:    enum.CityStatusSupported,
		Zone:      Zone,
		Icon:      Icon,
		Slug:      Slug,
		Timezone:  TimeZone,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.City{}, problems.RaiseInternal(
			fmt.Errorf("error creating city: %w", err),
		)
	}

	return models.City{
		ID:        cityID,
		CountryID: CountryID,
		Status:    enum.CityStatusSupported,
		Zone:      Zone,
		Icon:      Icon,
		Slug:      Slug,
		Timezone:  TimeZone,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Read methods for city

func (c City) GetCityByID(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	city, err := c.queries.New().FilterID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return models.City{}, problems.RaiseCityNotFoundByID(
				fmt.Errorf("city not found by id: %s, cause: %w", cityID, err),
				fmt.Sprintf("city not found by id: %s", cityID),
			)
		default:
			return models.City{}, problems.RaiseInternal(
				fmt.Errorf("get city by id: %s, cause: %w", cityID, err),
			)
		}
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    city.Status,
		Zone:      city.Zone,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}, nil
}

func (c City) GetNearestCity(ctx context.Context, lon, lat float64) (models.City, error) {
	city, err := c.queries.New().
		FilterStatus(enum.CityStatusSupported).
		OrderByDistance(lon, lat, false).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, problems.RaiseNotFound(
				fmt.Errorf("no supported cities found, cause: %w", err),
				fmt.Sprintf("city not found by id: %s", city.ID),
			)
		default:
			return models.City{}, problems.RaiseInternal(
				fmt.Errorf("get nearest city: %w", err),
			)
		}
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    city.Status,
		Zone:      city.Zone,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}, nil
}

func (c City) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	city, err := c.queries.New().FilterSlug(slug).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, problems.RaiseCityGovNotFound(
				fmt.Errorf("city not found by slug: %s, cause: %w", slug, err),
				fmt.Sprintf("city not found by slug: %s", slug),
			)
		default:
			return models.City{}, problems.RaiseInternal(
				fmt.Errorf("get city by slug: %w", err),
			)
		}
	}

	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Status:    city.Status,
		Zone:      city.Zone,
		Icon:      city.Icon,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}, nil
}

func (c City) SearchCityDetails(
	ctx context.Context,
	name string,
	statuses []string,
	countryIDs []uuid.UUID,
	pag pagi.Request,
) ([]models.City, pagi.Response, error) {
	queries := c.queries.New()

	if name != "" {
		queries = queries.FilterNameLike(name)
	}
	if len(statuses) > 0 {
		queries = queries.FilterStatus(statuses...)
	}
	if len(countryIDs) > 0 {
		queries = queries.FilterCountryID(countryIDs...)
	}

	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1 // +1 чтобы определить наличие next
	offset := (pag.Page - 1) * pag.Size

	rows, err := queries.Page(limit, offset).OrderByAlphabetical(true).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(
			fmt.Errorf("search city details: %w", err),
		)
	}

	prev := pag.Page > 1
	next := len(rows) > int(pag.Size)
	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	response := make([]models.City, 0, len(rows))
	for _, city := range rows {
		response = append(response, models.City{
			ID:        city.ID,
			CountryID: city.CountryID,
			Status:    city.Status,
			Zone:      city.Zone,
			Name:      city.Name,
			Icon:      city.Icon,
			Slug:      city.Slug,
			Timezone:  city.Timezone,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		})
	}

	return response, pagi.Response{
		Size: pag.Size,
		Page: pag.Page,
		Next: next,
		Prev: prev,
	}, nil
}

// Update methods for city

func (c City) UpdateCityStatus(ctx context.Context, cityID uuid.UUID, status string) error {
	cityStatus, err := enum.ParseCityStatus(status)
	if err != nil {
		return err
	}

	err = c.queries.New().FilterID(cityID).Update(ctx, map[string]any{
		"status":     cityStatus,
		"updated_at": time.Now(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("updating city status: %w", err),
		)
	}

	return nil
}

func (c City) UpdateCityZone(ctx context.Context, cityID uuid.UUID, zone orb.MultiPolygon) error {
	err := c.queries.New().FilterID(cityID).Update(ctx, map[string]any{
		"zone":       zone,
		"updated_at": time.Now(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("updating city zone: %w", err),
		)
	}

	return nil
}

func (c City) UpdateCityIcon(ctx context.Context, cityID uuid.UUID, icon string) error {
	err := c.queries.New().FilterID(cityID).Update(ctx, map[string]any{
		"icon":       icon,
		"updated_at": time.Now(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("updating city icon: %w", err),
		)
	}

	return nil
}

func (c City) UpdateCitySlug(ctx context.Context, cityID uuid.UUID, slug string) error {
	err := c.queries.New().FilterID(cityID).Update(ctx, map[string]any{
		"slug":       slug,
		"updated_at": time.Now(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("updating city slug: %w", err),
		)
	}

	return nil
}
