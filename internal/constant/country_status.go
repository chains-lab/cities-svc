package constant

import "fmt"

const (
	CountryStatusSupported   = "supported"
	CountryStatusUnsupported = "unsupported"
)

var countryStatuses = []string{
	CountryStatusSupported,
	CountryStatusUnsupported,
}

var ErrorCountryStatusNotSupported = fmt.Errorf("invalid country status must be one of: %v", GetAllCountriesStatuses())

func CheckCountryStatus(status string) error {
	for _, s := range countryStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorCountryStatusNotSupported)
}

func GetAllCountriesStatuses() []string {
	return countryStatuses
}
