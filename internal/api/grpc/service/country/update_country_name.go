package country

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) UpdateCountryName(ctx context.Context, req *svc.UpdateCountryNameRequest) (*svc.Country, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, responses.AppError(ctx, RequestID(ctx), errx.RaiseInternal(err))
	}

	if role != roles.Admin && role != roles.SuperUser {
		logger.Log(ctx, RequestID(ctx)).Warnf("user %s with role %s tried to update a country, but only admins and superusers can update countries",
			req.Initiator.Id, req.Initiator.Role)

		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"user %s with role %s is not allowed to update a country", req.Initiator.Id, req.Initiator.Role),
		)
	}

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.methods.UpdateCountryName(ctx, countryID, req.Name)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to create country")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.Country(country), nil
}
