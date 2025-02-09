package segment

import (
	"github.com/bytedance/sonic"
	"qq-bot-backend/utility/codec"
	"regexp"
	"strings"
)

var (
	cqCodeRe = regexp.MustCompile(`\[CQ:(\w+)(?:,([^]]+))?]`)
)

type messageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (segment *messageSegment) String() string {
	if segment == nil {
		return ""
	}

	if segment.Type == "text" {
		if text, ok := segment.Data["text"]; ok {
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

type messageSegments []*messageSegment

func (segments *messageSegments) String() string {
	var result strings.Builder

	for _, segment := range *segments {
		result.WriteString(segment.String())
	}

	return result.String()
}

func ParseMessage(message string) messageSegments {
	var segments messageSegments
	idxes := cqCodeRe.FindAllStringIndex(message, -1)

	lastEnd := 0
	for _, idx := range idxes {
		if lastEnd < idx[0] {
			if text := message[lastEnd:idx[0]]; text != "" {
				segments = append(segments, newTextSegment(text))
			}
		}

		if segment := newCQCodeSegment(message[idx[0]:idx[1]]); segment != nil {
			segments = append(segments, segment)
		}

		lastEnd = idx[1]
	}

	if lastEnd < len(message) {
		if text := message[lastEnd:]; text != "" {
			segments = append(segments, newTextSegment(text))
		}
	}

	return segments
}

func ParseJSON(jsonBytes []byte) (messageSegments, error) {
	var segments messageSegments
	if err := sonic.Unmarshal(jsonBytes, &segments); err != nil {
		return nil, err
	}
	return segments, nil
}

func NewTextSegment(text string) messageSegments {
	return messageSegments{newTextSegment(text)}
}

func newTextSegment(text string) *messageSegment {
	return &messageSegment{
		Type: "text",
		Data: map[string]string{"text": text},
	}
}

func newCQCodeSegment(cqCode string) *messageSegment {
	idxes := cqCodeRe.FindStringSubmatchIndex(cqCode)
	if len(idxes) < 6 {
		return nil
	}

	cqType := cqCode[idxes[2]:idxes[3]]
	var dataStr string
	if idxes[4] != -1 && idxes[5] != -1 {
		dataStr = cqCode[idxes[4]:idxes[5]]
	}

	data := make(map[string]string)
	beg := 0
	for beg < len(dataStr) {
		var param string
		if end := strings.IndexByte(dataStr[beg:], ','); end != -1 {
			end += beg
			param = dataStr[beg:end]
			beg = end + 1
		} else {
			param = dataStr[beg:]
			beg = len(dataStr)
		}

		if idx := strings.IndexByte(param, '='); idx != -1 {
			data[param[:idx]] = codec.DecodeCQCode(param[idx+1:])
		}
	}

	return &messageSegment{
		Type: cqType,
		Data: data,
	}
}
