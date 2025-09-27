package segment

import (
	"qq-bot-backend/utility/codec"
	"strings"
)

func NewTextSegments(text string) messageSegments {
	return messageSegments{newTextSegment(text)}
}

func NewAtSegments(qq string) messageSegments {
	return messageSegments{
		&messageSegment{
			Type: TypeAt,
			Data: map[string]string{"qq": qq},
		},
	}
}

func NewReplySegments(messageId string) messageSegments {
	return messageSegments{
		&messageSegment{
			Type: TypeReply,
			Data: map[string]string{"id": messageId},
		},
	}
}

func newTextSegment(text string) *messageSegment {
	return &messageSegment{
		Type: TypeText,
		Data: map[string]string{TypeText: text},
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
