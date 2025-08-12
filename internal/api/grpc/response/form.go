package response

import (
	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func Form(form models.Form) *svc.FormToCreateCity {
	return &svc.FormToCreateCity{
		Id:             form.ID.String(),
		Status:         form.Status,
		CityName:       form.CityName,
		CountryId:      form.CountryID.String(),
		InitiatorId:    form.InitiatorID.String(),
		ContactEmail:   form.ContactEmail,
		ContactPhone:   form.ContactPhone,
		Text:           form.Text,
		UserReviewedId: form.UserRevID.String(),
		CreatedAt:      form.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
