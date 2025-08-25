package entities

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Country struct {
	queries *dbx.Queries
}

func NewCountryService(db *pgxpool.Pool) *Country {
	return &Country{queries: dbx.New(db)}
}

// Create methods for countries

func (a Country) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	country, err := a.queries.CreateCountry(ctx, dbx.CreateCountryParams{
		Name:   name,
		Status: enum.CountryStatusUnsupported,
	})
	if err != nil {
		return models.Country{}, err
	}

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    string(country.Status),
		CreatedAt: country.CreatedAt.Time,
		UpdatedAt: country.UpdatedAt.Time,
	}, nil
}

// Read methods for countries

func (a Country) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.queries.GetCountryByID(ctx, ID)
	if err != nil {
		return models.Country{}, err
	}

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    string(country.Status),
		CreatedAt: country.CreatedAt.Time,
		UpdatedAt: country.UpdatedAt.Time,
	}, nil
}

func (a Country) SearchCountries(
	ctx context.Context,
	name string,
	status []string,
	pag pagi.Request,
) ([]models.Country, pagi.Response, error) {
	statuses := make([]dbx.CountryStatuses, 0, len(status))
	for _, s := range status {
		st, err := enum.ParseCountryStatus(s)
		if err != nil {
			return nil, pagi.Response{}, err
		}

		statuses = append(statuses, dbx.CountryStatuses(st))
	}

	countries, err := a.queries.SearchBeNameAndStatuses(ctx, dbx.SearchBeNameAndStatusesParams{
		Page:        int64(pag.Page),
		PageSize:    int64(pag.Size),
		NamePattern: name,
		Statuses:    statuses,
	})
	if err != nil {
		return nil, pagi.Response{}, err
	}

	var el []models.Country
	for _, country := range countries {
		el = append(el, models.Country{
			ID:        country.ID,
			Name:      country.Name,
			Status:    string(country.Status),
			CreatedAt: country.CreatedAt.Time,
			UpdatedAt: country.UpdatedAt.Time,
		})
	}

	total := 0
	if len(countries) > 0 {
		total = int(countries[0].TotalCount)
	}

	return el, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: uint64(total),
	}, nil
}

// Update methods for countries

func (a Country) UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status string) (models.Country, error) {
	st, err := enum.ParseCountryStatus(status)
	if err != nil {
		return models.Country{}, err
	}

	country, err := a.queries.UpdateCountryStatus(ctx, dbx.UpdateCountryStatusParams{
		ID:     ID,
		Status: dbx.CountryStatuses(st),
	})
	if err != nil {
		return models.Country{}, err
	}

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    string(country.Status),
		CreatedAt: country.CreatedAt.Time,
		UpdatedAt: country.UpdatedAt.Time,
	}, nil
}

func (a Country) UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error) {
	country, err := a.queries.UpdateCountryName(ctx, dbx.UpdateCountryNameParams{
		ID:   ID,
		Name: name,
	})
	if err != nil {
		return models.Country{}, err
	}

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    string(country.Status),
		CreatedAt: country.CreatedAt.Time,
		UpdatedAt: country.UpdatedAt.Time,
	}, nil
}
