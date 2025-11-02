package status

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/domain/models"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateStatus(ctx context.Context, status models.CityStatus) error
	GetStatusByID(ctx context.Context, ID string) (models.CityStatus, error)
	UpdateStatus(ctx context.Context, ID string, params UpdateParams) (models.CityStatus, error)
	DeleteStatus(ctx context.Context, ID string) error

	FilterStatuses(ctx context.Context, filters FilterStatusParams, page, size uint64) (models.CityStatusesCollection, error)
	GetStatusForCityID(ctx context.Context, cityID string) (models.CityStatus, error)

	ExistsCitiesWitStatus(ctx context.Context, ID string) (bool, error)
}
