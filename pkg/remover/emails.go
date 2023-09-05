package remover

import "regexp"

func RemoveEmails(text string) string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return emailRegex.ReplaceAllString(text, "")
}
