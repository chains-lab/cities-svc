package country

import (
	"context"

	countryProto "github.com/chains-lab/cities-proto/gen/go/svc/country"
	svc "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"
)

func (s Service) CreateCountry(ctx context.Context, req *svc.CreateCountryRequest) (*countryProto.Country, error) {
	user := meta.User(ctx)

	country, err := s.app.CreateCountry(ctx, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create country")

		return nil, err
	}

	logger.Log(ctx).Infof("created country with ID %s by user %s", country.ID, user.ID)
	return responses.Country(country), nil
}
