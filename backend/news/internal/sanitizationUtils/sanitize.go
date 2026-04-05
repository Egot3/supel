package sanitizationutils

import "strings"

func Slugify(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	parts := strings.Split(filename, ".")
	base := strings.ToLower(parts[0])
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, base)

	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
