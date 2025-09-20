package city

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/paulmach/orb"
)

type City struct {
	citiesQ dbx.CitiesQ

	slugRegexp *regexp.Regexp
	nameRegexp *regexp.Regexp
}

func NewCity(db *sql.DB) City {
	return City{
		citiesQ:    dbx.NewCitiesQ(db),
		slugRegexp: regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`),
		nameRegexp: regexp.MustCompile(`^[A-Za-z -]+$`),
	}
}

func (c City) validateTimezone(tz string) error {
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

func (c City) validatePoint(p orb.Point) error {
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

func (c City) validateSlug(slug string) error {
	if slug == "" {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("slug must not be empty"),
		)
	}
	if !c.slugRegexp.MatchString(slug) {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("invalid slug: %s", slug),
		)
	}
	return nil
}

func (c City) validateName(name string) error {
	if strings.Trim(name, " ") == "" {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("city name must not be empty"),
		)
	}
	if !c.nameRegexp.MatchString(name) {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("invalid city name: %s", name),
		)
	}
	return nil
}

func cityFromDb(c dbx.City) models.City {
	res := models.City{
		ID:        c.ID,
		CountryID: c.CountryID,
		Point:     c.Point,
		Status:    c.Status,
		Name:      c.Name,
		Timezone:  c.Timezone,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.Icon.Valid {
		res.Icon = &c.Icon.String
	}
	if c.Slug.Valid {
		res.Slug = &c.Slug.String
	}

	return res
}
