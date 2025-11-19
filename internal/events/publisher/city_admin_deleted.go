package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type DeletedCityAdminData struct {
	CityAdmin  models.CityAdmin   `json:"city_admin"`
	City       models.City        `json:"city"`
	Recipients *PayloadRecipients `json:"recipients,omitempty"`
}

const CityAdminEventDeleted = "city.admin.deleted"

func (s Service) PublishCityAdminDeleted(
	ctx context.Context,
	admin models.CityAdmin,
	city models.City,
	recipients ...uuid.UUID,
) error {
	event := events.Envelope[DeletedCityAdminData]{
		Event:     CityAdminEventDeleted,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: DeletedCityAdminData{
			CityAdmin: admin,
			City:      city,
		},
	}
	if len(recipients) > 0 {
		event.Data.Recipients = &PayloadRecipients{
			Users: recipients,
		}
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		fmt.Sprintf("%s:%s", admin.UserID.String(), city.ID.String()),
		event,
	)
}
