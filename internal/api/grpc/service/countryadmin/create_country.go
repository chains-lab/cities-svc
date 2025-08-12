package countryadmin

import (
	"context"

	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/countryadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
)

func (s Service) CreateCountry(ctx context.Context, req *svc.CreateCountryRequest) (*countryProto.Country, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "create profile",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	country, err := s.app.CreateCountry(ctx, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create country")

		return nil, err
	}

	logger.Log(ctx).Infof("created country with ID %s by user %s", country.ID, req.Initiator.UserId)
	return response.Country(country), nil
}
