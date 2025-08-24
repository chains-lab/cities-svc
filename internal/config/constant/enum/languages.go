package enum

import "fmt"

const (
	LanguageEnglish    = "en"
	LanguageSpanish    = "es"
	LanguageFrench     = "fr"
	LanguageGerman     = "de"
	LanguageItalian    = "it"
	LanguagePortuguese = "pt"
	LanguageUkrainian  = "uk"
)

var languages = []string{
	LanguageEnglish,
	LanguageSpanish,
	LanguageFrench,
	LanguageGerman,
	LanguageItalian,
	LanguagePortuguese,
	LanguageUkrainian,
}

var ErrorLanguageNotSupported = fmt.Errorf("language not supported must be one of: %v", GetAllLanguages())

func ParseLanguage(lang string) (string, error) {
	for _, l := range languages {
		if l == lang {
			return l, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", lang, ErrorLanguageNotSupported)
}

func GetAllLanguages() []string {
	return languages
}
