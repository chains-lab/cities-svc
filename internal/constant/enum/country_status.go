package enum

import "fmt"

const (
	CountryStatusSupported   = "supported"
	CountryStatusSuspended   = "suspended"
	CountryStatusUnsupported = "unsupported"
)

var countryStatuses = []string{
	CountryStatusSupported,
	CountryStatusSuspended,
	CountryStatusUnsupported,
}

var ErrorCountryStatusNotSupported = fmt.Errorf("invalid country status must be one of: %v", GetAllCountriesStatuses())

func ParseCountryStatus(status string) (string, error) {
	for _, s := range countryStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorCountryStatusNotSupported)
}

func GetAllCountriesStatuses() []string {
	return countryStatuses
}
