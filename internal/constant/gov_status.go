package constant

import "fmt"

const (
	GovStatusActive   = "active"
	GovStatusInactive = "inactive"
)

var AllGovStatuses = []string{
	GovStatusActive,
	GovStatusInactive,
}

var ErrorInvalidGovStatus = fmt.Errorf("invalid government status must be one of: %s", GetAllGovStatuses())

func CheckGovStatus(status string) error {
	for _, s := range AllGovStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorInvalidGovStatus)
}

func GetAllGovStatuses() []string {
	return AllGovStatuses
}
