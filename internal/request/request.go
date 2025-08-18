package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var newLine = []byte("\r\n")

const (
	stateInit int = iota
	stateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       int
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateInit {
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return n, err
		}
		if requestLine != nil {
			r.RequestLine = *requestLine
			r.state = stateDone
			return n, nil
		}
	}
	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state: stateInit,
	}

	var buf bytes.Buffer
	for request.state != stateDone {
		_, err := buf.ReadFrom(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read from reader: %w", err)
		}

		n, err := request.parse(buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to parse request: %w", err)
		}
		if n == 0 {
			continue
		}
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, newLine)
	if idx == -1 {
		return nil, 0, nil
	}

	parts := strings.Split(string(data[:idx]), " ")
	idx += len(newLine)
	if len(parts) != 3 {
		return nil, idx, errors.New("malformed request line")
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return nil, idx, errors.New("malformed request line: invalid method")
	}

	if parts[2] != "HTTP/1.1" {
		return nil, idx, errors.New("malformed request line: invalid protocol")
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.Split(parts[2], "/")[1],
	}, idx, nil
}
