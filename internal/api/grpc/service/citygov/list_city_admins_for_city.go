package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) ListCityAdminsForCity(ctx context.Context, req *svc.ListCityAdminsRequest) (*svc.ListCitiesAdmins, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	cityAdmins, err := s.app.GetCityAdmins(ctx, cityID, req.Pagination.Limit, req.Pagination.Page)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to list city admins")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CitiesAdminsList(cityAdmins), nil
}
