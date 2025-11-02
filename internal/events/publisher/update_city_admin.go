package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/events"
	"github.com/google/uuid"
)

type CityAdminUpdate struct {
	UserID uuid.UUID  `json:"user_id"`
	CityID *uuid.UUID `json:"city_id"`
	Role   *string    `json:"role"`
}

func (s *Service) UpdateCityAdmin(
	ctx context.Context,
	userID uuid.UUID,
	cityID *uuid.UUID,
	role *string,
) error {
	env := events.Envelope[CityAdminUpdate]{
		Event:     "city_admin.update",
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: CityAdminUpdate{
			UserID: userID,
			CityID: cityID,
			Role:   role,
		},
	}
	return s.publish(
		ctx,
		events.TopicCompaniesEmployeeV1,
		userID.String(),
		env,
	)
}
