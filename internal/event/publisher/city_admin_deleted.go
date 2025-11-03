package publisher

import (
	"context"
	"time"

	events "github.com/chains-lab/cities-svc/internal/event/contracts"
	"github.com/google/uuid"
)

type DeletedCityAdminData struct {
	UserID uuid.UUID
	CityID uuid.UUID
}

const CityAdminEventDeleted = "city.admin.deleted"

func (s Service) PublishCityAdminDeleted(
	ctx context.Context,
	userID uuid.UUID,
	cityID uuid.UUID,
) error {
	env := events.Envelope[DeletedCityAdminData]{
		Event:     CityAdminEventDeleted,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: DeletedCityAdminData{
			UserID: userID,
			CityID: cityID,
		},
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		userID.String(),
		env,
	)
}
