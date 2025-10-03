package enum

import "fmt"

const (
	InviteStatusSent     = "sent"
	InviteStatusAccepted = "accepted"
)

var allInviteStatuses = []string{
	InviteStatusSent,
	InviteStatusAccepted,
}

var ErrorInvalidInviteStatus = fmt.Errorf("invalid invite status")

func CheckInviteStatus(status string) error {
	for _, s := range allInviteStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorInvalidInviteStatus)
}

func GetAllInviteStatuses() []string {
	return allInviteStatuses
}
