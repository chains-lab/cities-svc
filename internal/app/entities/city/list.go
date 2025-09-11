package city

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterListParams struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Location  *FilterListDistance
}

type FilterListDistance struct {
	Point   orb.Point
	RadiusM uint64
}

func (c City) List(
	ctx context.Context,
	filters FilterListParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
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

	query := c.citiesQ.New()

	if filters.Name != nil {
		query = query.FilterNameLike(*filters.Name)
	}

	for _, s := range filters.Status {
		err := enum.CheckCityStatus(s)
		if err != nil {
			return nil, pagi.Response{}, errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status: %s, cause: %w", s, err),
			)
		}
	}
	if len(filters.Status) > 0 {
		query = query.FilterStatus(filters.Status...)
	}
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Location != nil {
		query = query.FilterWithinRadiusMeters(filters.Location.Point, filters.Location.RadiusM)
	}

	for _, s := range sort {
		switch s.Field {
		case "name":
			query = query.OrderByAlphabetical(s.Ascend)
		case "distance":
			if filters.Location != nil {
				query = query.OrderByNearest(filters.Location.Point, s.Ascend)
			}
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count cities, cause: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search cities, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	cities := make([]models.City, 0, len(rows))
	for _, city := range rows {
		cities = append(cities, cityFromDb(city))
	}

	return cities, pagi.Response{
		Size:  pag.Size,
		Page:  pag.Page,
		Total: total,
	}, nil
}
