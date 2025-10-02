package country

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
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

	CreateCountry(ctx context.Context, country models.Country) (models.Country, error)

	GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error)
	GetCountryByName(ctx context.Context, name string) (models.Country, error)

	FilterCountries(
		ctx context.Context,
		filters FilterParams,
		page, size uint64,
	) (models.CountriesCollection, error)

	UpdateCountry(
		ctx context.Context,
		countryID uuid.UUID,
		params UpdateParams,
		updatedAt time.Time,
	) error

	UpdateCountryStatus(ctx context.Context, countryID uuid.UUID, status string, updatedAt time.Time) error
	UpdateStatusForAllCountryCities(ctx context.Context, countryID uuid.UUID, status string, updatedAt time.Time) error
}
