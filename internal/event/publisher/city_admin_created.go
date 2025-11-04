package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/event/contracts"
	"github.com/google/uuid"
)

type CreatedCityAdminData struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	Position  *string
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const CityAdminCreatedEvent = "city.admin.create"

func (s Service) PublishCityAdminCreated(
	ctx context.Context,
	admin models.CityAdmin,
) error {
	env := events.Envelope[CreatedCityAdminData]{
		Event:     CityAdminCreatedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: CreatedCityAdminData{
			UserID:    admin.UserID,
			CityID:    admin.CityID,
			Role:      admin.Role,
			Position:  admin.Position,
			Label:     admin.Label,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		admin.UserID.String(),
		env,
	)
}
