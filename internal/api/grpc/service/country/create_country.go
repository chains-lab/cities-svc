package country

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) CreateCountry(ctx context.Context, req *svc.CreateCountryRequest) (*svc.Country, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, responses.AppError(ctx, RequestID(ctx), errx.RaiseInternal(err))
	}

	if role != roles.Admin && role != roles.SuperUser {
		logger.Log(ctx, RequestID(ctx)).Warnf("user %s with role %s tried to create a country, but only admins and superusers can create countries",
			req.Initiator.Id, req.Initiator.Role)

		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"user %s with role %s is not allowed to create a country", req.Initiator.Id, req.Initiator.Role),
		)
	}

	country, err := s.app.CreateCountry(ctx, req.Name)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to create country")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	logger.Log(ctx, RequestID(ctx)).Infof("created country with ID %s by user %s", country.ID, req.Initiator.Id)
	return responses.Country(country), nil
}
