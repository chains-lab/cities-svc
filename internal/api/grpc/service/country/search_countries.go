package country

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"

	"github.com/chains-lab/cities-svc/internal/problems"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SearchCountries(ctx context.Context, req *svc.SearchCountriesRequest) (*svc.CountriesList, error) {
	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problems.InvalidArgumentError(ctx, "invalid country status", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	countries, pag, err := s.app.SearchCountries(ctx, req.NameLike, status, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to search countries")

		return nil, err
	}

	return responses.CountriesList(countries, pag), nil
}
