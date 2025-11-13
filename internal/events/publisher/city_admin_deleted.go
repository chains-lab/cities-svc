package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type DeletedCityAdminData struct {
	CityAdmin  models.CityAdmin  `json:"city_admin"`
	City       models.City       `json:"city"`
	Recipients PayloadRecipients `json:"recipients"`
}

const CityAdminEventDeleted = "city.admin.deleted"

func (s Service) PublishCityAdminDeleted(
	ctx context.Context,
	admin models.CityAdmin,
	city models.City,
	recipients []uuid.UUID,
) error {
	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		admin.UserID.String(),
		events.Envelope[DeletedCityAdminData]{
			Event:     CityAdminEventDeleted,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: DeletedCityAdminData{
				CityAdmin: admin,
				City:      city,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
