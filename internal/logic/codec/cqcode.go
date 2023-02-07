package codec

import (
	"strings"
)

func (s *sCodec) DecodeCqCode(src string) (dest string) {
	dest = strings.Replace(src, "&#91;", "[", -1)
	dest = strings.Replace(dest, "&#93;", "]", -1)
	dest = strings.Replace(dest, "&#44;", ",", -1)
	// 必须最后一个
	dest = strings.Replace(dest, "&amp;", "&", -1)
	return
}

func (s *sCodec) EncodeCqCode(src string) (dest string) {
	// 必须第一个
	dest = strings.Replace(src, "&", "&amp;", -1)
	dest = strings.Replace(dest, "[", "&#91;", -1)
	dest = strings.Replace(dest, "]", "&#93;", -1)
	dest = strings.Replace(dest, ",", "&#44;", -1)
	return
}

func (s *sCodec) IsIncludeCqCode(str string) (yes bool) {
	return strings.Contains(str, "[CQ:")
}
