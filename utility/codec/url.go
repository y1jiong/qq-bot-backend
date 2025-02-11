package codec

import (
	"fmt"
	"net/url"
	"strings"
)

var blankDecoder = strings.NewReplacer(
	"%20", " ",
	"%25", "%",
)

func DecodeBlank(src string) string {
	return blankDecoder.Replace(src)
}

func GetRouteURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
}
