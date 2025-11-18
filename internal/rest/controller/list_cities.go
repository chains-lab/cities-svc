package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/chains-lab/restkit/pagi"
	"github.com/paulmach/orb"
)

func (s Service) ListCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters city.FilterParams

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &[]string{name}[0]
	}

	if sts := strings.TrimSpace(q.Get("status")); sts != "" {
		filters.Status = &[]string{sts}[0]
	}

	if q.Get("country_id") != "" {
		countryID := q.Get("country_id")
		filters.CountryID = &countryID
	}

	latStr, lonStr := q.Get("lat"), q.Get("lon")
	radM := q.Get("radius")

	if (latStr != "" || lonStr != "") || radM != "" {
		if latStr == "" || lonStr == "" {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"lat/lon": fmt.Errorf("both lat and lon are required with radius"),
			})...)
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil || math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"lat": fmt.Errorf("invalid latitude"),
			})...)
			return
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil || math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"lon": fmt.Errorf("invalid longitude"),
			})...)
			return
		}

		//var radius uint
		//switch {
		//case radKMStr != "":
		//	km, err := strconv.ParseFloat(radKMStr, 64)
		//	if err != nil || !(km > 0) {
		//		ape.RenderErr(w, problems.BadRequest(validation.Errors{
		//			"radius_km": fmt.Errorf("must be > 0"),
		//		})...)
		//		return
		//	}
		//	radius = uint(math.Round(km * 1000.0))
		//case radM != "":
		//default:
		//	ape.RenderErr(w, problems.BadRequest(validation.Errors{
		//		"radius_m/radius_km": fmt.Errorf("one of radius_m or radius_km is required"),
		//	})...)
		//	return
		//}

		rm, err := strconv.ParseUint(radM, 10, 64)
		if err != nil || rm == 0 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"radius": fmt.Errorf("must be > 0"),
			})...)
			return
		}
		radius := uint(rm)

		filters.Location = &city.FilterDistance{
			Point:   orb.Point{lon, lat},
			RadiusM: uint64(radius),
		}
	}

	page, size := pagi.GetPagination(r)

	cities, err := s.domain.city.Filter(ctx, filters, page, size)
	if err != nil {
		s.log.WithError(err).Error("failed to search cities")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CitiesCollection(cities))
}
