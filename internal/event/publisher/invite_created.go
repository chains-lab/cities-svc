package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/event/contracts"
	"github.com/google/uuid"
)

type InviteCreatedData struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CityID    uuid.UUID `json:"city_id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

const InviteCreatedEvent = "city.admin.create"

func (s Service) PublishInviteCreated(
	ctx context.Context,
	admin models.Invite,
) error {
	env := events.Envelope[InviteCreatedData]{
		Event:     CityAdminCreatedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: InviteCreatedData{
			ID:        admin.ID,
			UserID:    admin.UserID,
			CityID:    admin.CityID,
			Status:    admin.Status,
			Role:      admin.Role,
			ExpiresAt: admin.ExpiresAt,
			CreatedAt: admin.CreatedAt,
		},
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		admin.UserID.String(),
		env,
	)
}
