package country

import (
	"context"

	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) AdminUpdateCountryStatus(ctx context.Context, req *svc.UpdateCountryStatusRequest) (*countryProto.Country, error) {
	user := meta.User(ctx)

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid country id format", &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "request id", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	country, err := s.app.UpdateCountryStatus(ctx, countryID, status)
	if err != nil {
		return nil, err
	}

	logger.Log(ctx).Infof("country status updated by user %s for country ID %s", user.ID, country.ID)

	return responses.Country(country), nil
}
