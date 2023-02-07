package codec

import "strings"

func (s *sCodec) DecodeBlank(src string) (dest string) {
	dest = strings.Replace(src, "%20", " ", -1)
	dest = strings.Replace(dest, "%25", "%", -1)
	return
}
