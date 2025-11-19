package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type UpdatedCityAdminData struct {
	CityAdmin  models.CityAdmin   `json:"city_admin"`
	City       models.City        `json:"city"`
	Recipients *PayloadRecipients `json:"recipients,omitempty"`
}

const CityAdminUpdatedEvent = "city.admin.updated"

func (s Service) PublishCityAdminUpdated(
	ctx context.Context,
	admin models.CityAdmin,
	city models.City,
	recipients ...uuid.UUID,
) error {
	event := events.Envelope[UpdatedCityAdminData]{
		Event:     CityAdminUpdatedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: UpdatedCityAdminData{
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
