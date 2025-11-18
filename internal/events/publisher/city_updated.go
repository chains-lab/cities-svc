package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type CityUpdatedData struct {
	City       models.City        `json:"city"`
	Recipients *PayloadRecipients `json:"recipients,omitempty"`
}

const CityUpdateEvent = "city.admin.update"

func (s Service) PublishCityUpdated(
	ctx context.Context,
	city models.City,
	recipients ...uuid.UUID,
) error {
	event := events.Envelope[CityUpdatedData]{
		Event:     CityUpdateEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: CityUpdatedData{
			City: city,
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
		city.ID.String(),
		event,
	)
}
