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
)

type Country struct {
	queries dbx.CountryQ
}

func NewCountry(db *sql.DB) Country {
	return Country{
		queries: dbx.NewCountryQ(db),
	}
}

// Create methods for countries

func (a Country) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	now := time.Now().UTC()
	ID := uuid.New()
	err := a.queries.New().Insert(ctx, dbx.Country{
		ID:        ID,
		Name:      name,
		Status:    enum.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.Country{}, problems.RaiseInternal(
			fmt.Errorf("error creating country: %w", err),
		)
	}

	return models.Country{
		ID:        ID,
		Name:      name,
		Status:    enum.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Read methods for countries

func (a Country) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.queries.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, problems.RaiseCountryNotFound(
				fmt.Errorf("country not found: %w", err),
				fmt.Sprintf("country with ID %s not found", ID),
			)
		default:
			return models.Country{}, problems.RaiseInternal(
				fmt.Errorf("get country by ID: %w", err),
			)
		}
	}

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    country.Status,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}, nil
}

func (a Country) SearchCountries(
	ctx context.Context,
	name string,
	status []string,
	pag pagi.Request,
) ([]models.Country, pagi.Response, error) {
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

	rows, err := a.queries.New().FilterName(name).FilterStatus(status...).Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(
			fmt.Errorf("search countries: %w", err),
		)
	}

	prev := pag.Page > 1
	next := len(rows) > int(pag.Size)
	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	countries := make([]models.Country, 0, len(rows))
	for _, country := range rows {
		countries = append(countries, models.Country{
			ID:        country.ID,
			Name:      country.Name,
			Status:    country.Status,
			CreatedAt: country.CreatedAt,
			UpdatedAt: country.UpdatedAt,
		})
	}

	return countries, pagi.Response{
		Page: pag.Page,
		Size: pag.Size,
		Prev: prev,
		Next: next,
	}, nil
}

// Update methods for countries

func (a Country) UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status string) error {
	st, err := enum.ParseCountryStatus(status)
	if err != nil {
		return problems.RaiseInvalidCountryStatus(
			fmt.Errorf("parse country status: %w", err),
			fmt.Sprintf("invalid country status: %s", status),
		)
	}

	err = a.queries.New().FilterID(ID).Update(ctx, map[string]any{
		"status":     st,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("update country status: %w", err),
		)
	}

	return nil
}

func (a Country) UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) error {
	err := a.queries.New().FilterID(ID).Update(ctx, map[string]any{
		"name":       name,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("update country name: %w", err),
		)
	}

	return nil
}
