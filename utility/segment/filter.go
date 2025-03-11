package segment

func FilterCQCode(message []byte) []byte {
	if len(message) == 0 {
		return message
	}
	return cqCodeRe.ReplaceAll(message, []byte(nil))
}
