package codec

import "strings"

func DecodeCQCode(src string) (dest string) {
	dest = strings.ReplaceAll(src, "&#91;", "[")
	dest = strings.ReplaceAll(dest, "&#93;", "]")
	dest = strings.ReplaceAll(dest, "&#44;", ",")
	// 必须最后一个
	dest = strings.ReplaceAll(dest, "&amp;", "&")
	return
}

func EncodeCQCode(src string) (dest string) {
	// 必须第一个
	dest = strings.ReplaceAll(src, "&", "&amp;")
	dest = strings.ReplaceAll(dest, "[", "&#91;")
	dest = strings.ReplaceAll(dest, "]", "&#93;")
	dest = strings.ReplaceAll(dest, ",", "&#44;")
	return
}

func IsIncludeCQCode(str string) bool {
	return strings.Contains(str, "[CQ:")
}
