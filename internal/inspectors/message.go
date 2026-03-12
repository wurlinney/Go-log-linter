package inspectors

import "unicode"

// FirstLetterLower проверяет, что первая буква строки строчная.
func FirstLetterLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			continue
		}
		return unicode.IsLower(r)
	}
	return true
}

// IsAllowedChar сообщает, разрешён ли символ в сообщении логов.
func IsAllowedChar(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	if unicode.IsSpace(r) {
		return true
	}
	return false
}
