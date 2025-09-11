package country

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
)

type FilterListParams struct {
	Name     string
	Statuses []string
}

func (c Country) List(
	ctx context.Context,
	filters FilterListParams,
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
			if err := enum.CheckCountryStatus(s); err != nil {
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
