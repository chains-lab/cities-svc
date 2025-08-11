package city

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) CreateCity(ctx context.Context, req *svc.CreateCityRequest) (*svc.City, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, problems.UnauthenticatedError(ctx, "initiator role is invalid format")
	}

	if role != roles.Admin && role != roles.SuperUser {
		logger.Log(ctx, RequestID(ctx)).Warnf("user %s with role %s tried to create a city, but only admins and superusers can create cities",
			req.Initiator.Id, req.Initiator.Role)

		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"user %s with role %s is not allowed to create a city", req.Initiator.Id, req.Initiator.Role),
		)
	}

	CountryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, problems.InvalidArgumentError(ctx, fmt.Sprintf("country id is invalid"), &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	city, err := s.app.CreateCity(ctx, app.CreateCityInput{
		Name:      req.Name,
		CountryID: CountryID,
	})
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to create city")

		return nil, err
	}

	logger.Log(ctx, RequestID(ctx)).Infof("created city with ID %s by user %s", city.ID, req.Initiator.Id)

	return responses.City(city), nil
}
