package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func GenerateSlug(text string) string {
	var normalized strings.Builder
	for _, r := range strings.ToLower(text) {
		if norm := removeAccent(r); norm != "" {
			normalized.WriteString(norm)
		}
	}

	slug := normalized.String()

	slug = regexp.MustCompile(`[\s_]+`).ReplaceAllString(slug, "-")

	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")

	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

func removeAccent(r rune) string {
	// Convert accented characters to base
	switch unicode.ToLower(r) {
	case 'à', 'á', 'â', 'ã', 'ä', 'å':
		return "a"
	case 'è', 'é', 'ê', 'ë':
		return "e"
	case 'ì', 'í', 'î', 'ï':
		return "i"
	case 'ò', 'ó', 'ô', 'õ', 'ö':
		return "o"
	case 'ù', 'ú', 'û', 'ü':
		return "u"
	case 'ç':
		return "c"
	}
	return string(r)
}
