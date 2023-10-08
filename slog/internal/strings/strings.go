package strings

import (
	"strings"
)

func AddPrefix(delimiter rune, leaf string, prefix ...string) string {
	var sb strings.Builder
	for i, p := range prefix {
		if p != "" {
			if i > 0 {
				sb.WriteRune(delimiter)
			}
			sb.WriteString(p)
		}
	}
	if sb.Len() > 0 {
		sb.WriteRune(delimiter)
	}
	sb.WriteString(leaf)
	return sb.String()
}

func AddPrefixFaster(delimiter rune, leaf string, prefix string) string {
	if prefix == "" {
		if leaf == "" {
			return ""
		}
		return leaf
	}

	// return prefix + string(delimiter) + leaf

	var sb strings.Builder
	sb.Grow(len(prefix) + 1 + len(leaf))
	sb.WriteString(prefix)
	sb.WriteRune(delimiter)
	sb.WriteString(leaf)
	return sb.String()
}

func DotPrefix(leaf string, prefix ...string) string {
	if len(prefix) == 1 {
		return AddPrefixFaster('.', leaf, prefix[0])
	}
	return AddPrefix('.', leaf, prefix...)
}
