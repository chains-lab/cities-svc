package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityStatus(ctx context.Context, req *svc.UpdateCityStatusRequest) (*svc.City, error) {
	initiator, err := s.OnlyCityAdmin(ctx, req.Initiator.UserId, req.CityId, "update city name")
	if err != nil {
		return nil, err
	}

	status, err := enum.ParseCityStatus(req.Status)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problem.InvalidArgumentError(ctx, "city status is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	city, err := s.app.UpdateCitiesStatus(ctx, initiator.CityID, status)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city status")

		return nil, err
	}

	logger.Log(ctx).Infof("city status updated by user %s for city ID %s", initiator.UserID, city.ID)

	return response.City(city), nil
}
