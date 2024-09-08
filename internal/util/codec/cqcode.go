package codec

import "strings"

func DecodeCqCode(src string) (dest string) {
	dest = strings.ReplaceAll(src, "&#91;", "[")
	dest = strings.ReplaceAll(dest, "&#93;", "]")
	dest = strings.ReplaceAll(dest, "&#44;", ",")
	// 必须最后一个
	dest = strings.ReplaceAll(dest, "&amp;", "&")
	return
}

func EncodeCqCode(src string) (dest string) {
	// 必须第一个
	dest = strings.ReplaceAll(src, "&", "&amp;")
	dest = strings.ReplaceAll(dest, "[", "&#91;")
	dest = strings.ReplaceAll(dest, "]", "&#93;")
	dest = strings.ReplaceAll(dest, ",", "&#44;")
	return
}

func IsIncludeCqCode(str string) bool {
	return strings.Contains(str, "[CQ:")
}
