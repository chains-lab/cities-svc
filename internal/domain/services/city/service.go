package city

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{
		db: db,
	}
}

var slugRegexp = regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`)
var nameRegexp = regexp.MustCompile(`^[A-Za-z -]+$`)

func validateTimezone(tz string) error {
	if tz == "" {
		return errx.ErrorInvalidTimeZone.Raise(
			fmt.Errorf("timezone must not be empty"),
		)
	}
	_, err := time.LoadLocation(tz)
	if err != nil {
		return errx.ErrorInvalidTimeZone.Raise(
			fmt.Errorf("invalid timezone: %s", tz),
		)
	}
	return nil
}

func validatePoint(p orb.Point) error {
	lon, lat := p[0], p[1]

	if lon < -180 || lon > 180 {
		return errx.ErrorInvalidPoint.Raise(
			fmt.Errorf("invalid longitude: %.6f (must be between -180 and 180)", lon),
		)
	}
	if lat < -90 || lat > 90 {
		return errx.ErrorInvalidPoint.Raise(
			fmt.Errorf("invalid latitude: %.6f (must be between -90 and 90)", lat),
		)
	}
	return nil
}

func validateSlug(slug string) error {
	if slug == "" {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("slug must not be empty"),
		)
	}
	if !slugRegexp.MatchString(slug) {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("invalid slug: %s", slug),
		)
	}
	return nil
}

func validateName(name string) error {
	if strings.Trim(name, " ") == "" {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("city name must not be empty"),
		)
	}
	if !nameRegexp.MatchString(name) {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("invalid city name: %s", name),
		)
	}
	return nil
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetCountryByID(ctx context.Context, id uuid.UUID) (models.Country, error)

	CreateCity(ctx context.Context, m models.City) (models.City, error)

	GetCityByID(ctx context.Context, id uuid.UUID) (models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)
	GetCityByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error)

	FilterCities(ctx context.Context, filter FilterParams, page, size uint64) (models.CitiesCollection, error)

	UpdateCity(ctx context.Context, id uuid.UUID, m UpdateParams, updatedAt time.Time) error
	UpdateCityStatus(ctx context.Context, id uuid.UUID, status string, updatedAt time.Time) error

	DeleteGovForCity(ctx context.Context, cityID uuid.UUID) error
}

func (s Service) CountryIsSupported(ctx context.Context, countryID uuid.UUID) error {
	country, err := s.db.GetCountryByID(ctx, countryID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get country by ID: %w", err),
		)
	}
	if country.IsNil() {
		return errx.ErrorCountryNotFound.Raise(
			fmt.Errorf("country with ID %s not found", countryID),
		)
	}
	if country.Status != enum.CountryStatusSupported {
		return errx.ErrorCountryNotSupported.Raise(
			fmt.Errorf("country with ID %s is not active", countryID),
		)
	}

	return nil
}
