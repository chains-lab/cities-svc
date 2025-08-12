package admin

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateCity(ctx context.Context, req *svc.SendFormToCreateCityRequest) (*svc.City, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "send form to create city", roles.User)
	if err != nil {
		return nil, err
	}

	CountryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return nil, problem.InvalidArgumentError(ctx, fmt.Sprintf("country id is invalid"), &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	city, err := s.app.CreateCity(ctx, app.CreateCityInput{
		Name:      req.Name,
		CountryID: CountryID,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create city")

		return nil, err
	}

	logger.Log(ctx).Infof("created city with ID %s by user %s", city.ID, req.Initiator.UserId)

	return response.City(city), nil
}
