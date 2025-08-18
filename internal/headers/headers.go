package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, []byte(crlf))
	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return len(crlf) * 2, true, nil
	}

	fieldName, fieldValue, ok := strings.Cut(string(data[:idx]), ":")
	if !ok || len(fieldName) != len(strings.TrimSpace(fieldName)) {
		return 0, false, errors.New("malformed field line")
	}

	h.Set(fieldName, strings.TrimSpace(fieldValue))
	return len(data[:idx]) + len(crlf), false, nil
}
