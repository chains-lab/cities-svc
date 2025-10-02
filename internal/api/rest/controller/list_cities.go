package controller

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func (a Service) ListCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters city.FilterParams

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &[]string{name}[0]
	}

	if sts := strings.TrimSpace(q.Get("status")); sts != "" {
		filters.Status = &[]string{sts}[0]
	}

	if cid := strings.TrimSpace(q.Get("country_id")); cid != "" {
		id, err := uuid.Parse(cid)
		if err != nil {
			ape.RenderErr(w, problems.InvalidParameter("country_id", err))
			return
		}
		filters.CountryID = &id
	}

	latStr, lonStr := q.Get("lat"), q.Get("lon")
	radMStr, radKMStr := q.Get("radius_m"), q.Get("radius_km")

	if (latStr != "" || lonStr != "") || (radMStr != "" || radKMStr != "") {
		if latStr == "" || lonStr == "" {
			ape.RenderErr(w, problems.InvalidParameter("lat/lon", fmt.Errorf("both lat and lon are required when using radius")))
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil || math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 {
			ape.RenderErr(w, problems.InvalidParameter("lat", fmt.Errorf("invalid latitude")))
			return
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil || math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 {
			ape.RenderErr(w, problems.InvalidParameter("lon", fmt.Errorf("invalid longitude")))
			return
		}

		var radiusM uint
		switch {
		case radKMStr != "":
			km, err := strconv.ParseFloat(radKMStr, 64)
			if err != nil || !(km > 0) {
				ape.RenderErr(w, problems.InvalidParameter("radius_km", fmt.Errorf("must be > 0")))
				return
			}
			radiusM = uint(math.Round(km * 1000.0))
		case radMStr != "":
			rm, err := strconv.ParseUint(radMStr, 10, 64)
			if err != nil || rm == 0 {
				ape.RenderErr(w, problems.InvalidParameter("radius_m", fmt.Errorf("must be > 0")))
				return
			}
			radiusM = uint(rm)
		default:
			ape.RenderErr(w, problems.InvalidParameter("radius", fmt.Errorf("radius_km or radius_m is required with lat/lon")))
			return
		}

		filters.Location = &city.FilterDistance{
			Point:   orb.Point{lon, lat},
			RadiusM: uint64(radiusM),
		}
	}

	page, size := pagi.GetPagination(r)

	cities, err := a.domain.city.Filter(ctx, filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to search cities")
		switch {
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.InvalidParameter("status", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CitiesCollection(cities))
}
