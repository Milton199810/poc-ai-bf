package remover

import "regexp"

func RemoveEmojis(text string) string {
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}-\x{2B50}\x{1F900}-\x{1F9FF}\x{1F1E0}-\x{1F1FF}\x{1FAD0}-\x{1FAFF}\x{1F6F8}-\x{1F6FF}]`)
	return emojiRegex.ReplaceAllString(text, "")
}
