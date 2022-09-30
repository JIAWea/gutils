package httpx

import (
	"bytes"
	"sort"
)

func OrderParam(p map[string]string, bizKey string) string {
	keys := make([]string, 0, len(p))
	for k := range p {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		if p[k] == "" {
			continue
		}

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(p[k])
	}
	buf.WriteString(bizKey)
	return buf.String()
}
