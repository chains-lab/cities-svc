package enum

import "fmt"

const (
	FormToCreateCityStatusPending  = "pending"
	FormToCreateCityStatusAccepted = "accepted"
	FormToCreateCityStatusRejected = "rejected"
)

var formToCreateCityStatuses = []string{
	FormToCreateCityStatusPending,
	FormToCreateCityStatusAccepted,
	FormToCreateCityStatusRejected,
}

var ErrorFormToCreateCityStatusNotSupported = fmt.Errorf("'form to create city' status must be one of: %v", GetAllFormToCreateCityStatuses())

func ParseFormToCreateCityStatus(status string) (string, error) {
	for _, s := range formToCreateCityStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorFormToCreateCityStatusNotSupported)
}

func GetAllFormToCreateCityStatuses() []string {
	return formToCreateCityStatuses
}
