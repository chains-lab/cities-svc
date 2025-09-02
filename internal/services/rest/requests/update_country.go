package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/cities-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateCountry(r *http.Request) (req resources.UpdateCountry, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"/data/id":        validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UpdateCountryType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
