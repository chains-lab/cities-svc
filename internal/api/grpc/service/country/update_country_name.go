package country

import (
	"context"

	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCountryName(ctx context.Context, req *svc.UpdateCountryNameRequest) (*countryProto.Country, error) {
	user := meta.User(ctx)

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid country id format", &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.app.UpdateCountryName(ctx, countryID, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create country")

		return nil, err
	}

	logger.Log(ctx).Infof("country name updated by user %s for country ID %s", user.ID, country.ID)

	return responses.Country(country), nil
}
