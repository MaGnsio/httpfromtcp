package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLines := strings.Split(string(r), "\r\n")

	if len(requestLines) < 1 {
		return nil, errors.New("empty request")
	}

	requestLine := requestLines[0]
	requestLineParts := strings.Split(requestLine, " ")

	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request line")
	}

	if requestLineParts[0] != strings.ToUpper(requestLineParts[0]) {
		return nil, errors.New("invalid request method")
	}

	if requestLineParts[2] != "HTTP/1.1" {
		return nil, errors.New("invalid http version")
	}

	return &Request{
		RequestLine: RequestLine{
			Method:        requestLineParts[0],
			RequestTarget: requestLineParts[1],
			HttpVersion:   strings.TrimPrefix(requestLineParts[2], "HTTP/"),
		},
	}, nil
}
