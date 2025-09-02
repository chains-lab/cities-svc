package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type Country struct {
	countryQ dbx.CountryQ
}

func NewCountry(db *sql.DB) Country {
	return Country{
		countryQ: dbx.NewCountryQ(db),
	}
}

func (c Country) Create(ctx context.Context, name string, status string) (models.Country, error) {
	now := time.Now().UTC()
	ID := uuid.New()

	err := constant.CheckCountryStatus(status)
	if err != nil {
		return models.Country{}, errx.ErrorInvalidCountryStatus.Raise(
			fmt.Errorf("failed to parse country status: %w", err),
		)
	}

	err = c.countryQ.New().Insert(ctx, dbx.Country{
		ID:        ID,
		Name:      name,
		Status:    constant.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating country: %w", err),
		)
	}

	return models.Country{
		ID:        ID,
		Name:      name,
		Status:    constant.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Read methods for countries

func (c Country) GetByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := c.countryQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, errx.ErrorCountryNotFound.Raise(
				fmt.Errorf("country not found, cause: %w", err),
			)
		default:
			return models.Country{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get country by ID, cause: %w", err),
			)
		}
	}

	return countryFromDb(country), nil
}

func (c Country) GetByName(ctx context.Context, name string) (models.Country, error) {
	country, err := c.countryQ.New().FilterName(name).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, errx.ErrorCountryNotFound.Raise(
				fmt.Errorf("country not found, cause: %w", err),
			)
		default:
			return models.Country{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get country by name, cause: %w", err),
			)
		}
	}

	return countryFromDb(country), nil
}

// Select methods for countries

type FilterCountriesFilters struct {
	Name     string
	Statuses []string
}

func (c Country) Select(
	ctx context.Context,
	filters FilterCountriesFilters,
	pag pagi.Request,
	sort []pagi.SortField,
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

	query := c.countryQ.New()

	if filters.Statuses != nil {
		for _, s := range filters.Statuses {
			if err := constant.CheckCountryStatus(s); err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidCountryStatus.Raise(
					fmt.Errorf("failed to parse country status, cause: %w", err),
				)
			}
		}

		query = query.FilterStatus(filters.Statuses...)
	}

	if filters.Name != "" {
		query = query.FilterNameLike(filters.Name)
	}

	for _, sort := range sort {
		switch sort.Field {
		case "name":
			query = query.OrderAlphabetical(sort.Ascend)
		default:

		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count countries, cause: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search countries, cause: %w", err),
		)
	}

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
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

// UpdateOne methods for countries

type UpdateCountryParams struct {
	Name      *string
	Status    *string
	UpdatedAt time.Time
}

func (c Country) Update(ctx context.Context, ID uuid.UUID, params UpdateCountryParams) error {
	stmt := dbx.UpdateCountryParams{}

	if params.Name == nil && params.Status == nil {
		return nil
	}

	if params.Name != nil {
		stmt.Name = params.Name
	}
	if params.Status != nil {
		err := constant.CheckCountryStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidCountryStatus.Raise(
				fmt.Errorf("failed to invalid country status, cause: %w", err),
			)
		}
	}
	stmt.UpdatedAt = params.UpdatedAt

	err := c.countryQ.New().FilterID(ID).Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update country, cause: %w", err),
		)
	}

	return nil
}

func countryFromDb(c dbx.Country) models.Country {
	return models.Country{
		ID:        c.ID,
		Name:      c.Name,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
