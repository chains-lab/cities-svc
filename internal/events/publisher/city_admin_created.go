package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type CreatedCityAdminData struct {
	City       models.City       `json:"city"`
	Admin      models.CityAdmin  `json:"admin"`
	Recipients PayloadRecipients `json:"recipients"`
}

const CityAdminCreatedEvent = "city.admin.create"

func (s Service) PublishCityAdminCreated(
	ctx context.Context,
	admin models.CityAdmin,
	city models.City,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		admin.UserID.String(),
		events.Envelope[CreatedCityAdminData]{
			Event:     CityAdminCreatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CreatedCityAdminData{
				City:  city,
				Admin: admin,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
