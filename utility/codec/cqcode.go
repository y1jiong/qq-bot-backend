package codec

import "strings"

var cqCodeDecoder = strings.NewReplacer(
	"&#91;", "[",
	"&#93;", "]",
	"&#44;", ",",
	"&amp;", "&",
)

var cqCodeEncoder = strings.NewReplacer(
	"&", "&amp;",
	"[", "&#91;",
	"]", "&#93;",
	",", "&#44;",
)

func DecodeCQCode(src string) string {
	return cqCodeDecoder.Replace(src)
}

func EncodeCQCode(src string) string {
	return cqCodeEncoder.Replace(src)
}

func IsIncludeCQCode(str string) bool {
	return strings.Contains(str, "[CQ:")
}
