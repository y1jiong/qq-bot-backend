package segment

import (
	"qq-bot-backend/utility/codec"
	"regexp"
	"strings"
)

const (
	TypeAt    = "at"
	TypeText  = "text"
	TypeImage = "image"
	TypeReply = "reply"
)

var (
	cqCodeRe = regexp.MustCompile(`\[CQ:(\w+)(?:,([^]]+))?]`)
)

type messageSegments []*messageSegment

func (segments messageSegments) First() *messageSegment {
	if len(segments) == 0 {
		return nil
	}
	return segments[0]
}

func (segments messageSegments) String() string {
	var result strings.Builder

	for _, segment := range segments {
		result.WriteString(segment.String())
	}

	return result.String()
}

type messageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (segment *messageSegment) String() string {
	if segment == nil {
		return ""
	}

	if segment.Type == TypeText {
		if text, ok := segment.Data[TypeText]; ok {
			return text
		}
		return ""
	}

	if len(segment.Data) == 0 {
		return "[CQ:" + segment.Type + "]"
	}

	data := make([]string, 0, len(segment.Data))
	for k, v := range segment.Data {
		data = append(data, k+"="+codec.EncodeCQCode(v))
	}

	return "[CQ:" + segment.Type + "," + strings.Join(data, ",") + "]"
}
