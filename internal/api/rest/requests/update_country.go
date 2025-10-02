package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/cities-svc/resources"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateCountry(r *http.Request) (req resources.UpdateCountry, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UpdateCountryType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "country_id") != req.Data.Id.String() {
		errs["data/id"] = fmt.Errorf("country_id in path and body must match")
	}

	return req, errs.Filter()
}
