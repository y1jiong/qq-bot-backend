package segment

import (
	"fmt"
	"github.com/bytedance/sonic"
	"qq-bot-backend/internal/util/codec"
	"regexp"
	"strings"
)

var (
	cqCodeRe = regexp.MustCompile(`\[CQ:(\w+),([^]]+)]`)
)

type messageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (ms *messageSegment) String() string {
	if ms == nil {
		return ""
	}

	if ms.Type == "text" {
		if text, ok := ms.Data["text"]; ok {
			return text
		}
		return ""
	}

	data := make([]string, 0, len(ms.Data))
	for k, v := range ms.Data {
		data = append(data, fmt.Sprintf("%s=%s", k, codec.EncodeCqCode(v)))
	}

	return fmt.Sprintf("[CQ:%s,%s]", ms.Type, strings.Join(data, ","))
}

type messageSegments []*messageSegment

func (mss *messageSegments) String() string {
	var result strings.Builder

	for _, segment := range *mss {
		result.WriteString(segment.String())
	}

	return result.String()
}

func ParseMessage(message string) messageSegments {
	var segments messageSegments
	matches := cqCodeRe.Split(message, -1)
	cqMatches := cqCodeRe.FindAllString(message, -1)

	for i, match := range matches {
		if match != "" {
			segments = append(segments, newTextSegment(match))
		}
		if i < len(cqMatches) {
			segment := newCQCodeSegment(cqMatches[i])
			if segment != nil {
				segments = append(segments, segment)
			}
		}
	}

	return segments
}

func ParseJSON(jsonBytes []byte) (messageSegments, error) {
	var segments messageSegments
	if err := sonic.ConfigStd.Unmarshal(jsonBytes, &segments); err != nil {
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
	matches := cqCodeRe.FindStringSubmatch(cqCode)
	if len(matches) < 3 {
		return nil
	}

	cqType := matches[1]
	dataStr := matches[2]

	data := make(map[string]string)
	for _, param := range strings.Split(dataStr, ",") {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			data[kv[0]] = codec.DecodeCqCode(kv[1])
		}
	}

	return &messageSegment{
		Type: cqType,
		Data: data,
	}
}
