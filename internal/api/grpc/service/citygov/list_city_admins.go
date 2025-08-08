package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

func (s Service) ListCityAdmins(ctx context.Context, req *svc.ListCityAdminsRequest) (*svc.ListCitiesAdmins, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	cityAdmins, pag, err := s.app.GetCityAdmins(ctx, cityID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to list city admins")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CitiesAdminsList(cityAdmins, pag), nil
}
