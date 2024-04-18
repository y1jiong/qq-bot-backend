package codec

import "strings"

func (s *sCodec) DecodeBlank(src string) (dest string) {
	dest = strings.ReplaceAll(src, "%20", " ")
	dest = strings.ReplaceAll(dest, "%25", "%")
	return
}
