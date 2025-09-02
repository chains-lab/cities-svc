package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/cities-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateCity(r *http.Request) (req resources.UpdateCity, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"/data/id":        validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UpdateCityType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
