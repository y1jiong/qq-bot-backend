package codec

import (
	"fmt"
	"net/url"
	"strings"
)

func DecodeBlank(src string) (dest string) {
	dest = strings.ReplaceAll(src, "%20", " ")
	dest = strings.ReplaceAll(dest, "%25", "%")
	return
}

func GetRouteURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
}
