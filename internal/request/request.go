package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func (rl *RequestLine) ValidMethod() bool {
	return rl.Method == strings.ToUpper(rl.Method)
}

func (rl *RequestLine) ValidHTTP() bool {
	return rl.HTTPVersion == "1.1"
}

var SEPARATOR = "\r\n"
var (
	ErrorMalformedRequestLine = fmt.Errorf("malformed request-line")
	ErrorInvalidMethod        = fmt.Errorf("invalid method")
	ErrorInvalidHTTPVersion   = fmt.Errorf("invalid HTTP version")
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(string(data))

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}

func parseRequestLine(r string) (*RequestLine, error) {
	requestLines := strings.Split(r, SEPARATOR)
	if len(requestLines) == 0 {
		return nil, errors.New("empty request")
	}

	parts := strings.Split(requestLines[0], " ")
	if len(parts) != 3 {
		return nil, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HTTPVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}

	if !rl.ValidMethod() {
		return nil, ErrorInvalidMethod
	}
	if !rl.ValidHTTP() {
		return nil, ErrorInvalidHTTPVersion
	}

	return rl, nil
}
