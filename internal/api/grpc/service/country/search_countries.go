package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SearchCountries(ctx context.Context, req *svc.SearchCountriesRequest) (*svc.CountriesList, error) {
	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problem.InvalidArgumentError(ctx, "invalid country status", &errdetails.BadRequest_FieldViolation{
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

	return response.CountriesList(countries, pag), nil
}
