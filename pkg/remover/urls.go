package remover

import "regexp"

func RemoveUrls(text string) string {
	urlRegex := regexp.MustCompile(`(https?|ftp)://[^\s/$.?#].[^\s]*`)
	return urlRegex.ReplaceAllString(text, "")
}
