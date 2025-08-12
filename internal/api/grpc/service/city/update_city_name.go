package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
)

func (s Service) UpdateCityName(ctx context.Context, req *svc.UpdateCityNameRequest) (*svc.City, error) {
	initiator, err := s.OnlyCityAdmin(ctx, req.Initiator.UserId, req.CityId, "update city name")
	if err != nil {
		return nil, err
	}

	city, err := s.app.UpdateCityName(ctx, initiator.CityID, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city name")

		return nil, err
	}

	logger.Log(ctx).Infof("city name updated by user %s for city ID %s", initiator.UserID, city.ID)

	return response.City(city), nil
}
