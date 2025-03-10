package segment

func FilterCQCode(message string) string {
	if message == "" {
		return ""
	}
	return cqCodeRe.ReplaceAllString(message, "")
}
