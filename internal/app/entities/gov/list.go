package gov

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FiltersListParams struct {
	CityID *uuid.UUID
	Role   []string
}

func (g Gov) SelectGovs(
	ctx context.Context,
	filters FiltersListParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Gov, pagi.Response, error) {
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

	query := g.govQ.New()

	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil && len(filters.Role) > 0 {
		for _, r := range filters.Role {
			err := constant.CheckCityGovRole(r)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidGovRole.Raise(
					fmt.Errorf("invalid city gov role, cause: %w", err),
				)
			}
		}
		query = query.FilterRole(filters.Role...)
	}

	for _, s := range sort {
		switch s.Field {
		case "role":
			query = query.OrderByRole(s.Ascend)
		case "created_at":
			query = query.OrderByCreatedAt(s.Ascend)
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count city govs, cause: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to select city govs, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	res := make([]models.Gov, 0, len(rows))
	for _, cg := range rows {
		res = append(res, govFromDb(cg))
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}
