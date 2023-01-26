package module

import "strings"

func (s *sModule) DecodeCqCode(src string) (dst string) {
	dst = src
	dst = strings.Replace(dst, "&#91;", "[", -1)
	dst = strings.Replace(dst, "&#93;", "]", -1)
	dst = strings.Replace(dst, "&#44;", ",", -1)
	dst = strings.Replace(dst, "&amp;", "&", -1)
	return
}

func (s *sModule) EncodeCqCode(src string) (dst string) {
	dst = src
	dst = strings.Replace(dst, "&", "&amp;", -1)
	dst = strings.Replace(dst, "[", "&#91;", -1)
	dst = strings.Replace(dst, "]", "&#93;", -1)
	dst = strings.Replace(dst, ",", "&#44;", -1)
	return
}

func (s *sModule) IsIncludeCqCode(str string) (yes bool) {
	return strings.Contains(str, "[CQ:")
}
