package segment

import (
	json "github.com/bytedance/sonic"
)

func ToString(message any) string {
	if message == nil {
		return ""
	}
	if msg, ok := message.(string); ok {
		return msg
	}
	if segments, ok := message.(interface{ String() string }); ok {
		return segments.String()
	}

	segments, err := parseAny(message)
	if err != nil {
		return ""
	}
	return segments.String()
}

func ParseAny(message any) (messageSegments, error) {
	if message == nil {
		return nil, nil
	}
	if msg, ok := message.(string); ok {
		return ParseMessage(msg), nil
	}

	return parseAny(message)
}

func parseAny(message any) (messageSegments, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return ParseJSON(bytes)
}

func ParseJSON(jsonBytes []byte) (messageSegments, error) {
	var segments messageSegments
	if err := json.Unmarshal(jsonBytes, &segments); err != nil {
		return nil, err
	}
	return segments, nil
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
